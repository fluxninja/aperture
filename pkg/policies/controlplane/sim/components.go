package sim

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	rt "github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Input is a component that will emit signal from the given array – one
// Reading per tick. On subsequent ticks, it will emit Invalid readings.
type Input []rt.Reading

// NewInput creates an input from an array of floats.
func NewInput(values []float64) *Input {
	input := Input(ToRtReadings(NewReadings(values)))
	return &input
}

// NewConstantInput creates a constant-value input.
func NewConstantInput(value float64) rt.Component {
	return components.NewConstant(value)
}

// Name implements runtime.Component.
func (i *Input) Name() string { return "TestInput" }

// Type implements runtime.Component.
func (i *Input) Type() rt.ComponentType { return rt.ComponentTypeSource }

// Execute implements runtime.Component.
func (i *Input) Execute(_ rt.PortToValue, _ rt.TickInfo) (rt.PortToValue, error) {
	return rt.PortToValue{
		"output": []rt.Reading{i.execute()},
	}, nil
}

func (i *Input) execute() rt.Reading {
	if len(*i) == 0 {
		return rt.InvalidReading()
	}
	reading := (*i)[0]
	*i = (*i)[1:]
	return reading
}

// DynamicConfigUpdate implements runtime.Component.
func (i *Input) DynamicConfigUpdate(_ notifiers.Event, _ config.Unmarshaller) {}

// Output is a sink component to capture signal within a test.
type output struct {
	Readings []rt.Reading
}

// Name implements runtime.Component.
func (o *output) Name() string { return "TestOutput" }

// Type implements runtime.Component.
func (o *output) Type() rt.ComponentType { return rt.ComponentTypeSink }

// TakeReadings returns the list of readings since previous TakeReadings() call
// (or since start).
func (o *output) TakeReadings() []rt.Reading {
	readings := o.Readings
	o.Readings = nil
	return readings
}

// Execute implements runtime.Component.
func (o *output) Execute(ins rt.PortToValue, _ rt.TickInfo) (rt.PortToValue, error) {
	o.Readings = append(o.Readings, ins.ReadSingleValuePort("input"))
	return nil, nil
}

// DynamicConfigUpdate implements runtime.Component.
func (o *output) DynamicConfigUpdate(_ notifiers.Event, _ config.Unmarshaller) {}

/******** helpers ***********/

// ConfigureInputComponent takes an input component and creates a
// ConfiguredComponent which outputs signal with given name on its "output"
// port.
func ConfigureInputComponent(comp rt.Component, signal string) rt.ConfiguredComponent {
	return rt.ConfiguredComponent{
		Component: comp,
		PortMapping: rt.PortMapping{
			Outs: map[string][]rt.Port{
				"output": {{SignalName: &signal}},
			},
		},
	}
}

// ConfigureOutputComponent takes an output component and creates a
// ConfiguredComponent which reads signal with a given name on its "input"
// port.
func ConfigureOutputComponent(signal string, comp rt.Component) rt.ConfiguredComponent {
	return rt.ConfiguredComponent{
		Component: comp,
		PortMapping: rt.PortMapping{
			Ins: map[string][]rt.Port{
				"input": {{SignalName: &signal}},
			},
		},
	}
}
