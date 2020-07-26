package models

import (
	"testing"
)

func TestUpdateDiff(t *testing.T) {
	goldenNum := 50.0
	min := 5.0
	max := 20.0
	tests := []struct {
		submit  float64
		result1 float64
		result2 float64
	}{
		{60, min, max},
		{40, min, max},
		{54, 54 - goldenNum, max},
		{46, goldenNum - 46, max},
		{80, min, 80 - goldenNum},
		{20, min, goldenNum - 20},
	}

	for _, test := range tests {
		r1, r2 := updateDiff(min, max, test.submit, goldenNum)
		if r1 != test.result1 || r2 != test.result2 {
			t.Errorf("updateDiff(%v, %v, %v, %v) = %v, %v; expected: %v, %v", min, max, test.submit, goldenNum, r1, r2, test.result1, test.result2)
		}
	}
}

func TestCalculateScoreGet(t *testing.T) {
	userNum := 50
	goldenNum := 50.0
	min := 5.0
	max := 20.0
	tests := []struct {
		submit float64
		result int
	}{
		{60, 0},
		{40, 0},
		{goldenNum + min, userNum - 2},
		{goldenNum - min, userNum - 2},
		{goldenNum + max, -2},
		{goldenNum - max, -2},
	}

	for _, test := range tests {
		r := calculateScoreGet(min, max, test.submit, goldenNum, userNum)
		if r != test.result {
			t.Errorf("calculateScoreGet(%v, %v, %v, %v, %v) = %v; expected: %v", min, max, test.submit, goldenNum, userNum, r, test.result)
		}
	}
}
