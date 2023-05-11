package sim

import (
	"math"

	rt "github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// Reading is implementation of rt.Reading designed for tests.
//
// It avoids the problem of invalid readings being non-equal due to being
// represented as NaNs, so that comparing invalid readings works as expected.
type Reading struct {
	value   float64 // never NaN
	invalid bool    // value=0 if invalid=true
}

// Valid implements runtime.Reading.
func (r Reading) Valid() bool { return !r.invalid }

// Value implements runtime.Reading.
func (r Reading) Value() float64 { return r.value }

// NewReading creates a Reading from a float. NaN is treated as invalid reading.
func NewReading(value float64) Reading {
	if math.IsNaN(value) {
		return Reading{invalid: true}
	} else {
		return Reading{value: value}
	}
}

// InvalidReading creates a new invalid Reading.
func InvalidReading() Reading {
	return Reading{invalid: true}
}

// NewReadings creates a slice of readings from a slice of floats.
func NewReadings(values []float64) []Reading {
	readings := make([]Reading, 0, len(values))
	for _, value := range values {
		if math.IsNaN(value) {
			readings = append(readings, Reading{invalid: true})
		} else {
			readings = append(readings, Reading{value: value})
		}
	}
	return readings
}

// ReadingFromRt converts runtime.Reading to Reading.
func ReadingFromRt(rtReading rt.Reading) Reading {
	if rtReading.Valid() {
		return Reading{value: rtReading.Value()}
	} else {
		return Reading{invalid: true}
	}
}

// ReadingsFromRt converts runtime.Readings to Readings.
func ReadingsFromRt(rtReadings []rt.Reading) []Reading {
	readings := make([]Reading, 0, len(rtReadings))
	for _, rtReading := range rtReadings {
		readings = append(readings, ReadingFromRt(rtReading))
	}
	return readings
}

// ToRtReadings converts Readings to rt.Readings.
func ToRtReadings(readings []Reading) []rt.Reading {
	rtReadings := make([]rt.Reading, 0, len(readings))
	for _, reading := range readings {
		rtReadings = append(rtReadings, reading)
	}
	return rtReadings
}
