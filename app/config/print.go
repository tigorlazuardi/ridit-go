package config

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

var ErrNotSupported = errors.New("unsupported output format")

// Valid values: json, toml, yaml.
func FormatConfig(profile, format string) (val []byte, err error) {
	config, err := Load(profile)
	if err != nil {
		return
	}
	switch strings.ToLower(format) {
	case "json":
		val, err = json.MarshalIndent(config, "", "    ")
	case "yaml":
		val, err = yaml.Marshal(config)
	case "toml":
		val, err = toml.Marshal(config)
	default:
		err = ErrNotSupported
	}
	if err != nil {
		return
	}
	return
}
