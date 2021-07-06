package config

import (
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit-go/app/config/models"
	"github.com/tigorlazuardi/ridit-go/app/config/models/sort"
)

func defaultConfigValue() models.Config {
	path := filepath.Join(GetHomeFolder(), "Pictures", "ridit")
	config := models.Config{
		Download: models.Download{
			Path:           path,
			ConnectTimeout: models.Duration{Duration: time.Second * 5},
			Timeout:        models.Duration{Duration: time.Second * 30},
		},
		AspectRatio: models.AspectRatio{
			Enabled: true,
			Height:  9.0,
			Width:   16.0,
			Range:   0.5,
		},
		MinimumSize: models.MinimumSize{
			Enabled: true,
			Height:  1080,
			Width:   1920,
		},
		Subreddits: map[string]models.Subreddit{
			"wallpaper": {
				Sort: sort.New,
				NSFW: true,
			},
			"wallpapers": {
				Sort: sort.New,
				NSFW: true,
			},
		},
		Daemon: models.Daemon{
			Port:              10101,
			WallpaperChange:   true,
			WallpaperInterval: models.Duration{Duration: time.Minute * 10},
		},
	}
	return config
}

func GetHomeFolder() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

// Write default config file if config not found. Never ignore the returned value. Make sure to close the file.
func LoadConfigFile() (*os.File, bool, error) {
	loc := viper.GetString("configfile")
	f, err := os.Open(loc)
	if err != nil {
		w, err := os.Create(loc)
		if err != nil {
			return w, false, err
		}
		val, _ := toml.Marshal(defaultConfigValue())
		_, err = w.Write(val)
		if err != nil {
			logrus.WithError(err).
				WithField("location", loc).
				WithField("solution", "try cleanup the location or set correct file / folder permission. make sure the user has permission to folder location").
				Fatal("failed to write default value to config file")
		}
		return w, true, nil
	} else {
		return f, false, nil
	}
}
