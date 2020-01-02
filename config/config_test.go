package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_readConfigFile(t *testing.T) {
	c := &Config{}
	readConfigFile("../ph.example.conf", c)

	wantHostName := "host"
	if c.HostName != wantHostName {
		t.Errorf("HostName = %s, wants %s", c.HostName, wantHostName)
	}

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

	if c.Theme.Prompt != "powerline" {
		t.Errorf("theme.prompt = %s, wants powerline", c.Theme.Prompt)
	}

	if c.Scheme["simple"]["time"].Fg != "15" {
		t.Errorf("Scheme.simple.Time.Fg = %s, wants 15", c.Scheme["simple"]["time"].Fg)
	}
}

func Test_defaultConfig(t *testing.T) {
	RegisterDefaultScheme("fake", map[string]Color{
		"a": Color{Fg: "1"},
	})

	c := &Config{}
	c.defaultConfig()

	if c.Theme.Prompt != "powerline" {
		t.Errorf("theme.prompt = %s, wants powerline", c.Theme.Prompt)
	}

	if c.Scheme["fake"]["a"].Fg != "1" {
		t.Errorf("Scheme.fake.a.Fg = %s, wants 1", c.Scheme["fake"]["a"].Fg)
	}
}
