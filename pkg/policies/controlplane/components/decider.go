package components

import (
	"fmt"
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

type deciderState int8

const (
	decidedFalse deciderState = iota
	pendingFalse
	decidedTrue
	pendingTrue
)

// ComparisonOperator is the type of comparison operator.
type ComparisonOperator int8

//go:generate enumer -type=ComparisonOperator -transform=lower -output=decider-comparison-operator-string.go
const (
	UnknownComparison ComparisonOperator = iota
	GT
	LT
	GTE
	LTE
	EQ
	NEQ
)

// Decider controller for testing.
type Decider struct {
	// Time at which state became true pending
	truePendingSince time.Time
	// Time at which state became false pending
	falsePendingSince time.Time
	// The duration of time the condition must be met before transitioning to 1.0 signal
	trueForDuration time.Duration
	// The duration of time the condition must be unmet before transitioning to 0.0 signal
	falseForDuration time.Duration
	// The current error correction state
	state deciderState
	// The comparison operator can be greater-than, less-than, greater-than-or-equal, less-than-or-equal, equal, or not-equal.
	operator ComparisonOperator
}

// Name implements runtime.Component.
func (*Decider) Name() string { return "Decider" }

// Type implements runtime.Component.
func (*Decider) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (dec *Decider) ShortDescription() string {
	return fmt.Sprintf("%s for %s", dec.operator.String(), dec.trueForDuration.String())
}

// IsActuator implements runtime.Component.
func (*Decider) IsActuator() bool { return false }

// Make sure Decider complies with Component interface.
var _ runtime.Component = (*Decider)(nil)

// NewDeciderAndOptions creates timed controller and its fx options.
func NewDeciderAndOptions(deciderProto *policylangv1.Decider, _ runtime.ComponentID, _ iface.Policy) (runtime.Component, fx.Option, error) {
	operator, err := ComparisonOperatorString(deciderProto.Operator)
	if err != nil {
		return nil, fx.Options(), err
	}
	if operator == UnknownComparison {
		return nil, fx.Options(), fmt.Errorf("unknown operator")
	}
	timed := &Decider{
		trueForDuration:   deciderProto.TrueFor.AsDuration(),
		falseForDuration:  deciderProto.FalseFor.AsDuration(),
		operator:          operator,
		state:             decidedFalse,
		truePendingSince:  time.Time{},
		falsePendingSince: time.Time{},
	}
	return timed, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (dec *Decider) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	onTrue := runtime.NewReading(1.0)
	onFalse := runtime.NewReading(0.0)
	lhs := inPortReadings.ReadSingleReadingPort("lhs")
	rhs := inPortReadings.ReadSingleReadingPort("rhs")

	// Default currentDecision to False
	currentDecision := false

	if lhs.Valid() && rhs.Valid() {
		lhsVal, rhsVal := lhs.Value(), rhs.Value()
		switch dec.operator {
		case GT:
			currentDecision = (lhsVal > rhsVal)
		case LT:
			currentDecision = (lhsVal < rhsVal)
		case GTE:
			currentDecision = (lhsVal >= rhsVal)
		case LTE:
			currentDecision = (lhsVal <= rhsVal)
		case EQ:
			currentDecision = (lhsVal == rhsVal)
		case NEQ:
			currentDecision = (lhsVal != rhsVal)
		}
	}

	decisionType := dec.computeDecisionType(currentDecision, tickInfo)

	var output runtime.Reading

	switch decisionType {
	case currentDecided:
		if currentDecision {
			output = onTrue
		} else {
			output = onFalse
		}
	case currentPending:
		if currentDecision {
			output = onFalse
		} else {
			output = onTrue
		}
	default:
		output = onFalse
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for Decider.
func (dec *Decider) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}

type decisionType int

const (
	currentPending decisionType = iota
	currentDecided
)

func (dec *Decider) computeDecisionType(currentDecision bool, tickInfo runtime.TickInfo) decisionType {
	// Reset pending for opposite of current decision
	dec.resetPendingSince(!currentDecision)

	if dec.isCurrentStateAchieved(currentDecision) {
		return currentDecided
	} else {
		pendingSince := dec.getPendingSince(currentDecision, tickInfo)
		// check how much time has elapsed since the pending state was set
		if tickInfo.Timestamp().Sub(pendingSince) < dec.trueForDuration {
			dec.setPending(currentDecision)
			return currentPending
		} else {
			dec.resetPendingSince(currentDecision)
			dec.setDecided(currentDecision)
			return currentDecided
		}
	}
}

func (dec *Decider) isCurrentStateAchieved(currentDecision bool) bool {
	if currentDecision {
		return dec.state == decidedTrue
	} else {
		return dec.state == decidedFalse
	}
}

func (dec *Decider) resetPendingSince(currentDecision bool) {
	if currentDecision {
		dec.truePendingSince = time.Time{}
	} else {
		dec.falsePendingSince = time.Time{}
	}
}

func (dec *Decider) setPending(currentDecision bool) {
	if currentDecision {
		dec.state = pendingTrue
	} else {
		dec.state = pendingFalse
	}
}

func (dec *Decider) setDecided(currentDecision bool) {
	if currentDecision {
		dec.state = decidedTrue
	} else {
		dec.state = decidedFalse
	}
}

func (dec *Decider) getPendingSince(currentDecision bool, tickInfo runtime.TickInfo) time.Time {
	if currentDecision {
		if dec.truePendingSince.IsZero() {
			dec.truePendingSince = tickInfo.Timestamp()
		}
		return dec.truePendingSince
	} else {
		if dec.falsePendingSince.IsZero() {
			dec.falsePendingSince = tickInfo.Timestamp()
		}
		return dec.falsePendingSince
	}
}
