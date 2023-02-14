package runtime

import (
	"math"
)

var _ Reading = (*reading)(nil)

// Reading is the interface that reading implements which wraps a float Value which can be valid or invalid.
type Reading interface {
	Value() float64
	Valid() bool
}

type reading struct {
	value float64
}

// NewReading creates a reading with the given value.
func NewReading(value float64) Reading {
	return &reading{
		value: value,
	}
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

// Bool is a tri-state boolean: False, True or Unknown.
//
// Operations use truth tables as in
// https://en.wikipedia.org/wiki/Three-valued_logic#Kleene_and_Priest_logics.
type Bool int

const (
	// False is definitely false.
	False Bool = 0

	// Unknown is "maybe false, maybe true".
	Unknown Bool = 1

	// True is definitely true.
	True Bool = 2
)

// Not returns result of logical negation.
func (b Bool) Not() Bool {
	return Bool(2 - int(b))
}

// And implements "logical and with indeterminacy".
func (b Bool) And(rhs Bool) Bool {
	return Bool(min(int(b), int(rhs)))
}

// Or implements "logical or with indeterminacy".
func (b Bool) Or(rhs Bool) Bool {
	return Bool(max(int(b), int(rhs)))
}

// ToReading converts tri-state Bool to runtime.Reading, mapping False to 0,
// True to 1 and Unknown to invalid reading.
func (b Bool) ToReading() Reading {
	switch b {
	case True:
		return NewReading(1)
	case False:
		return NewReading(0)
	default:
		return InvalidReading()
	}
}

// FromReading interprets runtime.Reading as Bool, mapping 0 to False, any
// valid non-zero to True and Invalid to Unknown.
//
// (It's the same mapping as in ToReading, but allowing more truthy values).
func FromReading(reading Reading) Bool {
	if !reading.Valid() {
		return Unknown
	}
	if reading.Value() == 0 {
		return False
	}
	return True
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
