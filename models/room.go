package models

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Room hold user info
type Room struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Users        []User
	RoomHistorys []RoomHistory `json:",omitempty"`

	Interval   int
	RoundNow   int
	RoundTotal int
}

// RoomHistory holds room history
type RoomHistory struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	RoomID uint

	Round     int
	GoldenNum float64
}

type roomWorker struct {
	ch       chan int
	nextTime time.Time
	submit   sync.Map
}

const (
	roomIntervalDefault = 10
)

var (
	roomWorkers sync.Map
)

// Runner of room
func (r *Room) Runner(worker *roomWorker) {
	for {
		select {
		case _, ok := <-worker.ch:
			if !ok {
				return
			}
		default:
			if r.RoundNow >= r.RoundTotal {
				r.Stop()
				return
			}
			if r.Interval <= 0 {
				zap.S().Errorf("*Room.Runner, room interval invalid: %v <= 0", r.Interval)
				r.Interval = roomIntervalDefault
			}
			duration := time.Duration(r.Interval) * time.Second
			worker.nextTime = time.Now().Add(duration)
			time.Sleep(duration)

			ok := r.tick()

			var r2 Room
			if result := Db.First(&r2, r.ID); result.Error != nil {
				zap.S().Errorf("*Room.Runner, load: %v", result.Error)
			} else {
				*r = r2
			}

			if ok {
				r.RoundNow++
			}
			if result := Db.Save(r); result.Error != nil {
				zap.S().Errorf("*Room.Runner, save: %v", result.Error)
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
		zap.S().Errorf("*Room.Start, room already open, ID: %v", r.ID)
		return false
	}

	zap.S().Infof("*Room.Start, room open, ID: %v", r.ID)
	go r.Runner(worker)
	return true
}

// Stop the room
func (r *Room) Stop() bool {
	if value, ok := roomWorkers.Load(r.ID); ok {
		defer roomWorkers.Delete(r.ID)
		if worker, ok := value.(*roomWorker); ok && worker.ch != nil {
			close(worker.ch)
			zap.S().Infof("*Room.Stop, room stop, %v", r.String())
			return true
		}
	}
	zap.S().Errorf("*Room.Stop, room already closed, ID: %v", r.ID)
	return false
}

// String return formated room info
func (r *Room) String() string {
	return fmt.Sprintf("ID: %v, Round: %v/%v", r.ID, r.RoundNow, r.RoundTotal)
}

// GetUsers return room's users
func (r *Room) GetUsers() (users []User) {
	if result := Db.Model(r).Related(&users); result.Error != nil {
		zap.S().Errorf("*Room.GetUsers, %v", result.Error)
	}
	return
}

// GetHistory return room's history
func (r *Room) GetHistory() (history []RoomHistory) {
	if result := Db.Model(r).Related(&history); result.Error != nil {
		zap.S().Errorf("*Room.GetHistory, %v", result.Error)
	}
	return
}

// RoomUntilNextTick return time until next tick
func RoomUntilNextTick(id uint) (time.Duration, error) {
	if value, ok := roomWorkers.Load(id); ok {
		if worker, ok := value.(*roomWorker); ok {
			if !worker.nextTime.IsZero() {
				return time.Until(worker.nextTime), nil
			}
		}
	}
	return 0, fmt.Errorf("room %v not found or stopped", id)
}
