package components

import (
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Extrapolator takes an input signal and emits an output signal.
type Extrapolator struct {
	// Maximum time interval for each extrapolation of signal is done; Reading becomes invalid after this interval.
	maxExtrapolationInterval time.Duration
	// The last output that was emitted as an output signal.
	lastOutput runtime.Reading
	// The last valid timestamp.
	lastValidTimestamp time.Time
}

// Make sure Extrapolator complies with Component interface.
var _ runtime.Component = (*Extrapolator)(nil)

// NewExtrapolatorAndOptions creates a new Extrapolator Component.
func NewExtrapolatorAndOptions(extrapolatorProto *policylangv1.Extrapolator, componentIndex int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	exp := Extrapolator{
		maxExtrapolationInterval: extrapolatorProto.MaxExtrapolationInterval.AsDuration(),
		lastOutput:               runtime.InvalidReading(),
		lastValidTimestamp:       time.Time{},
	}

	return &exp, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (exp *Extrapolator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	input := inPortReadings.ReadSingleValuePort("input")
	output := runtime.InvalidReading()

	if input.Valid() {
		output = input
		exp.lastOutput = output
		exp.lastValidTimestamp = tickInfo.Timestamp()
	} else {
		// Check if elapsed since lastValidTimestamp has reached the maximum Extrapolation interval.
		if tickInfo.Timestamp().Sub(exp.lastValidTimestamp) <= exp.maxExtrapolationInterval {
			// If the signal is invalid, it repeats the last value for up to maxExtrapolationInterval.
			output = exp.lastOutput
		} else {
			// When maxExtrapolationInterval is reached, reading becomes invalid.
			exp.lastOutput = output
		}
	}

	// If the signal returns, it resumes mirroring the input signal as output signal.
	return runtime.PortToValue{
		"output": []runtime.Reading{output},
	}, nil
}
