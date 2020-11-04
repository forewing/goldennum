package models

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/forewing/goldennum/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// Models is connected to database
	Models *gorm.DB

	dbRetryTime = time.Second * 1
	dbMaxRetry  = 120

	mySQLMaxOpenConns    = 150
	mySQLMaxIdlePrec     = 0.25
	mySQLConnMaxLifetime = time.Minute * 5

	sqliteMaxOpenConns    = 1
	sqliteMaxIdleConns    = 1
	sqliteConnMaxLifetime = time.Second * -1
)

const (
	// defaultDbConfig = "file::memory:?cache=shared"
	defaultDbConfig = "./sqlite3.db"
)

// Load init Db from config
func Load() {
	if Models != nil {
		zap.S().Panicf("load init twice")
	}

	conf := config.Load()

	var db gorm.Dialector
	switch conf.Db.Type {
	case "sqlite3":
		db = sqlite.Open(conf.Db.Path)
	case "mysql":
		waitDatabase(conf.Db.Addr)
		url := fmt.Sprintf("%v:%v@(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
			conf.Db.User, conf.Db.Password, conf.Db.Addr, conf.Db.DbName)
		db = mysql.Open(url)
	default:
		zap.S().Errorf("load db config not found or invalid, using: %v", defaultDbConfig)
		db = sqlite.Open(defaultDbConfig)
	}

	var err error
	Models, err = gorm.Open(db, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	setDBConfig(conf.Db)

	Models.AutoMigrate(&User{}, &UserHistory{}, &Room{}, &RoomHistory{})
	RestartAllRooms()
}

func setDBConfig(conf config.Db) {
	db, err := Models.DB()
	if err != nil {
		panic(err)
	}

	if conf.Type == "sqlite3" {
		db.SetMaxOpenConns(sqliteMaxOpenConns)
		db.SetMaxIdleConns(sqliteMaxIdleConns)
		db.SetConnMaxLifetime(sqliteConnMaxLifetime)
		return
	}

	db.SetMaxOpenConns(mySQLMaxOpenConns)
	db.SetMaxIdleConns(int(float64(mySQLMaxOpenConns) * mySQLMaxIdlePrec))
	db.SetConnMaxLifetime(mySQLConnMaxLifetime)
	if conf.MaxConns > 0 {
		db.SetMaxOpenConns(conf.MaxConns)
	}
	if conf.MaxIdles > 0 {
		db.SetMaxIdleConns(conf.MaxIdles)
	}
	if conf.ConnLife > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(conf.ConnLife))
	}
}

// Close Db
func Close() {
}

func waitDatabase(addr string) {
	dial := func() bool {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}
	for i := 0; i < dbMaxRetry; i++ {
		if dial() {
			return
		}
		zap.S().Warnf("database %v down, retrying %v/%v", addr, i, dbMaxRetry)
		time.Sleep(dbRetryTime)
	}
	zap.S().Panicf("database %v still down after %v retries", dbMaxRetry)
}
