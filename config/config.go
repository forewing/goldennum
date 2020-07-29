package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/forewing/goldennum/utils"
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
	MaxConns int    `yaml:"max_conns"`
	MaxIdles int    `yaml:"max_idles"`
	ConnLife int    `yaml:"conn_life"`
}

const (
	defaultConfigPath = "conf/config.yml"
)

const (
	envUseEnvConfig = "USE_ENV_CONFIG"
	envDebug        = "DEBUG"
	envBind         = "BIND"
	envAdmin        = "ADMIN"
	envPassword     = "PASSWORD"
	envDbType       = "DB_TYPE"
	envDbPath       = "DB_PATH"
	envDbAddr       = "DB_ADDR"
	envDbUser       = "DB_USER"
	envDbPassword   = "DB_PASSWORD"
	envDbDbName     = "DB_NAME"
	envDbRedis      = "REDIS"
	envDbMaxConns   = "MAX_CONNS"
	envDbMaxIdles   = "MAX_IDLES"
	envDbConnLife   = "CONN_LIFE"
)

var (
	flagConf     = flag.String("conf", "", "Config file path. If set, will only use file config.")
	flagDebug    = flag.Bool("debug", false, "Set debug mode.")
	flagBind     = flag.String("bind", "localhost:8080", "Bind address.")
	flagAdmin    = flag.String("admin", "admin", "Admin username.")
	flagPassword = flag.String("password", "", "Admin password.")
	flagDbType   = flag.String("dbtype", "sqlite3", "[sqlite3, mysql]")
	flagDbPath   = flag.String("dbpath", "./sqlite3.db", "Path to sqlite3 database.")
	flagDbAddr   = flag.String("dbaddr", "localhost:3306", "Mysql server address.")
	flagDbUser   = flag.String("dbuser", "goldennum", "Database username.")
	flagDbPass   = flag.String("dbpass", "goldennum", "Database password.")
	flagDbName   = flag.String("dbname", "goldennum", "Database name.")
	flagRedis    = flag.String("redis", "", "Redis address, redis disabled if not set.")
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

	if len(*flagConf) > 0 {
		savedConfig.loadFromFile(*flagConf)
	} else {
		savedConfig.loadFromFlag()
	}
	savedConfig.completeFromEnv()

	savedConfig.complete()
	configLoaded = true
	log.Printf("Info: [config] Load: %+v\n", savedConfig)
	return &savedConfig
}

func (c *Config) loadFromFlag() {
	log.Println("[config] loadFromFlag")

	savedConfig = Config{
		Debug:    *flagDebug,
		Bind:     *flagBind,
		Admin:    *flagAdmin,
		Password: *flagPassword,
		Db: Db{
			Type:     *flagDbType,
			Path:     *flagDbPath,
			Addr:     *flagDbAddr,
			User:     *flagDbUser,
			Password: *flagDbPass,
			DbName:   *flagDbName,
			Redis:    *flagRedis,
		},
	}
}

func (c *Config) loadFromFile(path string) {
	log.Println("[config] loadFromFile: ", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("[config] loadFromConfig, fail to read %v: %v\n", path, err)
	}

	err = yaml.Unmarshal(data, &savedConfig)
	if err != nil {
		log.Panicln("[config] loadFromConfig, fail to parse config: ", err)
	}
}

func (c *Config) completeFromEnv() {
	if os.Getenv(envUseEnvConfig) != "true" {
		return
	}

	log.Println("[config] loadFromEnv")
	if os.Getenv(envDebug) == "true" {
		c.Debug = true
	}
	if s := os.Getenv(envAdmin); len(s) > 0 {
		c.Admin = s
	}
	if s := os.Getenv(envPassword); len(s) > 0 {
		c.Password = s
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
	if s := os.Getenv(envDbRedis); len(s) > 0 {
		c.Db.Redis = s
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
		c.Password, err = utils.RandomString(16)
		if err != nil {
			panic(err)
		}
	}
}
