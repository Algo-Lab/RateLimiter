package utils

import (
	"sync/atomic"
	"time"
)

type Timer struct {
	innerTimer *time.Timer
	stopped    int32
}

func NewTimer(d time.Duration, callback func()) *Timer {
	return &Timer{
		innerTimer: time.AfterFunc(d, callback),
	}
}

// Stop stops the timer.
func (t *Timer) Stop() {
	if t == nil {
		return
	}
	if !atomic.CompareAndSwapInt32(&t.stopped, 0, 1) {
		return
	}

	t.innerTimer.Stop()
}
