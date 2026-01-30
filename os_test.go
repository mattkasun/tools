package tools_test

import (
	"testing"

	"github.com/Kairum-Labs/should"
	"github.com/mattkasun/tools"
)

func TestUserDataDir(t *testing.T) {
	t.Run("XDG_DATA_HOME_set", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "/tmp/testing")
		dir, err := tools.UserDataDir()
		should.NotBeError(t, err)
		should.BeEqual(t, dir, "/tmp/testing")
	})
	t.Run("HomeSet", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		t.Setenv("HOME", "/tmp/home")
		dir, err := tools.UserDataDir()
		should.NotBeError(t, err)
		should.BeEqual(t, dir, "/tmp/home/.local/share")
	})
	t.Run("NotSet", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		t.Setenv("HOME", "")
		_, err := tools.UserDataDir()
		should.BeErrorIs(t, err, tools.ErrNotDefined)
	})
	t.Run("Relative", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "../home")
		_, err := tools.UserDataDir()
		should.BeErrorIs(t, err, tools.ErrRelative)
	})
}
