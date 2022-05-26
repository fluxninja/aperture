package runtime

import "time"

// TickInfo is a struct that contains information about the current tick.
type TickInfo struct {
	Timestamp     time.Time
	NextTimestamp time.Time
	Tick          int
	Interval      time.Duration
}
