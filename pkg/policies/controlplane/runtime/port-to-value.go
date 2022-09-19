package runtime

// PortToValue is a map from port name to a slice of Readings.
type PortToValue map[string][]Reading

// ReadSingleValuePort returns the reading of the first signal at port=portName. If no signal is found, InvalidReading() is returned.
func (ptv PortToValue) ReadSingleValuePort(portName string) Reading {
	retReading := InvalidReading()
	signals, ok := ptv[portName]
	if ok {
		if len(signals) > 0 {
			retReading = signals[0]
		}
	}
	return retReading
}

// ReadRepeatedValuePort returns the reading of all the signals at port=portName. If no signal is found, []Reading{} empty list is returned.
func (ptv PortToValue) ReadRepeatedValuePort(portName string) []Reading {
	retReadings := []Reading{}
	signals, ok := ptv[portName]
	if ok {
		retReadings = signals
	}
	return retReadings
}
