package config

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
	"github.com/tigorlazuardi/ridit/app/config/models"
)

func Modify(profile string, f func(*models.Config)) error {
	config, err := Load(profile)
	if err != nil {
		return err
	}
	f(&config)
	val, _ := toml.Marshal(config)
	path := GetConfigFilePath(profile)
	return ioutil.WriteFile(path, val, 0777)
}

func Load(profile string) (config models.Config, err error) {
	f, _, err := LoadConfigFile(profile)
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
