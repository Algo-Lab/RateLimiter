package limit

import (
	"fmt"
	"math"
	"testing"
	"time"
	"github.com/Algo-Lab/RateLimiter/pkg/utils"
)

func TestQpsLimiter_TryAcquire(t *testing.T) {
	limiter, err := NewQPSLimiter(0, 1000)
	if err != nil {
		t.Errorf("%v", err)
	}
	res := limiter.TryAcquire()
	if res {
		t.Errorf("false")
	} else {
		t.Log("ok")
	}
}

func TestQpsLimiter_TryAcquire1(t *testing.T) {
	limiter, err := NewQPSLimiter(1, 1000)
	if err != nil {
		t.Errorf("%v", err)
	}

	total := 0
	success := 0
	ticker := utils.NewTicker(func() {
		total++
		res := limiter.TryAcquire()
		if res {
			success++
		}
	})
	ticker.Start(time.Millisecond * 100)
	time.Sleep(5 * time.Second)
	ticker.Stop()
	if math.Abs(float64(success-5)) > 1 {
		t.Errorf("false, success=%d", success)
	} else {
		t.Log("total = ", total)
		t.Log("success = ", success)
	}
}

func TestQpsLimiter_TryAcquire2(t *testing.T) {
	maxAllowsList := []int64{1, 20, 500}
	periodMsList := []int64{1000, 3000}
	intervals := []time.Duration{500, 10}
	sleepSec := time.Duration(3)

	fmt.Println("start")
	for _, maxAllow := range maxAllowsList {
		for _, periodMs := range periodMsList {
			for _, interval := range intervals {
				fmt.Printf("maxAllow=%d, periodMs=%d, interval=%d", maxAllow, periodMs, interval)
				limiter, _ := NewQPSLimiter(maxAllow, periodMs)

				total := 0
				success := 0
				ticker := utils.NewTicker(func() {
					total++
					res := limiter.TryAcquire()
					if res {
						success++
					}
				})
				ticker.Start(time.Millisecond * interval)
				time.Sleep(sleepSec * time.Second)
				ticker.Stop()
				threshold := math.Min(float64(maxAllow*1000*int64(sleepSec))/float64(periodMs), float64(total))
				fmt.Printf("total=%d, success=%d, threshold=%f", total, success, threshold)
				if math.Abs(float64(success)-threshold) > 1 {
					t.Errorf("false, success=%d", success)
				}
			}
		}
	}
	t.Log("end")
}

