package models

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

// Room hold user info
type Room struct {
	gorm.Model

	Users        []User
	RoomHistorys []RoomHistory

	Interval   int
	RountNow   int
	RoundTotal int
}

// RoomHistory holds room history
type RoomHistory struct {
	gorm.Model

	RoomID uint

	Round     int
	Goldennum float64
}

const (
	roomIntervalDefault = 10
)

var (
	roomChans     = make(map[uint]chan int)
	roomChansLock sync.Mutex
)

// Runner of room
func (r *Room) Runner(ch chan int) {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		default:
			if r.RountNow >= r.RoundTotal {
				r.Stop()
				return
			}
			if r.Interval <= 0 {
				log.Printf("Error: [models] *Room.Runner, room interval invalid: %v <= 0\n", r.Interval)
				r.Interval = roomIntervalDefault
			}
			time.Sleep(time.Duration(r.Interval) * time.Second)
			log.Printf("Info: [models] *Room.Runner, room tick, %v\n", r.String())
			// TODO: implement game logic and save model here

			var r2 Room
			if result := Db.First(&r2, r.ID); result.Error != nil {
				log.Printf("Error: [models] *Room.Runner, load: %v\n", result.Error)
			} else {
				*r = r2
			}
			r.RountNow++
			if result := Db.Save(r); result.Error != nil {
				log.Printf("Error: [models] *Room.Runner, save: %v\n", result.Error)
			}
		}
	}
}

// Start the room
func (r *Room) Start() bool {
	roomChansLock.Lock()
	defer roomChansLock.Unlock()

	if _, ok := roomChans[r.ID]; ok {
		log.Printf("Error: [models] *Room.Start, room already open, ID: %v\n", r.ID)
		return false
	}
	log.Printf("Info: [models] *Room.Start, room open, ID: %v\n", r.ID)
	ch := make(chan int)
	roomChans[r.ID] = ch
	go r.Runner(ch)
	return true
}

// Stop the room
func (r *Room) Stop() bool {
	roomChansLock.Lock()
	defer roomChansLock.Unlock()

	if ch, ok := roomChans[r.ID]; ok && ch != nil {
		close(ch)
		delete(roomChans, r.ID)
		log.Printf("Info: [models] *Room.Stop, room stop, %v\n", r.String())
		return true
	}
	log.Printf("Error: [models] *Room.Stop, room already closed, ID: %v\n", r.ID)
	return false
}

// String return formated room info
func (r *Room) String() string {
	return fmt.Sprintf("ID: %v, Round: %v/%v", r.ID, r.RountNow, r.RoundTotal)
}

// GetUsers return room's users
func (r *Room) GetUsers(secret bool) (users []User) {
	if result := Db.Model(r).Related(&users); result.Error != nil {
		log.Printf("Error: [models] *Room.GetUsers, %v\n", result.Error)
		return
	}

	if !secret {
		for i := range users {
			users[i].FilterInfo(false)
		}
	}

	return
}

// GetHistory return room's history
func (r *Room) GetHistory() (history []RoomHistory) {
	if result := Db.Model(r).Related(&history); result.Error != nil {
		log.Printf("Error: [models] *Room.GetHistory, %v\n", result.Error)
	}
	return
}
