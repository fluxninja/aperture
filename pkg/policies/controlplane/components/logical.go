package components

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime/tristate"
)

// logicalCombinator is n-ary logical combinator used to implement And.
type logicalCombinator struct {
	op             func(tristate.Bool, tristate.Bool) tristate.Bool
	name           string
	neutralElement tristate.Bool
}

// Name implements runtime.Component.
func (c *logicalCombinator) Name() string { return c.name }

// Type implements runtime.Component.
func (c *logicalCombinator) Type() runtime.ComponentType {
	return runtime.ComponentTypeSignalProcessor
}

// ShortDescription implements runtime.Component.
func (c *logicalCombinator) ShortDescription() string { return "" }

// IsActuator implements runtime.Component.
func (*logicalCombinator) IsActuator() bool { return false }

// Execute implements runtime.Component.
func (c *logicalCombinator) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	inputs := inPortReadings.ReadRepeatedReadingPort("inputs")

	output := c.neutralElement
	for _, input := range inputs {
		output = c.op(output, tristate.FromReading(input))
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output.ToReading()},
	}, nil
}

// DynamicConfigUpdate is a no-op for logicalCombinator.
func (*logicalCombinator) DynamicConfigUpdate(notifiers.Event, config.Unmarshaller) {}

// NewAndAndOptions creates a new And Component.
func NewAndAndOptions(
	_ *policylangv1.And,
	_ runtime.ComponentID,
	_ iface.Policy,
) (runtime.Component, fx.Option, error) {
	return &logicalCombinator{
		neutralElement: tristate.True,
		op:             tristate.Bool.And,
		name:           "And",
	}, fx.Options(), nil
}

// NewOrAndOptions creates a new Or Component.
func NewOrAndOptions(
	_ *policylangv1.Or,
	_ runtime.ComponentID,
	_ iface.Policy,
) (runtime.Component, fx.Option, error) {
	return &logicalCombinator{
		neutralElement: tristate.False,
		op:             tristate.Bool.Or,
		name:           "Or",
	}, fx.Options(), nil
}

type inverter struct{}

// Name implements runtime.Component.
func (c *inverter) Name() string { return "Inverter" }

// Type implements runtime.Component.
func (c *inverter) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (c *inverter) ShortDescription() string { return "" }

// IsActuator implements runtime.Component.
func (*inverter) IsActuator() bool { return false }

// Execute implements runtime.Component.
func (c *inverter) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	input := inPortReadings.ReadSingleReadingPort("input")

	return runtime.PortToReading{
		"output": []runtime.Reading{
			tristate.FromReading(input).Not().ToReading(),
		},
	}, nil
}

// DynamicConfigUpdate is a no-op for not component.
func (*inverter) DynamicConfigUpdate(notifiers.Event, config.Unmarshaller) {}

// NewInverterAndOptions creates a new Inverter Component.
func NewInverterAndOptions(
	_ *policylangv1.Inverter,
	_ runtime.ComponentID,
	_ iface.Policy,
) (runtime.Component, fx.Option, error) {
	return &inverter{}, fx.Options(), nil
}
