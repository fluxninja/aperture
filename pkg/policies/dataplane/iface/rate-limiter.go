package iface

// RateLimiter interface.
type RateLimiter interface {
	Limiter
	TakeN(labels map[string]string, count int) (label string, ok bool, remaining int, current int)
}
