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
	HostName      string                      `yaml:"hostname"`
	VCS           map[string]interface{}      `yaml:"vcs"`
	PathShorterns map[string]string           `yaml:"path_shortern"`
	Theme         ColorTheme                  `yaml:"theme"`
	Scheme        map[string]map[string]Color `yaml:"scheme"`
}

var (
	defaultScheme = map[string]map[string]Color{}
)

func RegisterDefaultScheme(theme string, scheme map[string]Color) {
	defaultScheme[theme] = scheme
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
		conf.VCS = map[string]interface{}{}
	}

	if conf.PathShorterns == nil {
		conf.PathShorterns = map[string]string{}
	}

	setDefaultValue(&conf.Theme.Tmux, "powerline")
	setDefaultValue(&conf.Theme.Prompt, "powerline")
	setDefaultValue(&conf.Theme.PromptInTmux, "simple")

	if conf.Scheme == nil {
		conf.Scheme = map[string]map[string]Color{}
	}

	for theme, scheme := range defaultScheme {
		if _, ok := conf.Scheme[theme]; !ok {
			conf.Scheme[theme] = scheme
		} else {
			for k, v := range scheme {
				if _, ok := conf.Scheme[theme][k]; !ok {
					conf.Scheme[theme][k] = v
				}
			}
		}
	}
}

func (c *Config) GetColor(theme, key string) Color {
	return c.Scheme[theme][key]
}
