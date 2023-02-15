package runtime

// PortToReading is a map from port name to a slice of Readings.
type PortToReading map[string][]Reading

// ReadSingleReadingPort returns the reading of the first signal at port=portName. If no signal is found, InvalidReading() is returned.
func (ptr PortToReading) ReadSingleReadingPort(portName string) Reading {
	retReading := InvalidReading()
	signals, ok := ptr[portName]
	if ok {
		if len(signals) > 0 {
			retReading = signals[0]
		}
	}
	return retReading
}

// ReadRepeatedReadingPort returns the reading of all the signals at port=portName. If no signal is found, []Reading{} empty list is returned.
func (ptr PortToReading) ReadRepeatedReadingPort(portName string) []Reading {
	retReadings := []Reading{}
	signals, ok := ptr[portName]
	if ok {
		retReadings = signals
	}
	return retReadings
}
