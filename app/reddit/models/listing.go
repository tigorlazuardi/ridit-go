package models

import (
	"errors"
	"strings"

	confmodel "github.com/tigorlazuardi/ridit/app/config/models"
)

type Listing struct {
	Data Data `json:"data"`
}

func (l Listing) IntoDownloadMetas(config confmodel.Config) []DownloadMeta {
	result := []DownloadMeta{}
	for _, children := range l.Data.Children {
		data := children.Data
		sub := data.Subreddit
		if data.IsVideo {
			continue
		}

		if data.Over18 && !config.Subreddits[sub].NSFW {
			continue
		}

		if data.Preview == nil {
			continue
		}

		height, width, err := data.Preview.GetImageSize()
		if err != nil {
			continue
		}
		if !config.AspectRatio.Passed(height, width) {
			continue
		}

		if !config.MinimumSize.Passed(height, width) {
			continue
		}

		filename := getFilenameFromURL(data.URL)
		if !strings.HasSuffix(filename, ".png") && !strings.HasSuffix(filename, ".jpg") && !strings.HasSuffix(filename, ".jpeg") {
			continue
		}
		meta := DownloadMeta{
			SubredditName:   data.Subreddit,
			ImageHeight:     height,
			ImageWidth:      width,
			PostLink:        "https://reddit.com" + data.Permalink,
			URL:             data.URL,
			NSFW:            data.Over18,
			Title:           data.Title,
			Author:          data.Author,
			Filename:        filename,
			SuccessDownload: false,
		}
		result = append(result, meta)
	}
	return result
}

func getFilenameFromURL(url string) string {
	split := strings.Split(url, "/")
	last := split[len(split)-1]
	return strings.Split(last, "?")[0]
}

type Data struct {
	Modhash  string     `json:"modhash"`
	Dist     int64      `json:"dist"`
	After    string     `json:"after"`
	Children []Children `json:"children"`
}

type Children struct {
	Data ChildrenData `json:"data"`
}

type ChildrenData struct {
	Subreddit string   `json:"subreddit"`
	Title     string   `json:"title"`
	PostHint  *string  `json:"post_hint"`
	Created   float64  `json:"created"`
	Over18    bool     `json:"over_18"`
	Preview   *Preview `json:"preview"`
	ID        string   `json:"id"`
	Author    string   `json:"author"`
	Permalink string   `json:"permalink"`
	Sticked   bool     `json:"sticked"`
	URL       string   `json:"url"`
	IsVideo   bool     `json:"is_video"`
	IsGallery bool     `json:"is_gallery"`
}

type Preview struct {
	Enabled bool    `json:"enabled"`
	Images  []Image `json:"images"`
}

func (p Preview) GetImageSize() (height, width uint, err error) {
	if len(p.Images) == 0 {
		return height, width, errors.New("empty image preview list")
	}
	img := p.Images[0]
	return img.Source.Height, img.Source.Width, nil
}

type Image struct {
	Source      Source       `json:"source"`
	Resolutions []Resolution `json:"resolutions"`
	ID          string       `json:"id"`
}

type Source struct {
	URL    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}

type Resolution struct {
	URL    string `json:"url"`
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}
