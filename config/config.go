package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config save config
type Config struct {
	Debug    bool   `yaml:"debug"`
	Bind     string `yaml:"bind"`
	Admin    string `yaml:"admin"`
	Password string `yaml:"password"`

	Db Db
}

// Db saves db config
type Db struct {
	Type     string `yaml:"type"`
	Path     string `yaml:"path"`
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	Redis    string `yaml:"redis"`
}

const (
	defaultConfigPath = "conf/config.yml"

	envUseEnvConfig = "USE_ENV_CONFIG"
	envConfigPath   = "CONFIG_PATH"

	envDebug      = "DEBUG"
	envBind       = "BIND"
	envAdmin      = "ADMIN"
	envPassword   = "PASSWORD"
	envDbType     = "DB_TYPE"
	envDbPath     = "DB_PATH"
	envDbAddr     = "DB_ADDR"
	envDbUser     = "DB_USER"
	envDbPassword = "DB_PASSWORD"
	envDbDbName   = "DB_NAME"
	envDbRedis    = "DB_REDIS"
)

var savedConfig *Config

// Load config from config.yml
func Load() *Config {
	if savedConfig != nil {
		return savedConfig
	}

	if os.Getenv(envUseEnvConfig) == "true" {
		log.Println("Info: [config] Load from env")
		savedConfig = loadFromEnv()
	} else {
		path := os.Getenv(envConfigPath)
		if len(path) == 0 {
			path = defaultConfigPath
		}
		log.Println("Info: [config] Load from file: ", path)
		savedConfig = loadFromConfig(path)
	}

	log.Printf("Info: [config] Load, load: %+v\n", savedConfig)
	return savedConfig
}

func loadFromEnv() *Config {
	config := &Config{
		Debug:    false,
		Bind:     os.Getenv(envBind),
		Admin:    os.Getenv(envAdmin),
		Password: os.Getenv(envPassword),
		Db: Db{
			Type:     os.Getenv(envDbType),
			Path:     os.Getenv(envDbPath),
			Addr:     os.Getenv(envDbAddr),
			User:     os.Getenv(envDbUser),
			Password: os.Getenv(envDbPassword),
			DbName:   os.Getenv(envDbDbName),
			Redis:    os.Getenv(envDbRedis),
		},
	}

	if os.Getenv(envDebug) == "true" {
		config.Debug = true
	}

	return config
}

func loadFromConfig(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("[config] loadFromConfig, fail to read %v: %v\n", path, err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Panicln("[config] loadFromConfig, fail to parse config: ", err)
	}

	return config
}
