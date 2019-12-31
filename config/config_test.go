package config

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_readConfigFile(t *testing.T) {
	c := &Config{}
	readConfigFile("../ph.example.conf", c)

	wantVCS := map[string][]string{
		"git_skip":        []string{"/path/a"},
		"git_status_skip": []string{"/path/b", "/path/c"},
	}
	if diff := cmp.Diff(wantVCS, c.VCS); len(diff) != 0 {
		t.Errorf("VCS diff (-want, +got): %s", diff)
	}

	wantPathShorterns := map[string]string{
		"long/path": "shortpath",
	}
	if diff := cmp.Diff(wantPathShorterns, c.PathShorterns); len(diff) != 0 {
		t.Errorf("PathShorterns diff (-want, +got): %s", diff)
	}

	if c.Theme != "powerline" {
		t.Errorf("theme = %s, wants powerline", c.Theme)
	}

	if c.ColorScheme.Time.Fg != "15" {
		t.Errorf("ColorScheme.Time.Fg = %s, wants 15", c.ColorScheme.Time.Fg)
	}
}

func Test_defaultConfig(t *testing.T) {
	c := &Config{}
	c.defaultConfig()

	if c.Theme != "powerline" {
		t.Errorf("theme = %s, wants powerline", c.Theme)
	}

	if c.ColorScheme.Time.Fg != "15" {
		t.Errorf("ColorScheme.Time.Fg = %s, wants 15", c.ColorScheme.Time.Fg)
	}
}
