package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config save config
type Config struct {
	Debug    bool   `yaml:"debug"`
	Bind     string `yaml:"bind"`
	Admin    string `yaml:"admin"`
	Password string `yaml:"password"`

	Db struct {
		Type     string `yaml:"type"`
		Path     string `yaml:"path"`
		Addr     string `yaml:"addr"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	}
}

var savedConfig *Config

// Load config from config.yml
func Load() *Config {
	if savedConfig != nil {
		return savedConfig
	}

	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Panicln("[config] Load, fail to read config.yml: ", err)
	}

	savedConfig = &Config{}
	err = yaml.Unmarshal(data, savedConfig)
	if err != nil {
		log.Panicln("[config] Load, fail to parse config: ", err)
	}

	log.Printf("Info: [config] Load, load: %+v\n", savedConfig)

	return savedConfig
}
