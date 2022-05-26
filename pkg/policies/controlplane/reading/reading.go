package reading

import "math"

// Reading wraps a float Value which can be valid or invalid.
type Reading struct {
	Value float64
	Valid bool
}

// New creates a reading and checks if value is NaN.
func New(value float64) Reading {
	valid := !math.IsNaN(value)
	return Reading{Value: value, Valid: valid}
}

// NewInvalid crates a reading with 'valid' value set to false.
func NewInvalid() Reading {
	return Reading{Valid: false}
}
