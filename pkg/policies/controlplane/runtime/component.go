package runtime

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

// Component is the interface that all components must implement.
type Component interface {
	Execute(inPortReadings PortToValue, tickInfo TickInfo) (outPortReadings PortToValue, err error)
	DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller)
}
