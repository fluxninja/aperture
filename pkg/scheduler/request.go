package scheduler

// Request is metadata for request in a flow that is to be allowed or dropped based on controlled delay and queue limits.
type Request struct {
	FairnessLabel string // for enforcing fairness
	Tokens        uint64 // tokens (e.g. expected latency or complexity) for this request
	Priority      uint8  // larger values represent higher priority
}

// NewRequest calculates the inverse priority and returns a new Request.
func NewRequest(fairnessLabel string, tokens uint64, priority uint8) *Request {
	return &Request{
		FairnessLabel: fairnessLabel,
		Tokens:        tokens,
		Priority:      priority,
	}
}
