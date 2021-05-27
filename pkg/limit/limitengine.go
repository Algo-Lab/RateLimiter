package limit

import (
	"errors"
	logger "github.com/sirupsen/logrus"
)

type RuleConfig struct {
	Id      int    `json:"id"`
	Index   int    `json:"index"`
	Name    string `json:"name"`
	AppName string `json:"appName"`
	Enable  bool   `json:"enable"`
	RunMode string `json:"run_mode"`

	LimitConfig     LimitConfig      `json:"limit"`
	ResourceConfigs []ResourceConfig `json:"resources"`
}

// ComparisonCofig config
type ComparisonCofig struct {
	CompareType string `json:"compare_type"`
	Key         string `json:"key"`
	Value       string `json:"value"`
}

// ResourceConfig config
type ResourceConfig struct {
	Headers         []ComparisonCofig `json:"headers"`
	HeadersRelation string            `json:"headers_relation"`
	Params          []ComparisonCofig `json:"params"`
	ParamsRelation  string            `json:"params_relation"`
}

// LimitConfig config
type LimitConfig struct {
	LimitStrategy string  `json:"limit_strategy"`
	MaxBurstRatio float64 `json:"max_burst_ratio"`
	PeriodMs      int     `json:"period_ms"`
	MaxAllows     int     `json:"max_allows"`
}

// LimitEngine limit
type LimitEngine struct {
	RuleConfig *RuleConfig
	limiter    Limiter
}

// NewLimitEngine limit
func NewLimitEngine(ruleConfig *RuleConfig) (*LimitEngine, error) {
	l := &LimitEngine{
		RuleConfig: ruleConfig,
	}
	config := ruleConfig.LimitConfig
	if config.LimitStrategy == QPSStrategy {
		limiter, err := NewQPSLimiter(int64(config.MaxAllows), int64(config.PeriodMs))
		if err != nil {
			logger.Errorf("create NewQPSLimiter error, err: %s", err)
			return nil, err
		}
		l.limiter = limiter
		return l, nil
	} else if config.LimitStrategy == RateLimiterStrategy {
		limiter, err := NewRateLimiter(int64(config.MaxAllows), int64(config.PeriodMs), float64(config.MaxBurstRatio))
		if err != nil {
			logger.Errorf("create NewRateLimiter error, err: %s", err)
			return nil, err
		}
		l.limiter = limiter
		return l, nil
	}
	return nil, errors.New("Unknown LimitStrategy type:" + config.LimitStrategy)
}

// OverLimit check limit
func (engine *LimitEngine) OverLimit() bool {
	return !engine.limiter.TryAcquire()
}

