package models

import (
	"time"

	"github.com/tigorlazuardi/ridit-go/app/config/models/sort"
)

type Config struct {
	Profile     string               `db:"profile" json:"profile" toml:"profile" yaml:"profile"`
	Download    Download             `db:"download" json:"download" toml:"download" yaml:"download"`
	AspectRatio AspectRatio          `db:"aspect_ratio" json:"aspect_ratio" toml:"aspect_ratio" yaml:"aspect_ratio"`
	MinimumSize MinimumSize          `db:"minimum_size" json:"minimum_size" toml:"minimum_size" yaml:"minimum_size"`
	Subreddits  map[string]Subreddit `db:"subreddits" json:"subreddits" toml:"subreddits" yaml:"subreddits"`
	Daemon      Daemon               `db:"daemon" json:"daemon" toml:"daemon" yaml:"daemon"`
}

func (c Config) GetAllSubreddits() []string {
	result := []string{}
	for k := range c.Subreddits {
		result = append(result, k)
	}
	return result
}

type Subreddit struct {
	Sort sort.Sort `db:"sort" json:"sort" toml:"sort" yaml:"sort" comment:"valid values = 'new, top, controversial, hot, rising'"`
	NSFW bool      `db:"nsfw" json:"nsfw" toml:"nsfw" yaml:"nsfw"`
}

type AspectRatio struct {
	Enabled bool    `db:"enabled" json:"enabled" toml:"enabled" yaml:"enabled"`
	Height  float32 `db:"height" json:"height" toml:"height" yaml:"height" comment:"value of 0 or negative will resort to default. default 9"`
	Width   float32 `db:"width" json:"width" toml:"width" yaml:"width" comment:"value of 0 or negative will resort to default. default 16"`
	Range   float32 `db:"range" json:"range" toml:"range" yaml:"range" comment:"aspect ratio range to be considered valid. ar value gained by dividing width with height. default 0.5"`
}

func (ar AspectRatio) Passed(height, width uint) bool {
	if !ar.Enabled {
		return true
	}
	if ar.Height <= 0 {
		ar.Height = 9
	}
	if ar.Width <= 0 {
		ar.Width = 16
	}
	configAr := ar.Width / ar.Height
	imageAr := float32(width) / float32(height)
	minRatio := configAr - ar.Range
	maxRatio := configAr + ar.Range
	return imageAr >= minRatio && imageAr <= maxRatio
}

type MinimumSize struct {
	Enabled bool `db:"enabled" json:"enabled" toml:"enabled" yaml:"enabled"`
	Height  uint `db:"height" json:"height" toml:"height" yaml:"height" comment:"minimum pixel size in height"`
	Width   uint `db:"width" json:"width" toml:"width" yaml:"width" comment:"minimum pixel size in width"`
}

func (ms MinimumSize) Passed(height, width uint) bool {
	if !ms.Enabled {
		return true
	}
	return height >= ms.Height && width >= ms.Width
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
