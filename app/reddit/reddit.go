package reddit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/avast/retry-go"
	"github.com/sirupsen/logrus"
	confmodel "github.com/tigorlazuardi/ridit-go/app/config/models"
	"github.com/tigorlazuardi/ridit-go/app/reddit/models"
	"github.com/tigorlazuardi/ridit-go/pkg"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type Repository struct {
	client pkg.Doer
	config confmodel.Config
	bars   *mpb.Progress
}

type RepositoryError struct {
}

type ListingChan struct {
	Downloads []models.DownloadMeta
	Err       error
}

type DownloadChan struct {
	Download models.DownloadMeta
	Err      error
}

func NewRepository(client pkg.Doer, config confmodel.Config) Repository {
	return Repository{
		client: client,
		config: config,
		bars:   mpb.New(mpb.WithOutput(os.Stdout)),
	}
}

func (r Repository) Fetch(ctx context.Context) <-chan DownloadChan {
	lc := make(chan ListingChan, len(r.config.Subreddits))

	go func() {
		wg := sync.WaitGroup{}
		for subreddit, subconf := range r.config.Subreddits {
			wg.Add(1)
			go func(ctx context.Context, subreddit string, subconf confmodel.Subreddit) {
				defer wg.Done()
				ctx = pkg.ContextEntryWithFields(ctx, logrus.Fields{
					"subreddit":     subreddit,
					"configuration": subconf,
					"url":           fmt.Sprintf("https://reddit.com/r/%s/%s.json", subreddit, subconf.Sort),
				})
				meta, err := r.downloadListing(ctx, subreddit, subconf)
				lc <- ListingChan{
					Downloads: meta,
					Err:       err,
				}
			}(ctx, subreddit, subconf)
		}
		wg.Wait()
		close(lc)
	}()

	return r.downloadImages(ctx, lc)
}

func (r Repository) downloadImages(ctx context.Context, lc <-chan ListingChan) <-chan DownloadChan {
	dc := make(chan DownloadChan)
	go func(ctx context.Context) {
		entry := pkg.EntryFromContext(ctx)
		wg := sync.WaitGroup{}
		for listing := range lc {
			if listing.Err != nil {
				entry.WithError(listing.Err).Error(listing.Err)
				continue
			}
			wg.Add(1)
			go func(ctx context.Context, listing ListingChan) {
				defer wg.Done()
				wgg := &sync.WaitGroup{}
				for _, meta := range listing.Downloads {
					wgg.Add(1)
					go func(ctx context.Context, meta models.DownloadMeta) {
						defer wgg.Done()
						ctx = pkg.ContextEntryWithFields(ctx, logrus.Fields{"url": meta.URL})
						err := r.download(ctx, meta)
						dc <- DownloadChan{
							Download: meta,
							Err:      err,
						}
					}(ctx, meta)
				}
				wgg.Wait()
			}(ctx, listing)
		}

		wg.Wait()
		close(dc)
	}(ctx)
	return dc
}

func (r Repository) download(ctx context.Context, meta models.DownloadMeta) error {
	dir := filepath.Join(r.config.Download.Path, r.config.Profile, meta.SubredditName)
	_ = os.MkdirAll(dir, 0777)
	path := filepath.Join(dir, meta.Filename)
	entry := pkg.EntryFromContext(ctx).WithField("path", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0666)
	if err != nil {
		entry.Debug(err)
		return nil
	}
	f.Close()
	err = retry.Do(func() error {
		ctx, done := context.WithTimeout(ctx, r.config.Download.Timeout.Duration)
		defer done()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, meta.URL, nil)
		if err != nil {
			return retry.Unrecoverable(err)
		}
		res, err := r.client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		if res.StatusCode >= 500 {
			err := errors.New("source is unavailable")
			entry.WithField("status_code", res.StatusCode).Error(err)
			return retry.Unrecoverable(err)
		}

		if res.StatusCode >= 400 {
			err := errors.New("source returned status " + res.Status)
			entry.WithField("status_code", res.StatusCode).Error(err)
			return retry.Unrecoverable(err)
		}

		file, err := os.Create(path)
		if err != nil {
			return retry.Unrecoverable(err)
		}

		var reader = res.Body
		if isTerminal() {
			var bar *mpb.Bar
			leftDecor := mpb.PrependDecorators(
				decor.Name("["+meta.SubredditName+"]", decor.WCSyncSpaceR),
				decor.Name(meta.URL, decor.WCSyncSpaceR),
				decor.OnComplete(decor.Name("downloading", decor.WCSyncSpaceR), "done"),
			)
			if res.ContentLength < 0 {
				bar = r.bars.AddSpinner(res.ContentLength, leftDecor)
			} else {
				bar = r.bars.AddBar(res.ContentLength,
					leftDecor,
					mpb.AppendDecorators(
						decor.CountersKiloByte("%d /%d", decor.WCSyncSpaceR),
						decor.Percentage(decor.WCSyncSpaceR),
					),
				)
			}

			reader = bar.ProxyReader(res.Body)
		}
		_, err = io.Copy(file, reader)
		if err != nil {
			return err
		}

		if !isTerminal() {
			entry.Info("download success")
		}

		return nil
	}, retry.Attempts(3), retry.OnRetry(func(n uint, err error) {
		if !isTerminal() {
			entry.WithError(err).WithField("attempt", n).Error("failed downloading image. retrying...")
		}
	}))
	if err != nil {
		if !isTerminal() {
			entry.WithError(err).Error("remove fail download")
		}
		errRemove := os.Remove(path)
		if errRemove != nil {
			entry.WithError(errRemove).Debug("failed to remove")
		}
		return err
	}
	return nil
}

func isTerminal() bool {
	fileinfo, _ := os.Stdout.Stat()
	return (fileinfo.Mode() & os.ModeCharDevice) != 0
}

func (r Repository) downloadListing(ctx context.Context, subreddit string, subconf confmodel.Subreddit) ([]models.DownloadMeta, error) {
	entry := pkg.EntryFromContext(ctx)
	result := []models.DownloadMeta{}
	url := fmt.Sprintf("https://reddit.com/r/%s/%s.json", subreddit, subconf.Sort)
	entry.Trace(url)
	var resp *http.Response
	select {
	case <-ctx.Done():
		logrus.Panic("masuuuk")
	default:
	}
	err := retry.Do(func() error {
		// ctx, done := context.WithTimeout(ctx, r.config.Download.Timeout.Duration)
		// defer done()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return retry.Unrecoverable(err)
		}
		req.Header.Set("User-Agent", "ridit")
		res, err := r.client.Do(req)
		if err != nil {
			entry.WithError(err).Error(err)
			return err
		}
		if res.StatusCode >= 500 {
			err := errors.New("reddit is unavailable")
			entry.WithField("status_code", res.StatusCode).Error(err)
			return retry.Unrecoverable(err)
		}

		if res.StatusCode >= 400 {
			err := errors.New("reddit returned status " + res.Status)
			entry.WithField("status_code", res.StatusCode).Error(err)
			return err
		}
		resp = res
		return nil
	}, retry.Attempts(3), retry.OnRetry(func(n uint, err error) {
		entry.WithError(err).
			WithField("attempt", n).
			Error("failed fetching listing information. retrying...")
	}))
	if err != nil {
		return result, err
	}

	var listing models.Listing

	err = json.NewDecoder(resp.Body).Decode(&listing)
	if err != nil {
		entry.WithError(err).WithField("subreddit", subreddit).Error(err)
		return result, err
	}
	return listing.IntoDownloadMetas(r.config), nil
}
