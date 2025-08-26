// Package config reads a yaml config file from the XDG_CONFIG_HOME and unmarshals it into a user supplied struct
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

var (
	cached                any //nolint:gochecknoglobals
	errInteraceConversion = errors.New("interface conversion")
)

// Get returns the configuration data for the supplied configuration struct type T, caching it after first retrieval.
func Get[T any]() (*T, error) {
	if cached == nil {
		data, err := fromFile[T]()
		if err != nil {
			return nil, err
		}
		cached = data
	}
	data, ok := cached.(*T)
	if !ok {
		return data, fmt.Errorf("%w: wanted %T but cached type is %T", errInteraceConversion, data, cached)
	}
	return data, nil
}

// func fromFile reads the yaml configuration file and unmarshals it into a struct of type T
// config file location is $XDG_CONFIG_HOME/executable name/config.
func fromFile[T any]() (*T, error) {
	progName := filepath.Base(os.Args[0])
	xdg, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		xdg = os.Getenv("HOME") + "/.config"
	}
	cfgfile := xdg + "/" + progName + "/config"
	bytes, err := os.ReadFile(cfgfile) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("read config file %w", err)
	}
	var data T
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("unmarshal %w", err)
	}
	return &data, nil
}
