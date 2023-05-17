package ratetracker

// RateTracker is a generic limiter interface.
type RateTracker interface {
	Name() string
	Take(label string) (ok bool, remaining float64, current float64)
	TakeN(label string, count float64) (ok bool, remaining float64, current float64)
	GetRateLimitChecker() RateLimitChecker
	Close() error
}

// RateLimitChecker is a generic limit checker interface.
type RateLimitChecker interface {
	CheckRateLimit(label string, count float64) (ok bool, remaining float64)
	SetRateLimit(limit float64)
	GetRateLimit() float64
}
