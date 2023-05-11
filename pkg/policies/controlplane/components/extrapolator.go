package components

import (
	"fmt"
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// Extrapolator takes an input signal and emits an output signal.
type Extrapolator struct {
	// The last valid timestamp.
	lastValidTimestamp time.Time
	// The last output that was emitted as an output signal.
	lastOutput runtime.Reading
	// Maximum time interval for each extrapolation of signal is done; Reading becomes invalid after this interval.
	maxExtrapolationInterval time.Duration
}

// Name implements runtime.Component.
func (*Extrapolator) Name() string { return "Extrapolator" }

// Type implements runtime.Component.
func (*Extrapolator) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (exp *Extrapolator) ShortDescription() string {
	return fmt.Sprintf("for: %s", exp.maxExtrapolationInterval)
}

// IsActuator implements runtime.Component.
func (*Extrapolator) IsActuator() bool { return false }

// Make sure Extrapolator complies with Component interface.
var _ runtime.Component = (*Extrapolator)(nil)

// NewExtrapolatorAndOptions creates a new Extrapolator Component.
func NewExtrapolatorAndOptions(extrapolatorProto *policylangv1.Extrapolator, _ runtime.ComponentID, _ iface.Policy) (runtime.Component, fx.Option, error) {
	exp := Extrapolator{
		maxExtrapolationInterval: extrapolatorProto.Parameters.MaxExtrapolationInterval.AsDuration(),
		lastOutput:               runtime.InvalidReading(),
		lastValidTimestamp:       time.Time{},
	}

	return &exp, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (exp *Extrapolator) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	input := inPortReadings.ReadSingleReadingPort("input")
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
	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for Extrapolator.
func (exp *Extrapolator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
