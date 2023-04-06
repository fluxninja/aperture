package watchdogconfig

import "math"

// PolicyTempDisabled is a marker value for policies to signal that the policy
// is temporarily disabled. Use it when all hope is lost to turn around from
// significant memory pressure (such as when above an "extreme" watermark).
const PolicyTempDisabled uint64 = math.MaxUint64

// Policy is an interface for watchdog policies, eg. Watermarks or adaptive.
type Policy interface {
	// Returns next threshold or PolicyTempDisabled
	NextThreshold(total, used uint64) uint64
}
