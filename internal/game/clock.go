package game

import (
	"context"
	"sync"
	"time"
)

type Clock struct {
	currentTime time.Time
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	isPaused    bool
}

func (c *Clock) RunTime() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return

		case <-ticker.C:
			c.mu.Lock()

			if c.isPaused {
				c.mu.Unlock()
				continue
			}
			c.currentTime = c.currentTime.Add(1 * time.Second)
			c.mu.Unlock()
		}
	}
}

func NewClock(startTime time.Time) *Clock {
	TimeContext, TimeCancel := context.WithCancel(context.Background())

	clock := Clock{
		currentTime: startTime,
		mu:          sync.RWMutex{},
		ctx:         TimeContext,
		cancel:      TimeCancel,
	}

	return &clock
}

type StartTime struct {
	StartTime time.Time
}

func (c *Clock) GetTime() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.currentTime
}

func (c *Clock) StopTime() {
	c.cancel()
}

func (c *Clock) Pause() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isPaused = true
}

func (c *Clock) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isPaused = false
}

func (c *Clock) Add(duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.currentTime = c.currentTime.Add(duration)
}
