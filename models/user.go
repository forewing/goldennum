package models

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

// User hold user info
type User struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	RoomID       uint
	UserHistorys []UserHistory `json:",omitempty"`

	Name  string
	Score int

	Hashed string `json:"-"`

	Submit1 float64 `gorm:"-" json:"-"`
	Submit2 float64 `gorm:"-" json:"-"`
}

// UserHistory holds user history
type UserHistory struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	UserID uint

	Round    int
	Score    int
	ScoreGet int
	Submit1  float64
	Submit2  float64
}

type userSubmit struct {
	s1 float64
	s2 float64
}

const (
	userNameMinLength = 1
	userNameMaxLength = 32

	userPassMinLength = 8
	userPassMaxLength = 32

	userSubmitMin     = 0.0
	userSubmitMax     = 100.0
	userSubmitInvalid = -1.0

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
		zap.S().Errorf("UserNew, bcrypt: %v", err)
		return nil, err
	}
	user := User{
		RoomID: roomid,
		Name:   name,
		Hashed: string(hashed),
		Score:  0,
	}
	return &user, nil
}

// GetHistory return user's history
func (u *User) GetHistory() (history []UserHistory) {
	if result := Db.Model(u).Related(&history); result.Error != nil {
		zap.S().Errorf("*User.GetHistory, %v", result.Error)
	}
	return
}

// Submit user input
func (u *User) Submit(s1, s2 float64) error {
	if value, ok := roomWorkers.Load(u.RoomID); ok {
		if worker, ok := value.(*roomWorker); ok {
			worker.submit.Store(u.ID, userSubmit{
				s1: s1,
				s2: s2,
			})
			return nil
		}
	}
	return fmt.Errorf("Room %v closed", u.RoomID)
}
