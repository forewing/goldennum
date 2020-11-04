package models

import (
	"math"
	"strconv"

	"go.uber.org/zap"
)

const (
	tickMinUserCount = 3
	goldenNumRatio   = 0.618
)

func updateDiff(min, max, submit, goldenNum float64) (float64, float64) {
	diff := math.Abs(submit - goldenNum)
	if diff > max {
		max = diff
	}
	if diff < min {
		min = diff
	}
	return min, max
}

func calculateScoreGet(min, max, submit, goldenNum float64, userNum int) int {
	diff := math.Abs(submit - goldenNum)
	if diff == min {
		return userNum - 2
	}
	if diff == max {
		return -2
	}
	return 0
}

func (r *Room) tick() bool {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("*Room.tick recover: %v", r)
		}
	}()
	workerValue, ok := roomWorkers.Load(r.ID)
	if !ok {
		return false
	}
	var worker *roomWorker
	if worker, ok = workerValue.(*roomWorker); !ok {
		return false
	}

	userAll := r.GetUsers()
	var users []*User
	for i := range userAll {
		userAll[i].Submit1 = userSubmitInvalid
		userAll[i].Submit2 = userSubmitInvalid
		if submitValue, ok := worker.submit.Load(userAll[i].ID); ok {
			if submit, ok := submitValue.(userSubmit); ok {
				userAll[i].Submit1 = submit.s1
				userAll[i].Submit2 = submit.s2
			}
		}
		if UserSubmitValidate(userAll[i].Submit1) && UserSubmitValidate(userAll[i].Submit2) {
			users = append(users, &userAll[i])
		}
	}

	userNum := len(users)
	if userNum < tickMinUserCount {
		zap.S().Warnf("*Room.tick, no enough valid users, %v, len: %v", r.String(), userNum)
		return false
	}

	goldenNum := 0.0
	for _, user := range users {
		goldenNum += user.Submit1
		goldenNum += user.Submit2
	}
	goldenNum /= float64(2 * userNum)
	goldenNum *= goldenNumRatio

	minDiff := 100.0
	maxDiff := 0.0
	for _, user := range users {
		minDiff, maxDiff = updateDiff(minDiff, maxDiff, user.Submit1, goldenNum)
		minDiff, maxDiff = updateDiff(minDiff, maxDiff, user.Submit2, goldenNum)
	}

	userHistorys := []*UserHistory{}
	for _, user := range users {
		result := calculateScoreGet(minDiff, maxDiff, user.Submit1, goldenNum, userNum)
		if result2 := calculateScoreGet(minDiff, maxDiff, user.Submit2, goldenNum, userNum); result2 != result {
			result += result2
		}

		user.Score += result

		history := &UserHistory{
			UserID:   user.ID,
			Round:    r.RoundNow,
			ScoreGet: result,
			Score:    user.Score,
			Submit1:  user.Submit1,
			Submit2:  user.Submit2,
		}
		userHistorys = append(userHistorys, history)

		worker.submit.Store(user.ID, userSubmit{
			s1: userSubmitInvalid,
			s2: userSubmitInvalid,
		})
	}

	for _, user := range users {
		if err := Models.Save(user).Error; err != nil {
			zap.S().Errorf("*Room.tick, fail to save user: %v", user.String())
		}
	}

	for _, history := range userHistorys {
		if err := Models.Save(history).Error; err != nil {
			zap.S().Errorf("*Room.tick, fail to save userHistory, %+v", *history)
			continue
		}
		userHistoryCache.Delete(strconv.Itoa(int(history.UserID)))
	}

	savedUsers := []User{}
	worker.usersLock.Lock()
	defer worker.usersLock.Unlock()

	if err := Models.Model(r).Association(UsersName).Find(&savedUsers).Error; err == nil {
		worker.savedUsers.Store(savedUsers)
	} else {
		zap.S().Errorf("*Room.tick, refresh worker users cache failed, %v", err)
	}

	roomHistory := RoomHistory{
		RoomID:    r.ID,
		Round:     r.RoundNow,
		GoldenNum: goldenNum,
	}
	worker.historyLock.Lock()
	defer worker.historyLock.Unlock()
	worker.savedHistorys.Store([]RoomHistory{})

	if err := Models.Save(&roomHistory).Error; err != nil {
		zap.S().Errorf("*Room.tick, fail to save roomHistory, %+v", roomHistory)
	}

	zap.S().Infof("*Room.tick, success, %v, len: %v, goldenNum: %v, minDiff: %v, maxDiff: %v",
		r.String(), userNum, goldenNum, minDiff, maxDiff)

	return true
}
