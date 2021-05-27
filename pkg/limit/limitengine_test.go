package limit

import (
	"testing"
)

func TestNewLimitEngine(t *testing.T) {
	limitConfig := LimitConfig{
		LimitStrategy: QPSStrategy,
		MaxBurstRatio: 1.0,
		PeriodMs:      1000,
		MaxAllows:     10,
	}
	ruleConfig := &RuleConfig{
		LimitConfig: limitConfig,
	}

	limitEngine, err := NewLimitEngine(ruleConfig)
	if err != nil {
		t.Errorf("err=%s", err)
	}

	for i := 0; i < 10; i++ {
		ret := limitEngine.OverLimit()
		if ret {
			t.Errorf("false")
		}
	}
	for i := 0; i < 10; i++ {
		ret := limitEngine.OverLimit()
		if !ret {
			t.Errorf("false")
		}
	}

	limitConfig.MaxAllows = 0
	ruleConfig = &RuleConfig{
		LimitConfig: limitConfig,
	}
	limitEngine, err = NewLimitEngine(ruleConfig)
	ret := limitEngine.OverLimit()
	if !ret {
		t.Errorf("false")
	}
}
