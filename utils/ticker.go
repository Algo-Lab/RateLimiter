package utils

import (
	"runtime/debug"
	"sync/atomic"
	"time"
)

// Ticker is a thread-safe reusable ticker
type Ticker struct {
	innerTicker *time.Ticker
	interval    time.Duration
	callback    func()
	stopChan    chan bool
	started     int32
	stopped     int32
}

func NewTicker(callback func()) *Ticker {
	return &Ticker{
		callback: callback,
		stopChan: make(chan bool, 1),
	}
}

// Start starts a ticker running if it is not started
func (t *Ticker) Start(interval time.Duration) {
	if !atomic.CompareAndSwapInt32(&t.started, 0, 1) {
		return
	}

	if t.innerTicker == nil {
		t.innerTicker = time.NewTicker(interval)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
			}
			t.Close()
			atomic.StoreInt32(&t.started, 0)
			atomic.StoreInt32(&t.stopped, 0)
		}()

		for {
			select {
			case <-t.innerTicker.C:
				t.callback()
			case <-t.stopChan:
				t.innerTicker.Stop()
				return
			}
		}
	}()

}

// Stop stops the ticker.
func (t *Ticker) Stop() {
	if !atomic.CompareAndSwapInt32(&t.stopped, 0, 1) {
		return
	}

	t.stopChan <- true
}

// Close closes the ticker.
func (t *Ticker) Close() {
	close(t.stopChan)
}
