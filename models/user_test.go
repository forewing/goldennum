package models

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func testNewUser(t *testing.T, roomID uint, score int, submit1, submit2 float64) *User {
	u := &User{
		RoomID:  roomID,
		Score:   score,
		Submit1: submit1,
		Submit2: submit2,
	}

	if err := Db.Create(u).Error; err != nil {
		t.Fatalf("fail to create User, %v", err)
	}

	return u
}

func testNewUserHistory(t *testing.T, userID uint, round, score, scoreGet int, submit1, submit2 float64) *UserHistory {
	uh := &UserHistory{
		UserID:   userID,
		Round:    round,
		Score:    score,
		ScoreGet: scoreGet,
		Submit1:  submit1,
		Submit2:  submit2,
	}

	if err := Db.Create(uh).Error; err != nil {
		t.Fatalf("fail to create UserHistory, %v", err)
	}

	return uh
}

func TestUserNameValidate(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"user", true},
		{"", false},
		{strings.Repeat("x", 33), false},
		{"你好", false},
	}

	for _, test := range tests {
		r := UserNameValidate(test.input)
		if r != test.valid {
			t.Errorf("UserNameValidate(%v) = %v; expected: %v", test.input, r, test.valid)
		}
	}
}

func TestUserPassValidate(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"password", true},
		{"pass", false},
		{strings.Repeat("x", 33), false},
	}

	for _, test := range tests {
		r := UserPassValidate(test.input)
		if r != test.valid {
			t.Errorf("UserPassValidate(%v) = %v; expected: %v", test.input, r, test.valid)
		}
	}
}

func TestUserSubmitValidate(t *testing.T) {
	tests := []struct {
		input float64
		valid bool
	}{
		{12.6, true},
		{-1, false},
		{100.5, false},
	}

	for _, test := range tests {
		r := UserSubmitValidate(test.input)
		if r != test.valid {
			t.Errorf("UserSubmitValidate(%v) = %v; expected: %v", test.input, r, test.valid)
		}
	}
}

func TestUserAuth(t *testing.T) {
	tests := []struct {
		pass   string
		hashed string
		result error
	}{
		{"password", "$2a$10$mLHY719u2kmxzE/g/zP2o.ePCPI8LNTHsc49CctV80ywuseBHesLW", nil},
		{strings.Repeat("#", 30), "$2a$10$N3JDoveoksSMNZD48cvecuOszhzsQL5P4hRkLmun/qFvs0nuU.mMq", nil},
		{"password", "$2a$10$N3JDoveoksSMNZD48cvecuOszhzsQL5P4hRkLmun/qFvs0nuU.mMq", errors.New("")},
		{"password", "invalid", errors.New("")},
	}

	for _, test := range tests {
		u := User{Hashed: test.hashed}
		r := u.Auth(test.pass)
		if (test.result == nil && r != nil) || (test.result != nil && r == nil) {
			t.Errorf("User{Hashed: %v}.Auth(%v) = %v; expected: %v", test.hashed, test.pass, r, test.result)
		}
	}
}

func TestUserString(t *testing.T) {
	tests := []struct {
		id     uint
		name   string
		result string
	}{
		{1, "user1", "ID: 1, Name: user1"},
		{7984353, "user7984353", "ID: 7984353, Name: user7984353"},
	}

	for _, test := range tests {
		u := User{
			ID:   test.id,
			Name: test.name,
		}
		r := u.String()
		if r != test.result {
			t.Errorf("User{ID: %v, Name: %v}.String() = %v; expected: %v", test.id, test.name, r, test.result)
		}
	}
}

func TestUserNew(t *testing.T) {
	expected := &User{
		RoomID: 1,
		Name:   "user1",
		Score:  0,
	}

	ret, _ := UserNew(1, "user1", "password")
	if expected.RoomID != ret.RoomID ||
		expected.Name != ret.Name ||
		bcrypt.CompareHashAndPassword([]byte(ret.Hashed), []byte("password")) != nil ||
		expected.Score != ret.Score {
		t.Errorf("UserNew(%v, %v, %v) = %+v; expected: %+v", 1, "user1", "password", ret, expected)
	}
}

func TestUserGetHistory(t *testing.T) {
	r := testNewRoom(t, 10, 10)
	u := testNewUser(t, r.ID, 10, 1.1, 2.2)
	uh := []*UserHistory{}
	uh = append(uh, testNewUserHistory(t, u.ID, 0, 10, 4, 1.1, 1.2))
	uh = append(uh, testNewUserHistory(t, u.ID, 1, 14, -2, 2.1, 2.2))
	uh = append(uh, testNewUserHistory(t, u.ID, 2, 12, 6, 3.1, 3.2))
	uh = append(uh, testNewUserHistory(t, u.ID, 3, 18, -2, 4.1, 4.2))

	uh2 := u.GetHistory()

	if len(uh) != len(uh2) {
		t.Errorf("len(u.GetHistory()) = %v; expected: %v", len(uh2), len(uh))
		return
	}

	uh2map := make(map[int]*UserHistory)
	for i := range uh2 {
		uh2 := &uh2[i]
		uh2map[uh2.Round] = uh2
	}

	for i := range uh {
		uh := uh[i]
		uh2 := uh2map[uh.Round]
		if uh.UserID != uh2.UserID ||
			uh.Score != uh2.Score ||
			uh.ScoreGet != uh2.ScoreGet ||
			uh.Submit1 != uh2.Submit1 ||
			uh.Submit2 != uh2.Submit2 ||
			uh.Round != uh2.Round {
			t.Errorf("u.GetHistory()[%v] = %+v; expected: %+v", i, uh2, uh)
		}
	}
}
