package sim

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Input is a component that will emit signal from the given array – one
// Reading per tick. On subsequent ticks, it will emit Invalid readings.
type Input []runtime.Reading

// NewInput creates an input from an array of floats.
func NewInput(values []float64) *Input {
	input := Input(ToRtReadings(NewReadings(values)))
	return &input
}

// NewConstantInput creates a constant-value input.
func NewConstantInput(value float64) runtime.Component {
	return components.NewConstantSignal(value)
}

// Name implements runtime.Component.
func (i *Input) Name() string { return "TestInput" }

// Type implements runtime.Component.
func (i *Input) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// Execute implements runtime.Component.
func (i *Input) Execute(_ runtime.PortToReading, _ runtime.TickInfo) (runtime.PortToReading, error) {
	return runtime.PortToReading{
		"output": []runtime.Reading{i.execute()},
	}, nil
}

func (i *Input) execute() runtime.Reading {
	if len(*i) == 0 {
		return runtime.InvalidReading()
	}
	reading := (*i)[0]
	*i = (*i)[1:]
	return reading
}

// DynamicConfigUpdate implements runtime.Component.
func (i *Input) DynamicConfigUpdate(_ notifiers.Event, _ config.Unmarshaller) {}

// Output is a sink component to capture signal within a test.
type output struct {
	Readings []runtime.Reading
}

// Name implements runtime.Component.
func (o *output) Name() string { return "TestOutput" }

// Type implements runtime.Component.
func (o *output) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// TakeReadings returns the list of readings since previous TakeReadings() call
// (or since start).
func (o *output) TakeReadings() []runtime.Reading {
	readings := o.Readings
	o.Readings = nil
	return readings
}

// Execute implements runtime.Component.
func (o *output) Execute(ins runtime.PortToReading, _ runtime.TickInfo) (runtime.PortToReading, error) {
	o.Readings = append(o.Readings, ins.ReadSingleReadingPort("input"))
	return nil, nil
}

// DynamicConfigUpdate implements runtime.Component.
func (o *output) DynamicConfigUpdate(_ notifiers.Event, _ config.Unmarshaller) {}

/******** helpers ***********/

// ConfigureInputComponent takes an input component and creates a
// Component which outputs signal with given name on its "output"
// poruntime.
func ConfigureInputComponent(comp runtime.Component, signal runtime.SignalID) runtime.ConfiguredComponent {
	return runtime.ConfiguredComponent{
		Component: comp,
		PortMapping: runtime.PortMapping{
			Outs: map[string][]runtime.Signal{
				"output": {{
					SubCircuitID: signal.SubCircuitID,
					SignalName:   signal.SignalName,
				}},
			},
		},
	}
}

// ConfigureOutputComponent takes an output component and creates a
// Component which reads signal with a given name on its "input"
// poruntime.
func ConfigureOutputComponent(signal runtime.SignalID, comp runtime.Component) runtime.ConfiguredComponent {
	return runtime.ConfiguredComponent{
		Component: comp,
		PortMapping: runtime.PortMapping{
			Ins: map[string][]runtime.Signal{
				"input": {{
					SubCircuitID: signal.SubCircuitID,
					SignalName:   signal.SignalName,
				}},
			},
		},
	}
}
