package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Color struct {
	Fg string `yaml:"Fg"`
	Bg string `yaml:"Bg"`
}

type Scheme struct {
	Time        Color `yaml:"Time"`
	StatusGood  Color `yaml:"StatusGood"`
	StatusBad   Color `yaml:"StatusBad"`
	Status      Color `yaml:"Status"`
	Ssh         Color `yaml:"Ssh"`
	Os          Color `yaml:"Os"`
	Pwd         Color `yaml:"Pwd"`
	PrePwd      Color `yaml:"PrePwd"`
	DangerZone  Color `yaml:"DangerZone"`
	StatusNone  Color `yaml:"StatusNone"`
	StatusClean Color `yaml:"StatusClean"`
	StatusDirty Color `yaml:"StatusDirty"`
}

type Config struct {
	Includes      []string          `yaml:"Includes"`
	PathShorterns map[string]string `yaml:"PathShorterns"`
	Theme         string            `yaml:"Theme"`
	ColorScheme   Scheme            `yaml:"Scheme"`
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

	b, err := ioutil.ReadFile(dotFile)
	if err != nil {
		return conf
	}

	err = yaml.Unmarshal(b, conf)
	if err != nil {
		log.Fatalln("yaml.Unmarshal failed")
	}

	return conf
}

func setDefaultValue(field *string, defaultValue string) {
	if *field == "" {
		*field = defaultValue
	}
}

func (conf *Config) defaultConfig() {
	setDefaultValue(&conf.Theme, "powerline")

	scheme := &conf.ColorScheme
	setDefaultValue(&scheme.Time.Fg, "15")
	setDefaultValue(&scheme.StatusGood.Fg, "2")
	setDefaultValue(&scheme.StatusBad.Fg, "1")
	setDefaultValue(&scheme.Status.Fg, "1")
	setDefaultValue(&scheme.Status.Bg, "15")
	setDefaultValue(&scheme.Ssh.Fg, "252")
	setDefaultValue(&scheme.Ssh.Bg, "240")
	setDefaultValue(&scheme.Os.Fg, "15")
	setDefaultValue(&scheme.Os.Bg, "33")
	setDefaultValue(&scheme.Pwd.Fg, "15")
	setDefaultValue(&scheme.Pwd.Bg, "240")
	setDefaultValue(&scheme.PrePwd.Fg, "252")
	setDefaultValue(&scheme.DangerZone.Bg, "124")
	setDefaultValue(&scheme.StatusNone.Fg, "238")
	setDefaultValue(&scheme.StatusNone.Bg, "3")
	setDefaultValue(&scheme.StatusClean.Fg, "238")
	setDefaultValue(&scheme.StatusClean.Bg, "2")
	setDefaultValue(&scheme.StatusDirty.Fg, "15")
	setDefaultValue(&scheme.StatusDirty.Bg, "1")
}

/*
# Default
Scheme:
 # simple theme
 Time:
  Fg: 15
 StatusGood:
  Fg: 2
 StatusBad:
  Fg: 1
 # powerline theme
 Status:
  Fg: 1
  Bg: 15
 Ssh:
  Fg: 252
  Bg: 240
 Os:
  Fg: 15
  Bg: 33
 Pwd:
  Fg: 15
  Bg: 240
 PrePwd:
  Fg: 252
 DangerZone:
  Bg: 124
 StatusNone:
  Fg: 238
  Bg: 3
 StatusClean:
  Fg: 238
  Bg: 2
 StatusDirty:
  Fg: 15
  Bg: 1
*/
