package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGoWithRecover(t *testing.T) {
	debugIgnoreStdout = true
	panicStr := "test panic"
	panicHandler := func() {
		panic(panicStr)
	}
	output := ""
	recoverHandler := func(r interface{}) {
		output = fmt.Sprintf("%v", r)
	}
	GoWithRecover(panicHandler, recoverHandler)
	time.Sleep(time.Second) // wait panic goroutine
	if output != panicStr {
		t.Errorf("expected catch panic output, but got: %s", output)
	}
}

// recover handler panic, should not panic
func TestRecoverPanic(t *testing.T) {
	debugIgnoreStdout = true
	handler := func() {
		panic("1")
	}
	recoverHandler := func(r interface{}) {
		panic("2")
	}
	GoWithRecover(handler, recoverHandler)
}

// Example for how to recover with recover
type _run struct {
	count   int
	noPanic bool
}

func (r *_run) work() {
	GoWithRecover(r.exec, func(p interface{}) {
		r.work()
	})
}

func (r *_run) exec() {
	r.count++
	if r.count <= 2 {
		panic("panic")
	}
	r.noPanic = true
}

func TestGoWithRecoverAgain(t *testing.T) {
	debugIgnoreStdout = true
	r := &_run{}
	r.work()
	time.Sleep(time.Second)
	if !(r.noPanic && r.count == 3) {
		t.Errorf("panic handler is not restart expectedly, noPanic: %v, count: %d", r.noPanic, r.count)
	}
}

