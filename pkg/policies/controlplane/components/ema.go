package components

import (
	"errors"
	"fmt"
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

type stage int

const (
	warmUpStage stage = iota
	emaStage
)

// EMA is an Exponential Moving Average filter.
type EMA struct {
	lastGoodOutput runtime.Reading
	policyReadAPI  iface.Policy
	// The smoothing factor between 0-1. A higher value discounts older observations faster.
	alpha float64
	sum   float64
	count float64
	// The correction factor on the maximum relative to the signal
	correctionFactorOnMaxViolation float64
	// The correction factor on the minimum relative to the signal
	correctionFactorOnMinViolation float64
	currentStage                   stage
	// The initial value of EMA is the average of the first warmup_window number of observations.
	warmupWindow      uint32
	emaWindow         uint32
	warmupCount       uint32
	invalidCount      uint32
	validDuringWarmup bool
}

// Name implements runtime.Component.
func (*EMA) Name() string { return "EMA" }

// Type implements runtime.Component.
func (*EMA) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (ema *EMA) ShortDescription() string { return fmt.Sprintf("win: %v", ema.emaWindow) }

// IsActuator implements runtime.Component.
func (*EMA) IsActuator() bool { return false }

// Make sure EMA complies with Component interface.
var _ runtime.Component = (*EMA)(nil)

// NewEMAAndOptions returns a new EMA filter and its Fx options.
func NewEMAAndOptions(emaProto *policylangv1.EMA, _ runtime.ComponentID, policyReadAPI iface.Policy) (*EMA, fx.Option, error) {
	// period of tick
	evaluationPeriod := policyReadAPI.GetEvaluationInterval()
	params := emaProto.GetParameters()

	// number of ticks in emaWindow
	emaWindow := math.Ceil(float64(params.EmaWindow.AsDuration()) / float64(evaluationPeriod))

	alpha := 2.0 / (emaWindow + 1)
	warmupWindow := uint32(math.Ceil(float64(params.WarmupWindow.AsDuration()) / float64(evaluationPeriod)))

	ema := &EMA{
		correctionFactorOnMinViolation: params.CorrectionFactorOnMinEnvelopeViolation,
		correctionFactorOnMaxViolation: params.CorrectionFactorOnMaxEnvelopeViolation,
		alpha:                          alpha,
		warmupWindow:                   warmupWindow,
		emaWindow:                      uint32(emaWindow),
		policyReadAPI:                  policyReadAPI,
		validDuringWarmup:              params.ValidDuringWarmup,
	}
	ema.resetStages()
	return ema, fx.Options(), nil
}

func (ema *EMA) resetStages() {
	ema.currentStage = warmUpStage
	ema.warmupCount = 0
	ema.invalidCount = 0
	ema.sum = 0
	ema.count = 0
	ema.lastGoodOutput = runtime.InvalidReading()
}

// Execute implements runtime.Component.Execute.
func (ema *EMA) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	logger := ema.policyReadAPI.GetStatusRegistry().GetLogger()
	retErr := func(err error) (runtime.PortToReading, error) {
		return runtime.PortToReading{
			"output": []runtime.Reading{runtime.InvalidReading()},
		}, err
	}

	input := inPortReadings.ReadSingleReadingPort("input")
	maxEnvelope := inPortReadings.ReadSingleReadingPort("max_envelope")
	minEnvelope := inPortReadings.ReadSingleReadingPort("min_envelope")
	output := runtime.InvalidReading()

	switch ema.currentStage {
	case warmUpStage:
		ema.warmupCount++
		if input.Valid() {
			ema.sum += input.Value()
			ema.count++
			// Decide to switch to EMA stage
			if ema.warmupCount >= ema.warmupWindow {
				ema.currentStage = emaStage
			}
		} else {
			// Immediately reset on any missing values during warm-up.
			ema.resetStages()
		}
		// Emit valid output during emaStage or during warm-up if configured to do so.
		if ema.currentStage == emaStage || ema.validDuringWarmup {
			// Emit the avg value of input signal during the warm-up window.
			avg, err := ema.computeAverage()
			if err != nil {
				return retErr(err)
			}
			output = avg
		} else {
			output = runtime.InvalidReading()
		}
	case emaStage:
		if input.Valid() {
			if !ema.lastGoodOutput.Valid() {
				err := errors.New("ema: last good output is invalid")
				logger.Error().Err(err).Msg("This is unexpected!")
				return retErr(err)
			}
			// Compute the new outputValue.
			outputValue := (ema.alpha * input.Value()) + ((1 - ema.alpha) * ema.lastGoodOutput.Value())
			output = runtime.NewReading(outputValue)
		} else {
			ema.invalidCount++
			// If invalid count is greater than the ema window, reset the stages.
			if ema.invalidCount >= ema.emaWindow {
				ema.resetStages()
			}
			// emit last good EMA value
			output = ema.lastGoodOutput
		}
	default:
		logger.Panic().Msg("unexpected ema stage")
	}

	// Set the last good output
	if output.Valid() {
		// apply correction
		value := output.Value()
		if maxEnvelope.Valid() && value > maxEnvelope.Value() {
			value *= ema.correctionFactorOnMaxViolation
		}
		if minEnvelope.Valid() && value < minEnvelope.Value() {
			value *= ema.correctionFactorOnMinViolation
		}
		output = runtime.NewReading(value)
		ema.lastGoodOutput = output
	}
	// Returns Exponential Moving Average of a series of readings.
	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

func (ema *EMA) computeAverage() (runtime.Reading, error) {
	if ema.count > 0 {
		avg := ema.sum / (ema.count)
		return runtime.NewReading(avg), nil
	} else {
		return runtime.InvalidReading(), nil
	}
}

// DynamicConfigUpdate is a no-op for EMA.
func (ema *EMA) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}
