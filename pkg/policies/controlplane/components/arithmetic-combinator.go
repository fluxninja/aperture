package components

import (
	"errors"
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

type arithmeticOperator int8

//go:generate enumer -type=arithmeticOperator -output=arithmetic-operator-string.go
const (
	unknownArithmetic arithmeticOperator = iota
	add
	sub
	mul
	div
	xor
	lshift
	rshift
)

// ArithmeticCombinator takes lhs, rhs input signals and emits computed output via arithmetic operation.
type ArithmeticCombinator struct {
	// The arithmetic operation can be addition, subtraction, multiplication, division, XOR, left bit shift or right bit shift.
	operator arithmeticOperator
}

// Name implements runtime.Component.
func (*ArithmeticCombinator) Name() string { return "ArithmeticCombinator" }

// Type implements runtime.Component.
func (*ArithmeticCombinator) Type() runtime.ComponentType {
	return runtime.ComponentTypeSignalProcessor
}

// Make sure ArithmeticCombinator complies with Component interface.
var _ runtime.Component = (*ArithmeticCombinator)(nil)

// NewArithmeticCombinatorAndOptions returns a new ArithmeticCombinator and its Fx options.
func NewArithmeticCombinatorAndOptions(arithmeticCombinatorProto *policylangv1.ArithmeticCombinator, _ int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	operator, err := arithmeticOperatorString(arithmeticCombinatorProto.Operator)
	if err != nil {
		return nil, fx.Options(), err
	}
	if operator == unknownArithmetic {
		return nil, fx.Options(), errors.New("unknown arithmetic operator")
	}

	arith := ArithmeticCombinator{
		operator: operator,
	}
	return &arith, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (arith *ArithmeticCombinator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	lhs := inPortReadings.ReadSingleValuePort("lhs")
	rhs := inPortReadings.ReadSingleValuePort("rhs")
	output := runtime.InvalidReading()
	var err error

	if lhs.Valid() && rhs.Valid() {
		lhsVal, rhsVal := lhs.Value(), rhs.Value()
		switch arith.operator {
		case add:
			output = runtime.NewReading(lhsVal + rhsVal)
		case sub:
			output = runtime.NewReading(lhsVal - rhsVal)
		case mul:
			output = runtime.NewReading(lhsVal * rhsVal)
		case div:
			if rhsVal == 0 {
				err = fmt.Errorf("divide by zero")
			} else {
				output = runtime.NewReading(lhsVal / rhsVal)
			}
		case xor:
			output = runtime.NewReading(float64(int(lhsVal) ^ int(rhsVal)))
		case lshift:
			output = runtime.NewReading(float64(int(lhsVal) << int(rhsVal)))
		case rshift:
			output = runtime.NewReading(float64(int(lhsVal) >> int(rhsVal)))
		}
	}

	return runtime.PortToValue{
		"output": []runtime.Reading{output},
	}, err
}

// DynamicConfigUpdate is a no-op for ArithmeticCombinator.
func (arith *ArithmeticCombinator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
