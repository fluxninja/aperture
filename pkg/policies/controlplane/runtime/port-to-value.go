package runtime

import "github.com/FluxNinja/aperture/pkg/policies/controlplane/reading"

// PortToValue is a map from port name to a slice of Readings.
type PortToValue map[string][]reading.Reading

// ReadSingleValuePort returns the reading of the first signal at port=portName. If no signal is found, reading.NewInvalid() is returned.
func (ptv PortToValue) ReadSingleValuePort(portName string) reading.Reading {
	retReading := reading.NewInvalid()
	signals, ok := ptv[portName]
	if ok {
		if len(signals) > 0 {
			retReading = signals[0]
		}
	}
	return retReading
}

// ReadRepeatedValuePort returns the reading of all the signals at port=portName. If no signal is found, []reading.Reading{} empty list is returned.
func (ptv PortToValue) ReadRepeatedValuePort(portName string) []reading.Reading {
	retReadings := []reading.Reading{}
	signals, ok := ptv[portName]
	if ok {
		retReadings = signals
	}
	return retReadings
}
