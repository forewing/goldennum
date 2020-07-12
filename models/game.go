package models

import (
	"log"
	"math"
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
	userAll := r.GetUsers(true)
	var users []*User
	for i := range userAll {
		if UserSubmitValidate(userAll[i].Submit1) && UserSubmitValidate(userAll[i].Submit2) {
			users = append(users, &userAll[i])
		}
	}

	userNum := len(users)
	if userNum < tickMinUserCount {
		log.Printf("Info: [models] *Room.tick, no enough valid users, %v, len: %v\n", r.String(), userNum)
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
		result += calculateScoreGet(minDiff, maxDiff, user.Submit2, goldenNum, userNum)

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

		user.Submit1 = -1
		user.Submit2 = -1
	}

	for _, user := range users {
		if err := Db.Save(user).Error; err != nil {
			log.Printf("Error: [models] *Room.tick, fail to save user, %v\n", user.String())
		}
	}

	for _, history := range userHistorys {
		if err := Db.Save(history).Error; err != nil {
			log.Printf("Error: [models] *Room.tick, fail to save userHistory, %+v\n", *history)
		}
	}

	roomHistory := RoomHistory{
		RoomID:    r.ID,
		Round:     r.RoundNow,
		GoldenNum: goldenNum,
	}

	if err := Db.Save(&roomHistory).Error; err != nil {
		log.Printf("Error: [models] *Room.tick, fail to save roomHistory, %+v\n", roomHistory)
	}

	log.Printf("Info: [models] *Room.tick, room tick success, %v, len: %v, goldenNum: %v, minDiff: %v, maxDiff: %v",
		r.String(), userNum, goldenNum, minDiff, maxDiff)

	return true
}
