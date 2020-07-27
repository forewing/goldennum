package models

import (
	"github.com/jinzhu/gorm"
)

func init() {
	if Db != nil {
		return
	}

	var err error
	Db, err = gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic("fail to init database")
	}

	Db.AutoMigrate(&User{}, &UserHistory{}, &Room{}, &RoomHistory{})
}
