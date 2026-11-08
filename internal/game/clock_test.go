package game

import (
	"fmt"
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	clock := NewClock(time.Now())

	go clock.RunTime()

	time.Sleep(2 * time.Second)

	beforePause := clock.GetTime()

	clock.Pause()

	time.Sleep(3 * time.Second)

	afterPause := clock.GetTime()

	if beforePause != afterPause {
		t.Fatalf("lalala")
	}

	clock.Add(30 * time.Second)

	fmt.Println(clock.GetTime())

	clock.Resume()

	time.Sleep(3 * time.Second)

	fmt.Println(clock.GetTime())

	clock.StopTime()
}
