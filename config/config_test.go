package config //nolint:testpackage

import (
	"os"
	"path/filepath"
	"testing"
)

// test struct.
type testConfig struct {
	Name  string `yaml:"name"`
	Count int    `yaml:"count"`
}

// reset global cache between tests.
func resetCache() {
	cached = nil
}

// helper: write temp config file.
func writeTempConfig(t *testing.T, progName string, content string) string {
	t.Helper()

	dir := t.TempDir()
	err := os.MkdirAll(filepath.Join(dir, progName), 0o750)
	if err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(dir, progName, "config")
	err = os.WriteFile(cfgPath, []byte(content), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestFromFile_Success(t *testing.T) {
	defer resetCache()

	progName := "myprog"
	cfgContent := "name: test\ncount: 42\n"
	xdgDir := writeTempConfig(t, progName, cfgContent)

	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	origArgs := os.Args
	os.Args = []string{progName}
	defer func() { os.Args = origArgs }()

	cfg, err := Get[testConfig]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "test" || cfg.Count != 42 {
		t.Errorf("unexpected config values: %+v", cfg)
	}
}

func TestFromFile_FileMissing(t *testing.T) {
	defer resetCache()

	progName := "noexist"
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	origArgs := os.Args
	os.Args = []string{progName}
	defer func() { os.Args = origArgs }()

	_, err := Get[testConfig]()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFromFile_InvalidYAML(t *testing.T) {
	defer resetCache()

	progName := "badprog"
	xdgDir := writeTempConfig(t, progName, "not: [valid\n") // malformed YAML

	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	origArgs := os.Args
	os.Args = []string{progName}
	defer func() { os.Args = origArgs }()

	_, err := Get[testConfig]()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGet_CacheReuse(t *testing.T) {
	defer resetCache()

	progName := "cachedprog"
	cfgContent := "name: cached\ncount: 99\n"
	xdgDir := writeTempConfig(t, progName, cfgContent)

	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	origArgs := os.Args
	os.Args = []string{progName}
	defer func() { os.Args = origArgs }()

	first, err := Get[testConfig]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := Get[testConfig]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// should return same pointer now that we cache *T
	if first != second {
		t.Error("expected same cached pointer, got different")
	}
}

func TestGet_CacheTypeMismatch(t *testing.T) {
	defer resetCache()

	progName := "mismatch"
	cfgContent := "name: mismatch\ncount: 1\n"
	xdgDir := writeTempConfig(t, progName, cfgContent)

	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	origArgs := os.Args
	os.Args = []string{progName}
	defer func() { os.Args = origArgs }()

	// first, populate cache with *testConfig
	_, err := Get[testConfig]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// now request with a different type
	type otherConfig struct {
		Other string `yaml:"other"`
	}
	_, err = Get[otherConfig]()
	if err == nil {
		t.Fatal("expected type mismatch error, got nil")
	}
}

func TestFromFile_XDGConfigHomeUnset(t *testing.T) {
	defer resetCache()

	progName := "fallbackprog"
	cfgContent := "name: fallback\ncount: 7\n"

	// Create config under $HOME/.config/progName/config
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir) // make sure HOME points to temp dir
	old := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", old) //nolint:errcheck,usetesting
	if err := os.Unsetenv("XDG_CONFIG_HOME"); err != nil {
		t.Fatalf("faied to unset env var: %v", err)
	}

	cfgDir := filepath.Join(homeDir, ".config", progName)
	if err := os.MkdirAll(cfgDir, 0o750); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config")
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o600); err != nil {
		t.Fatal(err)
	}

	origArgs := os.Args
	os.Args = []string{progName}
	defer func() { os.Args = origArgs }()

	cfg, err := Get[testConfig]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "fallback" || cfg.Count != 7 {
		t.Errorf("unexpected config values: %+v", cfg)
	}
}
