// tristate is a helper package for tri-state boolean logic, which is used for
// logical combinator components.
package tristate

import "github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"

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
func (b Bool) ToReading() runtime.Reading {
	switch b {
	case True:
		return runtime.NewReading(1)
	case False:
		return runtime.NewReading(0)
	default:
		return runtime.InvalidReading()
	}
}

// FromReading interprets runtime.Reading as Bool, mapping 0 to False, any
// valid non-zero to True and Invalid to Unknown.
//
// (It's the same mapping as in ToReading, but allowing more truthy values).
func FromReading(reading runtime.Reading) Bool {
	if !reading.Valid() {
		return Unknown
	}
	if reading.Value() == 0 {
		return False
	}
	return True
}

// IsTrue returns true if Bool is True.
func (b Bool) IsTrue() bool {
	return b == True
}

// IsFalse returns true if Bool is False.
func (b Bool) IsFalse() bool {
	return b == False
}

// IsUnknown returns true if Bool is Unknown.
func (b Bool) IsUnknown() bool {
	return b == Unknown
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
