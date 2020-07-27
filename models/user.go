package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User hold user info
type User struct {
	gorm.Model

	RoomID       uint
	UserHistorys []UserHistory

	Hashed string

	Name    string
	Score   int
	Submit1 float64
	Submit2 float64
}

// UserHistory holds user history
type UserHistory struct {
	gorm.Model

	UserID uint

	Round    int
	Score    int
	ScoreGet int
	Submit1  float64
	Submit2  float64
}

const (
	userNameMinLength = 1
	userNameMaxLength = 32

	userPassMinLength = 8
	userPassMaxLength = 32

	userSubmitMin = 0.0
	userSubmitMax = 100.0

	uesrPassBcryptCost = bcrypt.DefaultCost
)

// UserNameValidate validate user name
func UserNameValidate(name string) bool {
	l := len(name)
	if l < userNameMinLength || l > userNameMaxLength {
		return false
	}
	return true
}

// UserPassValidate validate user password
func UserPassValidate(pass string) bool {
	l := len(pass)
	if l < userPassMinLength || l > userPassMaxLength {
		return false
	}
	return true
}

// UserSubmitValidate validate user submit
func UserSubmitValidate(submit float64) bool {
	if submit <= userSubmitMin || submit >= userSubmitMax {
		return false
	}
	return true
}

// Auth auth user
func (u *User) Auth(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Hashed), []byte(password))
}

// String return formatted user info
func (u *User) String() string {
	return fmt.Sprintf("ID: %v, Name: %v", u.ID, u.Name)
}

// UserNew build a new user
func UserNew(roomid uint, name, pass string) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), uesrPassBcryptCost)
	if err != nil {
		log.Printf("Error: [models] UserNew, bcrypt: %v\n", err)
		return nil, err
	}
	user := User{
		RoomID:  roomid,
		Name:    name,
		Hashed:  string(hashed),
		Score:   0,
		Submit1: -1,
		Submit2: -1,
	}
	return &user, nil
}

// FilterInfo delete secret info
func (u *User) FilterInfo(secret bool) {
	if !secret {
		u.Submit1 = -1
		u.Submit2 = -1
		u.Hashed = ""
	}
}

// GetHistory return user's history
func (u *User) GetHistory() (history []UserHistory) {
	if result := Db.Model(u).Related(&history); result.Error != nil {
		log.Printf("Error: [models] *User.GetHistory, %v\n", result.Error)
	}
	return
}
