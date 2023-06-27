package utils

import (
	"regexp"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

// SliceFind returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func SliceFind(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

// SliceContains tells whether a contains x.
func SliceContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// RemoveFromSlice removes x from a and returns the new slice.
func RemoveFromSlice(a []string, x string) []string {
	for i, n := range a {
		if x == n {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
}

// SliceToSet converts a slice of strings to a set.
func SliceToSet(a []string) Set[string] {
	m := make(map[string]struct{}, len(a))
	for _, n := range a {
		m[n] = struct{}{}
	}
	return Set[string](m)
}

// Set represents set of values.
type Set[T comparable] map[T]struct{}

// Contains returns whether the set contains given element.
func (s Set[T]) Contains(x T) bool {
	_, exists := s[x]
	return exists
}

// IsHTTPUrl returns true if the given string is an HTTP(S) URL.
func IsHTTPUrl(url string) bool {
	prefixHTTPRegex := "~^(?:f|ht)tps?://~i"
	matched, _ := regexp.MatchString(prefixHTTPRegex, url)
	return matched
}

// Mod is a normal modulo operation which does not produce negative values.
func Mod(a, b int) int {
	return (a%b + b) % b
}

// Shutdown invokes fx.Shutdowner and logs error if any.
func Shutdown(shutdowner fx.Shutdowner) {
	if err := shutdowner.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Unable to shutdown!")
	}
}
