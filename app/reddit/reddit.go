package reddit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	confmodel "github.com/tigorlazuardi/ridit-go/app/config/models"
	"github.com/tigorlazuardi/ridit-go/app/reddit/models"
	"github.com/tigorlazuardi/ridit-go/pkg"
	"github.com/vbauerster/mpb/v7"
)

type Repository struct {
	client pkg.Doer
	config confmodel.Config
	bars   *mpb.Progress
}

type ListingChan struct {
	Downloads []models.DownloadMeta
	Err       error
}

func NewRepository(client pkg.Doer, config confmodel.Config) Repository {
	return Repository{
		client: client,
		config: config,
		bars:   mpb.New(mpb.WithOutput(os.Stdout)),
	}
}

func (r Repository) GetListing(ctx context.Context) <-chan ListingChan {
	lc := make(chan ListingChan, len(r.config.Subreddits))

	go func() {
		wg := sync.WaitGroup{}
		for subreddit, subconf := range r.config.Subreddits {
			wg.Add(1)
			go func(ctx context.Context, subreddit string, subconf confmodel.Subreddit) {
				defer wg.Done()
				ctx, done := context.WithTimeout(ctx, r.config.Download.Timeout.Duration)
				meta, err := r.downloadListing(ctx, subreddit, subconf)
				lc <- ListingChan{
					Downloads: meta,
					Err:       err,
				}
				done()
			}(ctx, subreddit, subconf)
		}
		wg.Wait()
		close(lc)
	}()

	return lc
}

func (r Repository) downloadListing(ctx context.Context, subreddit string, subconf confmodel.Subreddit) ([]models.DownloadMeta, error) {
	entry := pkg.EntryFromContext(ctx)
	result := []models.DownloadMeta{}
	url := fmt.Sprintf("https://reddit.com/r/%s/%s.json", subreddit, subconf.Sort)
	entry.Debug(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("User-Agent", "ridit")
	res, err := r.client.Do(req)
	if err != nil {
		entry.WithError(err).Error(err)
		return result, err
	}
	if res.StatusCode >= 500 {
		err := errors.New("reddit is unavailable")
		entry.WithField("status_code", res.StatusCode).Error(err)
		return result, err
	}

	if res.StatusCode >= 400 {
		err := errors.New("reddit returned status " + res.Status)
		entry.WithField("status_code", res.StatusCode).Error(err)
		return result, err
	}
	var listing models.Listing
	err = json.NewDecoder(res.Body).Decode(&listing)
	if err != nil {
		entry.WithError(err).WithField("subreddit", subreddit).Error(err)
		return result, err
	}
	// entry.WithField("listing", listing).Debug("listign")

	return listing.IntoDownloadMetas(ctx, r.config), nil
}
