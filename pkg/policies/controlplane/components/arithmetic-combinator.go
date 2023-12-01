package components

import (
	"errors"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// ArithmeticOperator is the type of arithmetic operation.
type ArithmeticOperator int8

//go:generate enumer -type=ArithmeticOperator -transform=lower -output=arithmetic-operator-string.go
const (
	UnknownArithmetic ArithmeticOperator = iota
	Add
	Sub
	Mul
	Div
	Xor
	LShift
	RShift
)

// ArithmeticCombinator takes lhs, rhs input signals and emits computed output via arithmetic operation.
type ArithmeticCombinator struct {
	// The arithmetic operation can be addition, subtraction, multiplication, division, XOR, left bit shift or right bit shift.
	operator ArithmeticOperator
}

// Name implements runtime.Component.
func (*ArithmeticCombinator) Name() string { return "ArithmeticCombinator" }

// Type implements runtime.Component.
func (*ArithmeticCombinator) Type() runtime.ComponentType {
	return runtime.ComponentTypeSignalProcessor
}

// ShortDescription implements runtime.Component.
func (arith *ArithmeticCombinator) ShortDescription() string { return arith.operator.String() }

// IsActuator implements runtime.Component.
func (*ArithmeticCombinator) IsActuator() bool { return false }

// Make sure ArithmeticCombinator complies with Component interface.
var _ runtime.Component = (*ArithmeticCombinator)(nil)

// NewArithmeticCombinatorAndOptions returns a new ArithmeticCombinator and its Fx options.
func NewArithmeticCombinatorAndOptions(arithmeticCombinatorProto *policylangv1.ArithmeticCombinator, _ runtime.ComponentID, _ iface.Policy) (runtime.Component, fx.Option, error) {
	operator, err := ArithmeticOperatorString(arithmeticCombinatorProto.Operator)
	if err != nil {
		return nil, fx.Options(), err
	}
	if operator == UnknownArithmetic {
		return nil, fx.Options(), errors.New("unknown arithmetic operator")
	}

	arith := ArithmeticCombinator{
		operator: operator,
	}
	return &arith, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (arith *ArithmeticCombinator) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	lhs := inPortReadings.ReadSingleReadingPort("lhs")
	rhs := inPortReadings.ReadSingleReadingPort("rhs")
	output := runtime.InvalidReading()
	var err error

	if lhs.Valid() && rhs.Valid() {
		lhsVal, rhsVal := lhs.Value(), rhs.Value()
		switch arith.operator {
		case Add:
			output = runtime.NewReading(lhsVal + rhsVal)
		case Sub:
			output = runtime.NewReading(lhsVal - rhsVal)
		case Mul:
			output = runtime.NewReading(lhsVal * rhsVal)
		case Div:
			if rhsVal == 0 {
				output = runtime.InvalidReading()
			} else {
				output = runtime.NewReading(lhsVal / rhsVal)
			}
		case Xor:
			output = runtime.NewReading(float64(int(lhsVal) ^ int(rhsVal)))
		case LShift:
			output = runtime.NewReading(float64(int(lhsVal) << int(rhsVal)))
		case RShift:
			output = runtime.NewReading(float64(int(lhsVal) >> int(rhsVal)))
		}
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, err
}

// DynamicConfigUpdate is a no-op for ArithmeticCombinator.
func (arith *ArithmeticCombinator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
