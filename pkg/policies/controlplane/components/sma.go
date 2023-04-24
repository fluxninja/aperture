package components

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// SMA is a Simple Moving Average filter.
type SMA struct {
	policyReadAPI     iface.Policy
	sum               float64
	window            int
	buffer            []float64
	lastGoodOutput    runtime.Reading
	validDuringWarmup bool
	invalidCounter    int
}

// Make sure SMA complies with Component interface.
var _ runtime.Component = (*SMA)(nil)

// NewSMAAndOptions returns a new SMA filter and its Fx options.
func NewSMAAndOptions(smaProto *policylangv1.SMA, _ string, policyReadAPI iface.Policy) (*SMA, fx.Option, error) {
	params := smaProto.GetParameters()

	sma := &SMA{
		policyReadAPI:     policyReadAPI,
		window:            int(params.SmaWindow.AsDuration() / policyReadAPI.GetEvaluationInterval()),
		buffer:            make([]float64, 0),
		lastGoodOutput:    runtime.InvalidReading(),
		validDuringWarmup: params.ValidDuringWarmup,
		invalidCounter:    0,
	}

	return sma, fx.Options(), nil
}

// Name returns the name of the component.
func (sma *SMA) Name() string {
	return "SMA"
}

// Type returns the type of the component.
func (sma *SMA) Type() runtime.ComponentType {
	return runtime.ComponentTypeSignalProcessor
}

// ShortDescription returns a short description of the component.
func (sma *SMA) ShortDescription() string {
	return fmt.Sprintf("window: %d", sma.window)
}

// IsActuator returns whether this component is a actuator or not.
func (*SMA) IsActuator() bool { return false }

// Execute implements runtime.Component.Execute.
func (sma *SMA) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	input := inPortReadings.ReadSingleReadingPort("input")
	var output runtime.Reading

	if input.Valid() {
		sma.buffer = append(sma.buffer, input.Value())
		if len(sma.buffer) > sma.window {
			sma.sum -= sma.buffer[0]
			sma.buffer = sma.buffer[1:]
		}
		sma.sum += input.Value()
		if len(sma.buffer) == sma.window || sma.validDuringWarmup {
			avg := sma.computeAverage()
			output = runtime.NewReading(avg)
			sma.invalidCounter = 0
		} else {
			output = runtime.InvalidReading()
		}
	} else {
		sma.invalidCounter++
		if sma.invalidCounter >= sma.window {
			sma.buffer = []float64{}
			sma.sum = 0
			sma.invalidCounter = 0
			output = runtime.InvalidReading()
		} else {
			output = sma.lastGoodOutput
		}
	}

	// Set the last good output
	if output.Valid() {
		sma.lastGoodOutput = output
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for SMA.
func (sma *SMA) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}

func (sma *SMA) computeAverage() float64 {
	if len(sma.buffer) > 0 {
		return sma.sum / float64(len(sma.buffer))
	}
	return 0
}
