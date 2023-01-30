package runtime

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

// Component is the interface that all components must implement.
type Component interface {
	Execute(inPortReadings PortToReading, tickInfo TickInfo) (outPortReadings PortToReading, err error)
	DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller)
	// Generic name of the component, eg. "Gradient"
	Name() string
	Type() ComponentType
}

// ComponentType describes the type of a component based on its connectivity in the circuit.
type ComponentType string

const (
	// ComponentTypeStandAlone is a component that does not accept or emit any
	// signals.
	ComponentTypeStandAlone ComponentType = "StandAlone"
	// ComponentTypeSource is a component that emits output signal(s) but does
	// not accept an input signal.
	ComponentTypeSource ComponentType = "Source"
	// ComponentTypeSink is a component that accepts input signal(s) but does
	// not emit an output signal.
	ComponentTypeSink ComponentType = "Sink"
	// ComponentTypeSignalProcessor is a component that accepts input signal(s)
	// and emits output signal(s).
	ComponentTypeSignalProcessor ComponentType = "SignalProcessor"
)

// NewDummyComponent creates a component with provided name and type.
func NewDummyComponent(name string, componentType ComponentType) Component {
	return dummyComponent{
		name:          name,
		componentType: componentType,
	}
}

type dummyComponent struct {
	name          string
	componentType ComponentType
}

// Execute implements runtime.Component.
func (c dummyComponent) Execute(inPortReadings PortToReading, tickInfo TickInfo) (outPortReadings PortToReading, err error) {
	return nil, nil
}

// DynamicConfigUpdate implements runtime.Component.
func (c dummyComponent) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

// Name implements runtime.Component.
func (c dummyComponent) Name() string { return c.name }

// Type implements runtime.Component.
func (c dummyComponent) Type() ComponentType { return c.componentType }
