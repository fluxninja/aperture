package runtime

import (
	"fmt"
	"time"

	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ TickInfo = (*tickInfo)(nil)

// TickInfo is the interface that trackInfo implements which contains information about the current tick.
type TickInfo interface {
	Timestamp() time.Time
	Tick() int
	Interval() time.Duration
	Serialize() *policysyncv1.TickInfo
	String() string
}

type tickInfo struct {
	timestamp time.Time
	tick      int
	interval  time.Duration
}

// NewTickInfo returns a Tickinfo.
func NewTickInfo(timestamp time.Time, tick int, interval time.Duration) TickInfo {
	return &tickInfo{
		timestamp: timestamp,
		tick:      tick,
		interval:  interval,
	}
}

// Timestamp returns the timestamp of the tick.
func (tickInfo *tickInfo) Timestamp() time.Time {
	return tickInfo.timestamp
}

// Tick returns the tick of the tickInfo.
func (tickInfo *tickInfo) Tick() int {
	return tickInfo.tick
}

// Interval returns the interval of the tickInfo.
func (tickInfo *tickInfo) Interval() time.Duration {
	return tickInfo.interval
}

// Serialize returns the proto serialized version of the tickInfo.
func (tickInfo *tickInfo) Serialize() *policysyncv1.TickInfo {
	return &policysyncv1.TickInfo{
		Timestamp: timestamppb.New(tickInfo.timestamp),
		Tick:      (int64)(tickInfo.tick),
		Interval:  durationpb.New(tickInfo.interval),
	}
}

// String returns the string representation of the tickInfo.
func (tickInfo *tickInfo) String() string {
	return fmt.Sprintf("TickInfo: {timestamp: %s, tick: %d, interval: %s}",
		tickInfo.timestamp, tickInfo.tick, tickInfo.interval)
}
