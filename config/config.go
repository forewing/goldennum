package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/forewing/goldennum/version"
	"go.uber.org/zap"

	"gopkg.in/yaml.v2"
)

// Config save config
type Config struct {
	Debug    bool   `yaml:"debug"`
	Bind     string `yaml:"bind"`
	Admin    string `yaml:"admin"`
	Password string `yaml:"password"`
	BaseURL  string `yaml:"base_url"`

	Db Db
}

// Db saves db config
type Db struct {
	Type     string `yaml:"type"`
	Path     string `yaml:"path"`
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
	MaxConns int    `yaml:"max_conns"`
	MaxIdles int    `yaml:"max_idles"`
	ConnLife int    `yaml:"conn_life"`
}

const (
	envUseEnvConfig = "USE_ENV_CONFIG"
	envDebug        = "DEBUG"
	envBind         = "BIND"
	envAdmin        = "ADMIN"
	envPassword     = "PASSWORD"
	envBaseURL      = "BASE_URL"
	envDbType       = "DB_TYPE"
	envDbPath       = "DB_PATH"
	envDbAddr       = "DB_ADDR"
	envDbUser       = "DB_USER"
	envDbPassword   = "DB_PASSWORD"
	envDbDbName     = "DB_NAME"
	envDbMaxConns   = "MAX_CONNS"
	envDbMaxIdles   = "MAX_IDLES"
	envDbConnLife   = "CONN_LIFE"
)

var (
	flagConf     = flag.String("conf", "", "Config file path. If set, load config from file instead.")
	flagDebug    = flag.Bool("debug", false, "Set debug mode.")
	flagBind     = flag.String("bind", "0.0.0.0:8080", "Bind address.")
	flagAdmin    = flag.String("admin", "admin", "Admin username.")
	flagPassword = flag.String("password", "", "Admin password. Random if empty.")
	flagBaseURL  = flag.String("base-url", "", "Base URL. If you are using reverse proxy to redirect \"//PUBLIC_HOST/PREFIX/uri\" to \"//REAL_HOST/url\", it should be set to \"/PREFIX\"")
	flagDbType   = flag.String("db-type", "sqlite3", "[sqlite3, mysql]")
	flagDbPath   = flag.String("db-path", "./sqlite3.db", "Path to sqlite3 database.")
	flagDbAddr   = flag.String("db-addr", "localhost:3306", "Mysql server address.")
	flagDbUser   = flag.String("db-user", "goldennum", "Database username.")
	flagDbPass   = flag.String("db-pass", "goldennum", "Database password.")
	flagDbName   = flag.String("db-name", "goldennum", "Database name.")

	flagVersion = flag.Bool("version", false, "Display versions.")
)

var (
	configLoaded = false
	savedConfig  Config
)

// Load config from config.yml
func Load() *Config {
	if configLoaded {
		return &savedConfig
	}

	flag.Parse()
	if *flagVersion {
		version.Display()
		os.Exit(0)
	}

	if len(*flagConf) > 0 {
		savedConfig.loadFromFile(*flagConf)
	} else {
		savedConfig.loadFromFlag()
	}
	savedConfig.completeFromEnv()

	savedConfig.complete()
	configLoaded = true
	zap.S().Infof("loaded: %+v", savedConfig)
	return &savedConfig
}

func (c *Config) loadFromFlag() {
	zap.S().Debugf("load from flag")

	savedConfig = Config{
		Debug:    *flagDebug,
		Bind:     *flagBind,
		Admin:    *flagAdmin,
		Password: *flagPassword,
		BaseURL:  *flagBaseURL,
		Db: Db{
			Type:     *flagDbType,
			Path:     *flagDbPath,
			Addr:     *flagDbAddr,
			User:     *flagDbUser,
			Password: *flagDbPass,
			DbName:   *flagDbName,
		},
	}
}

func (c *Config) loadFromFile(path string) {
	zap.S().Debugf("load from file: %v", path)

	data, err := os.ReadFile(path)
	if err != nil {
		zap.S().Panicf("load from config, fail to read %v: %v", path, err)
	}

	err = yaml.Unmarshal(data, &savedConfig)
	if err != nil {
		zap.S().Panicf("load from config, fail to parse config: %v", err)
	}
}

func (c *Config) completeFromEnv() {
	if os.Getenv(envUseEnvConfig) != "true" {
		return
	}

	zap.S().Debugf("load from environment variable")
	if os.Getenv(envDebug) == "true" {
		c.Debug = true
	}
	if s := os.Getenv(envBind); len(s) > 0 {
		c.Bind = s
	}
	if s := os.Getenv(envAdmin); len(s) > 0 {
		c.Admin = s
	}
	if s := os.Getenv(envPassword); len(s) > 0 {
		c.Password = s
	}
	if s := os.Getenv(envBaseURL); len(s) > 0 {
		c.BaseURL = s
	}
	if s := os.Getenv(envDbType); len(s) > 0 {
		c.Db.Type = s
	}
	if s := os.Getenv(envDbPath); len(s) > 0 {
		c.Db.Path = s
	}
	if s := os.Getenv(envDbAddr); len(s) > 0 {
		c.Db.Addr = s
	}
	if s := os.Getenv(envDbUser); len(s) > 0 {
		c.Db.User = s
	}
	if s := os.Getenv(envDbPassword); len(s) > 0 {
		c.Db.Password = s
	}
	if s := os.Getenv(envDbDbName); len(s) > 0 {
		c.Db.DbName = s
	}
	if n, err := strconv.ParseInt(os.Getenv(envDbMaxConns), 10, 64); err == nil {
		c.Db.MaxConns = int(n)
	}
	if n, err := strconv.ParseInt(os.Getenv(envDbMaxIdles), 10, 64); err == nil {
		c.Db.MaxIdles = int(n)
	}
	if n, err := strconv.ParseInt(os.Getenv(envDbConnLife), 10, 64); err == nil {
		c.Db.ConnLife = int(n)
	}
}

func (c *Config) complete() {
	if len(c.Password) == 0 {
		var err error
		c.Password, err = randomString(16)
		if err != nil {
			panic(err)
		}
	}
}
