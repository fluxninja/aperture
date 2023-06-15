package aperture

// FlowDecision represents the decision on the flow execution made by Agent.
type FlowDecision uint8

//go:generate enumer -type=FlowDecision -output=flow-decision-string.go
const (
	// Accepted indicates flow should be allowed.
	Accepted FlowDecision = iota
	// Rejected indicates flow should be rejected.
	Rejected
	// Unreachable indicates there was no response from the Agent.
	Unreachable
)
