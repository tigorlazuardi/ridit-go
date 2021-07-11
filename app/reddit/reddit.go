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
	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	confmodel "github.com/tigorlazuardi/ridit/app/config/models"
	"github.com/tigorlazuardi/ridit/app/reddit/models"
	"github.com/tigorlazuardi/ridit/pkg"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type Repository struct {
	client pkg.Doer
	config confmodel.Config
	bars   *mpb.Progress
	sem    chan struct{}
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
		sem:    make(chan struct{}, viper.GetUint("concurrency")),
	}
}

func (r Repository) Fetch(ctx context.Context) <-chan DownloadChan {
	lc := make(chan ListingChan, len(r.config.Subreddits))

	go func() {
		wg := sync.WaitGroup{}
		for subreddit, subconf := range r.config.Subreddits {
			wg.Add(1)
			select {
			case <-ctx.Done():
				return
			default:
			}
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
			select {
			case <-ctx.Done():
				return
			default:
			}
			if listing.Err != nil {
				entry.WithError(listing.Err).Error(listing.Err)
				continue
			}
			wg.Add(1)
			go func(ctx context.Context, listing ListingChan) {
				defer wg.Done()
				wgg := &sync.WaitGroup{}
				for _, meta := range listing.Downloads {
					select {
					case <-ctx.Done():
						return
					default:
					}
					wgg.Add(1)
					r.sem <- struct{}{}
					go func(ctx context.Context, meta models.DownloadMeta) {
						defer wgg.Done()
						defer func() {
							<-r.sem
						}()
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
	select {
	case <-ctx.Done():
		return errors.New("download canceled")
	default:
	}
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
	var i uint = 0
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
		// cut the https://
		if pkg.IsTerminal() {
			job := meta.URL[8:]
			if len(job) >= 27 {
				job = job[:27]
			}
			job = color.FgWhite.Render(job)
			sub := color.FgLightBlue.Render("[" + meta.SubredditName + "]")
			downloading := color.FgLightYellow.Render("downloading")
			retrying := color.FgLightRed.Render("retrying")
			done := color.FgGreen.Render("done")
			var bar *mpb.Bar
			leftDecor := mpb.PrependDecorators(
				decor.Name(sub, decor.WC{W: 23, C: decor.DidentRight}),
				decor.Name(job, decor.WC{W: 28, C: decor.DidentRight}),
				decor.OnComplete(decor.Name(downloading, decor.WC{W: 11, C: decor.DidentRight}), done),
			)
			if i >= 1 {
				leftDecor = mpb.PrependDecorators(
					decor.Name("["+meta.SubredditName+"]", decor.WC{W: 23, C: decor.DidentRight}),
					decor.Name(job, decor.WC{W: 28, C: decor.DidentRight}),
					decor.OnComplete(decor.Name(retrying, decor.WC{W: 11, C: decor.DidentRight}), done),
				)
			}
			if res.ContentLength < 0 {
				bar = r.bars.AddSpinner(res.ContentLength, leftDecor)
			} else {
				bar = r.bars.AddBar(res.ContentLength,
					leftDecor,
					mpb.AppendDecorators(
						decor.CountersKiloByte("%d / %d", decor.WC{W: 14, C: decor.DidentRight}),
						decor.Percentage(decor.WC{W: 5, C: decor.DidentRight}),
					),
				)
			}
			reader = bar.ProxyReader(res.Body)
		} else {
			entry.WithField("subreddit", meta.SubredditName)
		}
		_, err = io.Copy(file, reader)
		file.Close()
		if err == context.Canceled {
			return retry.Unrecoverable(errors.New("download canceled"))
		}
		if err == context.DeadlineExceeded {
			return errors.New("download timeout")
		}
		if err != nil {
			return err
		}

		if !pkg.IsTerminal() {
			entry.Info("download success")
		}

		return nil
	}, retry.Attempts(3), retry.OnRetry(func(n uint, err error) {
		i = n
		if !pkg.IsTerminal() {
			entry.WithError(err).WithField("attempt", n).Error("failed downloading image. retrying...")
		}
	}))
	if err != nil {
		if !pkg.IsTerminal() {
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

func (r Repository) downloadListing(ctx context.Context, subreddit string, subconf confmodel.Subreddit) ([]models.DownloadMeta, error) {
	entry := pkg.EntryFromContext(ctx)
	result := []models.DownloadMeta{}
	url := fmt.Sprintf("https://reddit.com/r/%s/%s.json?limit=100", subreddit, subconf.Sort)
	entry.Trace(url)
	var resp *http.Response
	select {
	case <-ctx.Done():
		return nil, errors.New("download canceled")
	default:
	}
	err := retry.Do(func() error {
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

type Check struct {
	Name  string
	Exist bool
	Err   error
}

// Creates client if Doer is nil
func CheckSubredditExist(client pkg.Doer, ctx context.Context, subreddits []string) <-chan Check {
	cc := make(chan Check)
	go func() {
		if client == nil {
			client = http.DefaultClient
		}
		wg := sync.WaitGroup{}
		for _, v := range subreddits {
			wg.Add(1)
			go func(subreddit string) {
				defer wg.Done()
				url := "https://reddit.com/r/" + subreddit + ".json"
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				req.Header.Add("User-Agent", "ridit")
				res, err := client.Do(req)
				if err != nil {
					cc <- Check{Name: subreddit, Err: err}
					return
				}
				defer res.Body.Close()

				if res.StatusCode >= 300 {
					cc <- Check{Name: subreddit}
					return
				}
				cc <- Check{Name: subreddit, Exist: true}
			}(v)
		}
		wg.Wait()
		close(cc)
	}()
	return cc
}
