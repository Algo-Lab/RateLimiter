package pkg

// Limiter limit
type Limiter interface {
	TryAcquire() bool
}

// limit strategy
const (
	QPSStrategy         = "QPS"
	RateLimiterStrategy = "RateLimiter"
)

