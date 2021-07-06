package models

import (
	"time"

	"github.com/tigorlazuardi/ridit-go/app/config/models/sort"
)

type Config struct {
	Download    Download             `db:"download" json:"download" toml:"download" yaml:"download"`
	AspectRatio AspectRatio          `db:"aspect_ratio" json:"aspect_ratio" toml:"aspect_ratio" yaml:"aspect_ratio"`
	MinimumSize MinimumSize          `db:"minimum_size" json:"minimum_size" toml:"minimum_size" yaml:"minimum_size"`
	Subreddits  map[string]Subreddit `db:"subreddits" json:"subreddits" toml:"subreddits" yaml:"subreddits"`
	Daemon      Daemon               `db:"daemon" json:"daemon" toml:"daemon" yaml:"daemon"`
}

type Subreddit struct {
	Sort sort.Sort `db:"sort" json:"sort" toml:"sort" yaml:"sort" comment:"valid values = 'new, top, controversial, hot, rising'"`
	Nsfw bool      `db:"nsfw" json:"nsfw" toml:"nsfw" yaml:"nsfw"`
}

type AspectRatio struct {
	Enabled bool    `db:"enabled" json:"enabled" toml:"enabled" yaml:"enabled"`
	Height  uint    `db:"height" json:"height" toml:"height" yaml:"height"`
	Width   uint    `db:"width" json:"width" toml:"width" yaml:"width"`
	Range   float64 `db:"range" json:"range" toml:"range" yaml:"range"`
}

type MinimumSize struct {
	Enabled bool `db:"enabled" json:"enabled" toml:"enabled" yaml:"enabled"`
	Height  uint `db:"height" json:"height" toml:"height" yaml:"height"`
	Width   uint `db:"width" json:"width" toml:"width" yaml:"width"`
}

type Download struct {
	Path           string   `db:"path" json:"path" toml:"path" yaml:"path" comment:"download path"`
	ConnectTimeout Duration `db:"connect_timeout" json:"connect_timeout" toml:"connect_timeout" yaml:"connect_timeout" comment:"duration waiting for opening connection"`
	Timeout        Duration `db:"timeout" json:"timeout" toml:"timeout" yaml:"timeout" comment:"duration for download timeout. increase value on slow connection"`
}

type Daemon struct {
	Port              uint     `db:"port" json:"port" toml:"port" yaml:"port" comment:"daemon port. set between 1 to 65536. port 1 - 1023 requires root privilege"`
	WallpaperChange   bool     `db:"wallpaper_change" json:"wallpaper_change" toml:"wallpaper_change" yaml:"wallpaper_change" comment:"enable wallpaper changes"`
	WallpaperInterval Duration `db:"wallpaper_interval" json:"wallpaper_interval" toml:"wallpaper_interval" yaml:"wallpaper_interval" comment:"duration between wallpaper changes"`
}

type Duration struct{ time.Duration }

func (d *Duration) UnmarshalText(text []byte) error {
	dur, err := time.ParseDuration(string(text))
	d.Duration = dur
	return err
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}
