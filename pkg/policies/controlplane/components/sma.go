package components

import (
	"fmt"
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// SMA is a Simple Moving Average filter.
type SMA struct {
	lastGoodOutput    runtime.Reading
	policyReadAPI     iface.Policy
	sum               float64
	count             uint32
	runs              uint32
	invalidCount      uint32
	validDuringWarmup bool
}

// Make sure EMA complies with Component interface.
var _ runtime.Component = (*EMA)(nil)

// NewSMAAndOptions returns a new SMA filter and its Fx options.
func NewSMAAndOptions(smaProto *policylangv1.SMA, _ string, policyReadAPI iface.Policy) (*SMA, fx.Option, error) {
	evaluationInterval := policyReadAPI.GetEvaluationInterval()
	params := smaProto.GetParameters()

	runs := math.Ceil(float64(params.Window.AsDuration()) / float64(evaluationInterval))

	sma := &SMA{
		policyReadAPI:     policyReadAPI,
		runs:              uint32(runs),
		validDuringWarmup: params.ValidDuringWarmup,
	}
	sma.reset()
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
	return fmt.Sprintf("runs: %v", sma.runs)
}

// Execute implements runtime.Component.Execute.
func (sma *SMA) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	// logger := sma.policyReadAPI.GetStatusRegistry().GetLogger()
	retErr := func(err error) (runtime.PortToReading, error) {
		return runtime.PortToReading{
			"output": []runtime.Reading{runtime.InvalidReading()},
		}, err
	}

	input := inPortReadings.ReadSingleReadingPort("input")
	output := runtime.InvalidReading()

	if input.Valid() {
		sma.sum += input.Value()
		sma.count++
		if sma.count <= sma.runs {
			if sma.validDuringWarmup {
				avg, err := sma.computeAverage()
				if err != nil {
					return retErr(err)
				}
				output = avg
			} else {
				output = runtime.InvalidReading()
			}
		}
		if sma.count == sma.runs {
			sma.reset()
		}
	} else {
		sma.invalidCount++
		if sma.invalidCount >= sma.runs {
			sma.reset()
		}
		output = sma.lastGoodOutput
	}

	// Set the last good output
	if output.Valid() {
		sma.lastGoodOutput = output
	}
	reading := runtime.PortToReading{"output": []runtime.Reading{output}}
	return reading, nil
}

// DynamicConfigUpdate is a no-op for SMA.
func (sma *SMA) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}

func (sma *SMA) reset() {
	sma.invalidCount = 0
	sma.sum = 0
	sma.count = 0
	sma.lastGoodOutput = runtime.InvalidReading()
}

func (sma *SMA) computeAverage() (runtime.Reading, error) {
	if sma.count > 0 {
		avg := sma.sum / float64(sma.count)
		return runtime.NewReading(avg), nil
	}
	return runtime.InvalidReading(), nil
}
