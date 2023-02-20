package constraints

import (
	"errors"
	"math"
)

// MinMaxConstraints grouping constraints.
type MinMaxConstraints struct {
	max float64
	min float64
}

// ErrMinMaxConstraint custom error.
var ErrMinMaxConstraint = errors.New("min constraint cannot be greater than max constraint and vice versa")

// NewMinMaxConstraints Default to largest float64 for max and smallest float64 for min.
func NewMinMaxConstraints() *MinMaxConstraints {
	constraints := &MinMaxConstraints{}
	constraints.Reset()
	return constraints
}

// Reset reset min and max to default.
func (mmc *MinMaxConstraints) Reset() {
	mmc.ResetMax()
	mmc.ResetMin()
}

// AdjustMax updates max constraint if current max constraint is greater than provided max.
func (mmc *MinMaxConstraints) AdjustMax(max float64) error {
	if mmc.max > max {
		// Make sure max is greater than min
		if mmc.min > max {
			return ErrMinMaxConstraint
		}
		mmc.max = max
	}
	return nil
}

// ResetMax reset max constraint to largest float64.
func (mmc *MinMaxConstraints) ResetMax() {
	mmc.max = math.MaxFloat64
}

// AdjustMin updates min constraint if current min constraint is less than provided min.
func (mmc *MinMaxConstraints) AdjustMin(min float64) error {
	if mmc.min < min {
		// Make sure min is less than max
		if mmc.max < min {
			return ErrMinMaxConstraint
		}
		mmc.min = min
	}
	return nil
}

// ResetMin reset min constraint to smallest float64.
func (mmc *MinMaxConstraints) ResetMin() {
	mmc.min = -math.MaxFloat64
}

// GetMax returns max.
func (mmc *MinMaxConstraints) GetMax() float64 {
	return mmc.max
}

// GetMin returns min.
func (mmc *MinMaxConstraints) GetMin() float64 {
	return mmc.min
}

// ConstraintType indicates the type of constraint.
type ConstraintType int

const (
	// None indicates no constraint.
	None ConstraintType = iota
	// MinConstraint indicates min constraint.
	MinConstraint
	// MaxConstraint indicates max constraint.
	MaxConstraint
)

// Constrain sets constraints.
func (mmc *MinMaxConstraints) Constrain(value float64) (float64, ConstraintType) {
	if value > mmc.max {
		return mmc.max, MaxConstraint
	}
	if value < mmc.min {
		return mmc.min, MinConstraint
	}
	return value, None
}
