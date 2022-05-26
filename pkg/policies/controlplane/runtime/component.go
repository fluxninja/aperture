package runtime

// Component is the interface that all components must implement.
type Component interface {
	Execute(inPortReadings PortToValue, tickInfo TickInfo) (outPortReadings PortToValue, err error)
}
