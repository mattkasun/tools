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

// Get returns the cached value of the configuration struct if previously read
// else it reads struct from configuration path.
func Get[T any]() (*T, error) {
	if cached != nil {
		data, ok := cached.(*T)
		if !ok {
			return data, fmt.Errorf("%w: wanted %T but cached type is %T", errInteraceConversion, data, cached)
		}
		return data, nil
	}
	data, err := fromFile[T]()
	if err != nil {
		return nil, err
	}
	return data, err
}

// func fromFile populates the config struct from a yaml file in the XDG_CONFIG_HOME dir
// config file location is $XDG_CONFIG_HOME/executable name/config
// decoded structure is cached for faster subsequent retrievals.
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
	cached = &data
	return &data, nil
}
