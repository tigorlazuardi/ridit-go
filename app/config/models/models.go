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
}

type Subreddit struct {
	Sort sort.Sort `db:"sort" json:"sort" toml:"sort" yaml:"sort"`
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
	Path           string   `db:"path" json:"path" toml:"path" yaml:"path"`
	ConnectTimeout Duration `db:"connect_timeout" json:"connect_timeout" toml:"connect_timeout" yaml:"connect_timeout"`
	Timeout        Duration `db:"timeout" json:"timeout" toml:"timeout" yaml:"timeout"`
}

type Duration struct{ time.Duration }

func (d *Duration) UnmarshalText(text []byte) error {
	dur, err := time.ParseDuration(string(text))
	d.Duration = dur
	return err
}
