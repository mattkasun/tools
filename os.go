package tools

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrNotDefined = errors.New("neither $XDG_DATA_HOME nor $HOME are defined") //nolint:revive
	ErrRelative   = errors.New("path in $XDG_DATA_HOME is relative")
)

// UserDataDir returns the default root directory to use for user-specific data.
// Users should create their own application-specific subdirectory within this
// one and use that.
//
// On Unix systems, it returns $XDG_DATA_HOME as specified by
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.local/share.
//
// If the location cannot be determined (for example, $HOME is not defined) or
// the path in $XDG_DATA_HOME is relative, then it will return an error.
func UserDataDir() (string, error) {
	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", ErrNotDefined
		}
		dir += "/.local/share"
	} else if !filepath.IsAbs(dir) {
		return "", ErrRelative
	}
	return dir, nil
}
