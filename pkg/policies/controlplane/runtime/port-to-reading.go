package runtime

import (
	"fmt"
	"log"
	"reflect"
)

// PortToReading is a map from port name to a slice of Readings.
type PortToReading map[string][]Reading

func validatePortName(portName string, protoStruct interface{}) {
	matched := false
	v := reflect.ValueOf(protoStruct)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == portName {
			matched = true
			break
		}
	}

	if !matched {
		err := fmt.Errorf("mismatched portName: %s", portName)
		log.Fatal(err)
	}
}

// ReadSingleReadingPort returns the reading of the first signal at port=portName. If no signal is found, InvalidReading() is returned.
func (ptr PortToReading) ReadSingleReadingPort(portName string, protoStruct interface{}) Reading {
	validatePortName(portName, protoStruct)

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
func (ptr PortToReading) ReadRepeatedReadingPort(portName string, protoStruct interface{}) []Reading {
	validatePortName(portName, protoStruct)

	retReadings := []Reading{}
	signals, ok := ptr[portName]
	if ok {
		retReadings = signals
	}
	return retReadings
}
