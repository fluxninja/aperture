package runtime

import (
	"math"
)

var _ Reading = (*reading)(nil)

// Reading is the interface that reading implements which wraps a float64 value which can be valid or invalid.
type Reading interface {
	Value() float64
	Valid() bool
}

type reading struct {
	value float64
}

// NewReading creates a reading with the given float64 value.
func NewReading(value float64) Reading {
	return &reading{
		value: value,
	}
}

// NewBoolReading creates a reading with the given bool value.
func NewBoolReading(value bool) Reading {
	if value {
		return NewReading(1)
	}
	return NewReading(0)
}

// InvalidReading creates a reading with 'value' set to math.NaN(), invalid value.
func InvalidReading() Reading {
	return &reading{
		value: math.NaN(),
	}
}

// Value returns the value of the reading.
func (r *reading) Value() float64 {
	return r.value
}

// Valid checks whether the value of the reading is valid or invalid.
func (r *reading) Valid() bool {
	return !math.IsNaN(r.value)
}
