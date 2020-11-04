package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	if Models != nil {
		return
	}

	var err error
	Models, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("fail to init database")
	}

	Models.AutoMigrate(&User{}, &UserHistory{}, &Room{}, &RoomHistory{})
}
