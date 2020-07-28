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
	RoundNow   int
	RoundTotal int
}

// RoomHistory holds room history
type RoomHistory struct {
	gorm.Model

	RoomID uint

	Round     int
	GoldenNum float64
}

type roomWorker struct {
	ch     chan int
	submit sync.Map
}

const (
	roomIntervalDefault = 10
)

var (
	roomWorkers sync.Map
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
			if r.RoundNow >= r.RoundTotal {
				r.Stop()
				return
			}
			if r.Interval <= 0 {
				log.Printf("Error: [models] *Room.Runner, room interval invalid: %v <= 0\n", r.Interval)
				r.Interval = roomIntervalDefault
			}
			time.Sleep(time.Duration(r.Interval) * time.Second)

			// log.Printf("Info: [models] *Room.Runner, room tick, %v\n", r.String())
			ok := r.tick()

			var r2 Room
			if result := Db.First(&r2, r.ID); result.Error != nil {
				log.Printf("Error: [models] *Room.Runner, load: %v\n", result.Error)
			} else {
				*r = r2
			}

			if ok {
				r.RoundNow++
			}
			if result := Db.Save(r); result.Error != nil {
				log.Printf("Error: [models] *Room.Runner, save: %v\n", result.Error)
			}
		}
	}
}

// Start the room
func (r *Room) Start() bool {
	worker := &roomWorker{
		ch: make(chan int),
	}
	if _, ok := roomWorkers.LoadOrStore(r.ID, worker); ok {
		log.Printf("Error: [models] *Room.Start, room already open, ID: %v\n", r.ID)
		return false
	}

	log.Printf("Info: [models] *Room.Start, room open, ID: %v\n", r.ID)
	go r.Runner(worker.ch)
	return true
}

// Stop the room
func (r *Room) Stop() bool {
	if value, ok := roomWorkers.Load(r.ID); ok {
		defer roomWorkers.Delete(r.ID)
		if worker, ok := value.(*roomWorker); ok && worker.ch != nil {
			close(worker.ch)
			log.Printf("Info: [models] *Room.Stop, room stop, %v\n", r.String())
			return true
		}
	}
	log.Printf("Error: [models] *Room.Stop, room already closed, ID: %v\n", r.ID)
	return false
}

// String return formated room info
func (r *Room) String() string {
	return fmt.Sprintf("ID: %v, Round: %v/%v", r.ID, r.RoundNow, r.RoundTotal)
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
