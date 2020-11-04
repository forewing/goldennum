package models

import (
	"testing"
)

func testNewRoom(t *testing.T, interval int, roundTotal int) *Room {
	r := &Room{
		Interval:   interval,
		RoundTotal: roundTotal,
	}

	if err := Models.Create(r).Error; err != nil {
		t.Fatalf("fail to create Room, %v", err)
	}

	return r
}
