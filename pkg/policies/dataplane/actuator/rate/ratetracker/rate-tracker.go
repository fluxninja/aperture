package ratetracker

// RateTracker is a generic limiter interface.
type RateTracker interface {
	Name() string
	Take(label string) (ok bool, remaining int, current int)
	TakeN(label string, count int) (ok bool, remaining int, current int)
	GetRateLimitChecker() RateLimitChecker
	Close() error
}

// RateLimitChecker is a generic limit checker interface.
type RateLimitChecker interface {
	CheckRateLimit(label string, count int) (ok bool, remaining int)
	SetRateLimit(limit int)
	GetRateLimit() int
}
