package models

import confmodel "github.com/tigorlazuardi/ridit-go/app/config/models"

type Listing struct {
	Data Data `json:"data"`
}

func (l Listing) IntoDownloadMetas(config confmodel.Config) []DownloadMeta {
	result := []DownloadMeta{}
	for _, children := range l.Data.Children {
		data := children.Data
		if data.IsVideo {
			continue
		}

	}
	return result
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
	Created   uint64   `json:"created"`
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

type Image struct {
	Source      Source       `json:"source"`
	Resolutions []Resolution `json:"resolutions"`
	ID          string       `json:"id"`
}

type Source struct {
	URL    string `json:"url"`
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}

type Resolution struct {
	URL    string `json:"url"`
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}
