package deadlinemargin

import (
	"context"
	"time"
)

const (
	lookupMargin = 10 * time.Millisecond
)

// IsMarginExceeded returns true if the deadline will be passed in the next 10ms.
func IsMarginExceeded(ctx context.Context) bool {
	deadline, deadlineOK := ctx.Deadline()
	if deadlineOK {
		// check if deadline will be passed in the next 10ms
		deadline = deadline.Add(-lookupMargin)
		return time.Now().After(deadline)
	}
	return false
}
