package iface

import "github.com/fluxninja/aperture/pkg/selectors"

// RateLimiter interface.
type RateLimiter interface {
	Limiter
	TakeN(labels selectors.Labels, count int) (label string, ok bool, remaining int, current int)
}
