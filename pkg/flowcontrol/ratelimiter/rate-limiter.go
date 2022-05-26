package ratelimiter

// RateLimiter is a generic limiter interface.
type RateLimiter interface {
	Name() string
	Take(label string) (ok bool, remaining int, current int)
	TakeN(label string, count int) (ok bool, remaining int, current int)
	GetRateLimitCheck() RateLimitCheck
	Close() error
}

// RateLimitCheck is a generic limit checker interface.
type RateLimitCheck interface {
	CheckRateLimit(label string, count int) (ok bool, remaining int)
	SetRateLimit(limit int)
	GetRateLimit() int
}
