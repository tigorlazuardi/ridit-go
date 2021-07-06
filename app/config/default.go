package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func defaultConfigValue() string {
	const defaultValue = `
[download]
path = "%s"
connect_timeout = "5s"
timeout = "30s"

[aspect_ratio]
enabled = true
height = 16
width = 9
range = 0.5

[minimum_size]
enabled = true
height = 1080
width = 1920

[subreddits]
[subreddits.wallpaper]
sort = "new"
nsfw = true

[subreddits.wallpapers]
sort = "new"
nsfw = true

[daemon]
port = 10101
wallpaper_change = true
wallpaper_interval = "10m"
`

	path := filepath.Join(GetHomeFolder(), "Pictures", "ridit")
	val := fmt.Sprintf(defaultValue, path)
	return strings.TrimSpace(val)
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
		_, err = w.WriteString(defaultConfigValue())
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
