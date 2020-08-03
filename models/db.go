package models

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/forewing/goldennum/config"
	"github.com/jinzhu/gorm"

	// database drivers
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	// Db saves database
	Db *gorm.DB

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
	if Db != nil {
		zap.S().Panicf("load init twice")
	}

	conf := config.Load()

	var err error
	switch conf.Db.Type {
	case "sqlite3":
		Db, err = gorm.Open("sqlite3", conf.Db.Path)
	case "mysql":
		waitDatabase(conf.Db.Addr)
		url := fmt.Sprintf("%v:%v@(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
			conf.Db.User, conf.Db.Password, conf.Db.Addr, conf.Db.DbName)
		Db, err = gorm.Open("mysql", url)
	default:
		zap.S().Errorf("load db config not found or invalid, using: %v", defaultDbConfig)
		Db, err = gorm.Open("sqlite3", defaultDbConfig)
	}

	if err != nil {
		panic(err)
	}

	Db.DB().SetMaxOpenConns(mySQLMaxOpenConns)
	Db.DB().SetMaxIdleConns(int(float64(mySQLMaxOpenConns) * mySQLMaxIdlePrec))
	Db.DB().SetConnMaxLifetime(mySQLConnMaxLifetime)
	if conf.Db.MaxConns > 0 {
		Db.DB().SetMaxOpenConns(conf.Db.MaxConns)
	}
	if conf.Db.MaxIdles > 0 {
		Db.DB().SetMaxIdleConns(conf.Db.MaxIdles)
	}
	if conf.Db.ConnLife > 0 {
		Db.DB().SetConnMaxLifetime(time.Second * time.Duration(conf.Db.ConnLife))
	}
	if conf.Db.Type == "sqlite3" {
		Db.DB().SetMaxOpenConns(sqliteMaxOpenConns)
		Db.DB().SetMaxIdleConns(sqliteMaxIdleConns)
		Db.DB().SetConnMaxLifetime(sqliteConnMaxLifetime)
	}

	Db.AutoMigrate(&User{}, &UserHistory{}, &Room{}, &RoomHistory{})
}

// Close Db
func Close() {
	Db.Close()
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
