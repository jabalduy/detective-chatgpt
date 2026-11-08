package game

import (
	"testing"
	"time"
)

func TestPerformActionAddsTime(t *testing.T) {
	start := time.Date(2026, 1, 1, 22, 0, 0, 0, time.UTC)

	clock := NewClock(start)

	state := &GameState{
		Clock: clock,
	}

	action := Action{
		Requirements: []Requirement{},
		Effects:      []Effect{},
		TimeCost:     30,
	}

	if !PerformAction(state, action) {
		t.Fatalf("не сработал performAction")

	}

	currentTime := clock.GetTime()

	if currentTime.Sub(start) != 30*time.Second {
		t.Fatalf("ожидалось 30 секунд, получено %v", currentTime.Sub(start))

	}

}
