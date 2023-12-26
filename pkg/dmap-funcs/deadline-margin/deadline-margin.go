package deadlinemargin

import (
	"context"
	"time"
)

const (
	lookupMargin = 10 * time.Millisecond
)

// IsMarginExceeded returns true if the deadline will be passed within the lookup margin.
func IsMarginExceeded(ctx context.Context) bool {
	return IsWaitMarginExceeded(ctx, 0)
}

// IsWaitMarginExceeded returns true if the deadline will be exceeded in the given wait time and lookup margin.
func IsWaitMarginExceeded(ctx context.Context, wait time.Duration) bool {
	deadline, deadlineOK := ctx.Deadline()
	if deadlineOK {
		// check if deadline will be passed in the next 10ms
		deadline = deadline.Add(-lookupMargin)
		return time.Now().Add(wait).After(deadline)
	}
	return false
}
