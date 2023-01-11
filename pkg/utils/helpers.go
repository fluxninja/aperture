package utils

import (
	"regexp"
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

// SliceToMap converts a slice of string to a map[string]bool.
func SliceToMap(a []string) map[string]bool {
	m := make(map[string]bool)
	for _, n := range a {
		m[n] = true
	}
	return m
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
