package ratelimiter

// CheckRateLimit checks the limit.
func CheckRateLimit(count, limit float64) (bool, float64) {
	// limit < 0 means that there is no limit
	if limit < 0 {
		return true, -1
	}
	if count > limit {
		return false, 0
	}
	return true, limit - count
}
