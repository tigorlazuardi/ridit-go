package models

import (
	"path"
	"path/filepath"
)

type DownloadMeta struct {
	URL             string `db:"url" json:"url"`
	SubredditName   string `db:"subreddit_name" json:"subreddit_name"`
	ImageHeight     uint   `db:"image_height" json:"image_height"`
	ImageWidth      uint   `db:"image_width" json:"image_width"`
	PostLink        string `db:"post_link" json:"post_link"`
	NSFW            bool   `db:"nsfw" json:"nsfw"`
	Filename        string `db:"filename" json:"filename"`
	Title           string `db:"title" json:"title"`
	Author          string `db:"author" json:"author"`
	Profile         string `db:"profile" json:"profile"`
	SuccessDownload bool   `db:"success_download" json:"success_download"`
}

func (dm DownloadMeta) GetFileLocation(baseLocation string) (string, error) {
	p := path.Join(baseLocation, dm.SubredditName, dm.Filename)
	return filepath.Abs(p)
}
