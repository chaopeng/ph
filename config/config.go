package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Color struct {
	Fg string `yaml:"fg"`
	Bg string `yaml:"bg"`
}

type ColorTheme struct {
	Tmux         string `yaml:"tmux"`
	Prompt       string `yaml:"prompt"`
	PromptInTmux string `yaml:"prompt_in_tmux"`
}

type Config struct {
	VCS           map[string][]string `yaml:"vcs"`
	PathShorterns map[string]string   `yaml:"path_shortern"`
	Theme         ColorTheme          `yaml:"theme"`
	Scheme        map[string]Color    `yaml:"scheme"`
}

func ReadConfig() *Config {
	conf := &Config{}
	defer conf.defaultConfig()

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		log.Fatalln("$HOME not set")
	}

	dotFile := homeDir + "/.ph"
	if _, err := os.Stat(dotFile); os.IsNotExist(err) {
		return conf
	}

	readConfigFile(dotFile, conf)

	return conf
}

func readConfigFile(path string, conf *Config) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, conf)
	if err != nil {
		log.Fatalln("yaml.Unmarshal failed")
	}
}

func setDefaultValue(field *string, defaultValue string) {
	if *field == "" {
		*field = defaultValue
	}
}

func (conf *Config) defaultConfig() {
	if conf.VCS == nil {
		conf.VCS = map[string][]string{}
	}

	if conf.PathShorterns == nil {
		conf.PathShorterns = map[string]string{}
	}

	setDefaultValue(&conf.Theme.Tmux, "powerline")
	setDefaultValue(&conf.Theme.Prompt, "powerline")
	setDefaultValue(&conf.Theme.PromptInTmux, "simple")

	defaultScheme := map[string]Color{
		// simple theme
		"simple/time": Color{
			Fg: "15",
		},
		"simple/good": Color{
			Fg: "2",
		},
		"simple/bad": Color{
			Fg: "1",
		},
		// powerline theme
		"power/status": Color{
			Fg: "1",
			Bg: "15",
		},
		"power/ssh": Color{
			Fg: "252",
			Bg: "240",
		},
		"power/os": Color{
			Fg: "15",
			Bg: "33",
		},
		"power/pwd": Color{
			Fg: "15",
			Bg: "240",
		},
		"power/pre_pwd": Color{
			Fg: "252",
		},
		"power/danger_zone": Color{
			Bg: "124",
		},
		"power/vcs_status_none": Color{
			Fg: "238",
			Bg: "3",
		},
		"power/vcs_status_clean": Color{
			Fg: "238",
			Bg: "2",
		},
		"power/vcs_status_dirty": Color{
			Fg: "15",
			Bg: "1",
		},
	}

	if conf.Scheme == nil {
		conf.Scheme = map[string]Color{}
	}

	for k, v := range defaultScheme {
		if _, ok := conf.Scheme[k]; !ok {
			conf.Scheme[k] = v
		}
	}
}
