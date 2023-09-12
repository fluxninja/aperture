package aperture

// FlowStatus represents status of flow execution.
type FlowStatus uint8

// User passes a code to indicate status of flow execution.
//
//go:generate enumer -type=FlowStatus -output=flow-status-string.go
const (
	// OK indicates successful flow execution.
	OK FlowStatus = iota
	// Error indicate error on flow execution.
	Error
)
