package models

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

const (
	roomStatusDefault  = 0
	roomStatusDisabled = 1
)

const (
	// UsersName is the col name of Users
	UsersName = "Users"
	// RoomHistorysName is the col name of RoomHistorys
	RoomHistorysName = "RoomHistorys"
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

	Status int
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

// roomWorker must not be copy
type roomWorker struct {
	ch       chan int
	nextTime time.Time
	submit   sync.Map

	historyLock   sync.RWMutex
	savedHistorys atomic.Value // []RoomHistory

	usersLock  sync.Mutex
	savedUsers atomic.Value // []User
}

const (
	roomIntervalDefault = 10
)

var (
	roomWorkers sync.Map
)

func getRoomWorker(id uint) (worker *roomWorker) {
	if value, ok := roomWorkers.Load(id); ok {
		if worker, ok := value.(*roomWorker); ok {
			return worker
		}
	}
	return nil
}

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
			if result := Models.First(&r2, r.ID); result.Error != nil {
				zap.S().Errorf("*Room.Runner, load: %v", result.Error)
			} else {
				*r = r2
			}

			if ok {
				r.RoundNow++
			}
			if result := Models.Save(r); result.Error != nil {
				zap.S().Errorf("*Room.Runner, save: %v", result.Error)
			}
		}
	}
}

// RestartAllRooms restart all not disabled rooms
func RestartAllRooms() {
	go func() {
		rooms := []Room{}
		Models.Not("Status", roomStatusDisabled).Find(&rooms)
		for _, room := range rooms {
			if room.RoundNow >= room.RoundTotal {
				continue
			}
			zap.S().Infof("RestartAll restarting room: %v", room.String())
			room.Start()
			time.Sleep(time.Millisecond * 500)
		}
	}()
}

// Start the room
func (r *Room) Start() bool {
	Models.Model(r).Update("Status", roomStatusDefault)
	worker := &roomWorker{
		ch: make(chan int),
	}
	if _, ok := roomWorkers.LoadOrStore(r.ID, worker); ok {
		zap.S().Errorf("*Room.Start, room already open, ID: %v", r.ID)
		return false
	}

	zap.S().Infof("*Room.Start, room open, ID: %v", r.ID)
	r2 := *r
	go r2.Runner(worker)
	return true
}

// Stop the room
func (r *Room) Stop() bool {
	Models.Model(r).Update("Status", roomStatusDisabled)
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
	if result := Models.Model(r).Association(UsersName).Find(&users); result != nil {
		zap.S().Errorf("*Room.GetUsers, %v", result)
	}
	return
}

// GetUsersWithCache return room's users from room cache
func (r *Room) GetUsersWithCache() (users []User) {
	worker := getRoomWorker(r.ID)
	if worker == nil {
		return r.GetUsers()
	}
	if v, ok := worker.savedUsers.Load().([]User); ok {
		return v
	}
	worker.usersLock.Lock()
	defer worker.usersLock.Unlock()
	users = r.GetUsers()
	worker.savedUsers.Store(users)
	return
}

// GetHistory return room's history
func (r *Room) GetHistory() (history []RoomHistory) {
	worker := getRoomWorker(r.ID)
	if worker != nil {
		if historys, ok := worker.savedHistorys.Load().([]RoomHistory); ok && len(history) > 0 {
			return historys
		}
	}
	if worker != nil {
		worker.historyLock.RLock()
		defer worker.historyLock.RUnlock()
	}
	if result := Models.Model(r).Association(RoomHistorysName).Find(&history); result != nil {
		zap.S().Errorf("*Room.GetHistory, %v", result)
	} else if worker != nil {
		worker.savedHistorys.Store(history)
	}
	return
}

// RoomUntilNextTick return time until next tick
func RoomUntilNextTick(id uint) (time.Duration, error) {
	if worker := getRoomWorker(id); worker != nil && !worker.nextTime.IsZero() {
		return time.Until(worker.nextTime), nil
	}
	return 0, fmt.Errorf("room %v not found or stopped", id)
}
