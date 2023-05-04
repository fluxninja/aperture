package components

import (
	"errors"
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// UnaryOp is the type of unary operation.
type UnaryOp int8

//go:generate enumer -type=UnaryOp -output=unary-op-string.go
const (
	UnknownUnary UnaryOp = iota
	Abs
	Acos
	Acosh
	Asin
	Asinh
	Atan
	Atanh
	Cbrt
	Ceil
	Cos
	Cosh
	Erf
	Erfc
	Erfcinv
	Erfinv
	Exp
	Exp2
	Expm1
	Floor
	Gamma
	J0
	J1
	Lgamma
	Log
	Log10
	Log1p
	Log2
	Round
	RoundToEven
	Sin
	Sinh
	Sqrt
	Tan
	Tanh
	Trunc
	Y0
	Y1
)

// UnaryOperator takes an input signal and emits Square Root of it multiplied by scale as output.
type UnaryOperator struct {
	operator UnaryOp
}

// Name implements runtime.Component.
func (*UnaryOperator) Name() string { return "UnaryOperator" }

// Type implements runtime.Component.
func (*UnaryOperator) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (unaryOperator *UnaryOperator) ShortDescription() string {
	return unaryOperator.operator.String()
}

// IsActuator implements runtime.Component.
func (*UnaryOperator) IsActuator() bool { return false }

// Make sure UnaryOperator complies with Component interface.
var _ runtime.Component = (*UnaryOperator)(nil)

// NewUnaryOperatorAndOptions creates a new UnaryOperator Component.
func NewUnaryOperatorAndOptions(unaryOperatorProto *policylangv1.UnaryOperator, _ runtime.ComponentID, _ iface.Policy) (runtime.Component, fx.Option, error) {
	operator, err := UnaryOpString(unaryOperatorProto.Operator)
	if err != nil {
		return nil, fx.Options(), err
	}
	if operator == UnknownUnary {
		return nil, fx.Options(), errors.New("unknown unary operator")
	}

	unary := &UnaryOperator{
		operator: operator,
	}
	return unary, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (unaryOperator *UnaryOperator) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	input := inPortReadings.ReadSingleReadingPort("input")
	output := runtime.InvalidReading()

	switch unaryOperator.operator {
	case Abs:
		output = runtime.NewReading(math.Abs(input.Value()))
	case Acos:
		output = runtime.NewReading(math.Acos(input.Value()))
	case Acosh:
		output = runtime.NewReading(math.Acosh(input.Value()))
	case Asin:
		output = runtime.NewReading(math.Asin(input.Value()))
	case Asinh:
		output = runtime.NewReading(math.Asinh(input.Value()))
	case Atan:
		output = runtime.NewReading(math.Atan(input.Value()))
	case Atanh:
		output = runtime.NewReading(math.Atanh(input.Value()))
	case Cbrt:
		output = runtime.NewReading(math.Cbrt(input.Value()))
	case Ceil:
		output = runtime.NewReading(math.Ceil(input.Value()))
	case Cos:
		output = runtime.NewReading(math.Cos(input.Value()))
	case Cosh:
		output = runtime.NewReading(math.Cosh(input.Value()))
	case Erf:
		output = runtime.NewReading(math.Erf(input.Value()))
	case Erfc:
		output = runtime.NewReading(math.Erfc(input.Value()))
	case Erfcinv:
		output = runtime.NewReading(math.Erfcinv(input.Value()))
	case Erfinv:
		output = runtime.NewReading(math.Erfinv(input.Value()))
	case Exp:
		output = runtime.NewReading(math.Exp(input.Value()))
	case Exp2:
		output = runtime.NewReading(math.Exp2(input.Value()))
	case Expm1:
		output = runtime.NewReading(math.Expm1(input.Value()))
	case Floor:
		output = runtime.NewReading(math.Floor(input.Value()))
	case Gamma:
		output = runtime.NewReading(math.Gamma(input.Value()))
	case J0:
		output = runtime.NewReading(math.J0(input.Value()))
	case J1:
		output = runtime.NewReading(math.J1(input.Value()))
	case Lgamma:
		out, _ := math.Lgamma(input.Value())
		output = runtime.NewReading(out)
	case Log:
		output = runtime.NewReading(math.Log(input.Value()))
	case Log10:
		output = runtime.NewReading(math.Log10(input.Value()))
	case Log1p:
		output = runtime.NewReading(math.Log1p(input.Value()))
	case Log2:
		output = runtime.NewReading(math.Log2(input.Value()))
	case Round:
		output = runtime.NewReading(math.Round(input.Value()))
	case RoundToEven:
		output = runtime.NewReading(math.RoundToEven(input.Value()))
	case Sin:
		output = runtime.NewReading(math.Sin(input.Value()))
	case Sinh:
		output = runtime.NewReading(math.Sinh(input.Value()))
	case Sqrt:
		output = runtime.NewReading(math.Sqrt(input.Value()))
	case Tan:
		output = runtime.NewReading(math.Tan(input.Value()))
	case Tanh:
		output = runtime.NewReading(math.Tanh(input.Value()))
	case Trunc:
		output = runtime.NewReading(math.Trunc(input.Value()))
	case Y0:
		output = runtime.NewReading(math.Y0(input.Value()))
	case Y1:
		output = runtime.NewReading(math.Y1(input.Value()))
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for UnaryOperator.
func (unaryOperator *UnaryOperator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
