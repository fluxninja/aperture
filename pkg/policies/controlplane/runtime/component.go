package runtime

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

// Component is the interface that all components must implement.
type Component interface {
	Execute(inPortReadings PortToValue, tickInfo TickInfo) (outPortReadings PortToValue, err error)
	DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller)
	// Generic name of the component, eg. "Gradient"
	Name() string
	Type() ComponentType
}

// NewDummyComponent creates a standalone component which won't accept or emit any signals
func NewDummyComponent(name string) Component { return dummyComponent{name: name} }

type dummyComponent struct{ name string }

func (c dummyComponent) Execute(inPortReadings PortToValue, tickInfo TickInfo) (outPortReadings PortToValue, err error) {
	return nil, nil
}
func (c dummyComponent) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (c dummyComponent) Name() string        { return c.name }
func (c dummyComponent) Type() ComponentType { return ComponentTypeStandAlone }
