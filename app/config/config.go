package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit-go/app/config/models"
)

func Modify(f func(*models.Config)) error {
	config, err := Load()
	if err != nil {
		return err
	}
	f(&config)
	val, _ := toml.Marshal(config)
	configdir := configdir.LocalConfig("ridit", viper.GetString("profile"))
	path := filepath.Join(configdir, viper.GetString("profile")+".toml")
	return ioutil.WriteFile(path, val, 0777)
}

func Load() (config models.Config, err error) {
	f, _, err := LoadConfigFile()
	if err != nil {
		return
	}
	defer f.Close()
	tom, err := toml.LoadReader(f)
	if err != nil {
		return
	}
	err = tom.Unmarshal(&config)
	if err != nil {
		return
	}
	return
}
