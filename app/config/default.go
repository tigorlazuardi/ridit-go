package config

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit-go/app/config/models"
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
`

	path := path.Join(GetHomeFolder(), "Pictures", "ridit")
	val := fmt.Sprintf(defaultValue, path)
	return strings.TrimSpace(val)
}

func GetHomeFolder() string {
	usr, _ := user.Current()
	return strings.ReplaceAll(usr.HomeDir, string(os.PathSeparator), "/")
}

// Write default config file if config not found. Never ignore the returned value. Make sure to close the file.
func LoadConfigFile() (*os.File, bool) {
	loc := viper.GetString("configfile")
	f, err := os.Open(loc)
	if err != nil {
		w, err := os.Create(loc)
		if err != nil {
			logrus.WithError(err).WithField("location", loc).Fatal("failed to create config file")
		}
		_, err = w.WriteString(defaultConfigValue())
		if err != nil {
			logrus.WithError(err).WithField("location", loc).Fatal("failed to write default value to config file")
		}
		return w, true
	} else {
		return f, false
	}
}

func Load() (config models.Config) {
	f, _ := LoadConfigFile()
	defer f.Close()
	tom, err := toml.LoadReader(f)
	if err != nil {
		logrus.Fatal(err)
	}
	err = tom.Unmarshal(&config)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}
