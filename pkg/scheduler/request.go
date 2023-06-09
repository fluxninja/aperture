package scheduler

import "math"

// Request is metadata for request in a flow that is to be allowed or dropped based on controlled delay and queue limits.
type Request struct {
	FairnessLabel  string  // for enforcing fairness
	Tokens         uint64  // tokens (e.g. expected latency or complexity) for this request
	WeightedTokens float64 // weighted tokens for this request
	Priority       uint8   // larger values represent higher priority
}

// NewRequest calculates the weighted tokens for a request based on its priority, and returns a new Request.
func NewRequest(fairnessLabel string, tokens uint64, priority uint8) *Request {
	invPriority := float64(1 / (uint64(math.MaxUint8-priority) + 1))
	weightedTokens := float64(tokens) * invPriority
	return &Request{
		FairnessLabel:  fairnessLabel,
		Tokens:         tokens,
		Priority:       priority,
		WeightedTokens: weightedTokens,
	}
}
