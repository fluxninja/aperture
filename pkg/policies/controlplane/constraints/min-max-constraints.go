package constraints

import (
	"errors"
	"math"
)

// MinMaxConstraints grouping constraints.
type MinMaxConstraints struct {
	Max float64
	Min float64
}

// ErrMinMaxConstraint custom error.
var ErrMinMaxConstraint = errors.New("min constraint cannot be greater than max constraint")

// NewMinMaxConstraints Default to largest float64 for max and smallest float64 for min.
func NewMinMaxConstraints() *MinMaxConstraints {
	return &MinMaxConstraints{Max: math.MaxFloat64, Min: -math.MaxFloat64}
}

// SetMax update max constraint if current max constraint is greater than provided max.
func (mmc *MinMaxConstraints) SetMax(max float64) error {
	if mmc.Max > max {
		// Make sure max is greater than min
		if mmc.Min > max {
			return ErrMinMaxConstraint
		}
		mmc.Max = max
	}
	return nil
}

// SetMin update min constraint if current min constraint is less than provided min.
func (mmc *MinMaxConstraints) SetMin(min float64) error {
	if mmc.Min < min {
		// Make sure min is less than max
		if mmc.Max < min {
			return ErrMinMaxConstraint
		}
		mmc.Min = min
	}
	return nil
}

// GetMax returns max.
func (mmc *MinMaxConstraints) GetMax() float64 {
	return mmc.Max
}

// GetMin returns min.
func (mmc *MinMaxConstraints) GetMin() float64 {
	return mmc.Min
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
	if value > mmc.Max {
		return mmc.Max, MaxConstraint
	}
	if value < mmc.Min {
		return mmc.Min, MinConstraint
	}
	return value, None
}
