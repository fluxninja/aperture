package component

import (
	"errors"
	"math"

	"github.com/fluxninja/aperture/pkg/log"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/apis/policyapi"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/constraints"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/reading"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

type stage int

const (
	warmUpStage stage = iota
	emaStage
)

// EMA is an Exponential Moving Average filter.
type EMA struct {
	lastGoodOutput reading.Reading
	// The smoothing factor between 0-1. A higher value discounts older observations faster.
	alpha float64
	sum   float64
	count float64
	// The initial value of EMA is the average of the first warm_up_window number of observations.
	warmUpWindow uint32
	emaWindow    uint32
	warmupCount  uint32
	invalidCount uint32
	// The correction factor on the maximum relative to the signal
	correctionFactorOnMaxViolation float64
	// The correction factor on the minimum relative to the signal
	correctionFactorOnMinViolation float64
	currentStage                   stage
}

// Make sure EMA complies with Component interface.
var _ runtime.Component = (*EMA)(nil)

// NewEMAAndOptions returns a new EMA filter and its Fx options.
func NewEMAAndOptions(emaProto *policylangv1.EMA,
	_ int,
	policyReadAPI policyapi.PolicyReadAPI,
) (*EMA, fx.Option, error) {
	// period of tick
	evaluationPeriod := policyReadAPI.GetEvaluationInterval()
	// number of ticks in emaWindow
	emaWindow := math.Ceil(float64(emaProto.EmaWindow.AsDuration()) / float64(evaluationPeriod))

	alpha := 2.0 / (emaWindow + 1)
	warmUpWindow := uint32(math.Ceil(float64(emaProto.WarmUpWindow.AsDuration()) / float64(evaluationPeriod)))

	ema := &EMA{
		correctionFactorOnMinViolation: emaProto.CorrectionFactorOnMinEnvelopeViolation,
		correctionFactorOnMaxViolation: emaProto.CorrectionFactorOnMaxEnvelopeViolation,
		alpha:                          alpha,
		warmUpWindow:                   warmUpWindow,
		emaWindow:                      uint32(emaWindow),
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
	ema.lastGoodOutput = reading.NewInvalid()
}

// Execute implements runtime.Component.Execute.
func (ema *EMA) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	retErr := func(err error) (runtime.PortToValue, error) {
		return runtime.PortToValue{
			"output": []reading.Reading{reading.NewInvalid()},
		}, err
	}

	input := inPortReadings.ReadSingleValuePort("input")
	maxEnvelope := inPortReadings.ReadSingleValuePort("max_envelope")
	minEnvelope := inPortReadings.ReadSingleValuePort("min_envelope")
	output := reading.NewInvalid()

	switch ema.currentStage {
	case warmUpStage:
		ema.warmupCount++
		if input.Valid {
			ema.sum += input.Value
			ema.count++
			// Emit the avg for the valid values in the warmup window.
			avg, err := ema.computeAverage(minEnvelope, maxEnvelope)
			if err != nil {
				return retErr(err)
			}
			output = avg
			// Decide to switch to EMA stage
			if ema.warmupCount >= ema.warmUpWindow {
				ema.currentStage = emaStage
			}
		} else {
			// Emit the avg for the valid values in the warmup window.
			avg, err := ema.computeAverage(minEnvelope, maxEnvelope)
			if err != nil {
				return retErr(err)
			}
			output = avg
			// Immediately reset on any missing values during warmup.
			ema.resetStages()
		}
	case emaStage:
		if input.Valid {
			if !ema.lastGoodOutput.Valid {
				err := errors.New("ema: last good output is invalid")
				log.Error().Err(err).Msg("This is unexpected!")
				return retErr(err)
			}
			// Compute the new outputValue.
			outputValue := (ema.alpha * input.Value) + ((1 - ema.alpha) * ema.lastGoodOutput.Value)
			output = reading.New(outputValue)
		} else {
			ema.invalidCount++
			// emit last good EMA value
			output = ema.lastGoodOutput
			// If invalid count is greater than the ema window, reset the stages.
			if ema.invalidCount >= ema.emaWindow {
				ema.resetStages()
			}
		}
	default:
		log.Panic().Msg("unexpected ema stage")
	}

	// Set the last good output
	if output.Valid {
		ema.lastGoodOutput = output
	}
	// Returns Exponential Moving Average of a series of readings.
	return runtime.PortToValue{
		"output": []reading.Reading{output},
	}, nil
}

func (ema *EMA) computeAverage(minEnvelope, maxEnvelope reading.Reading) (reading.Reading, error) {
	if ema.count > 0 {
		avg := ema.sum / (ema.count)
		envelopedAvg, err := ema.applyEnvelope(avg, minEnvelope, maxEnvelope)
		if err != nil {
			return reading.NewInvalid(), err
		}
		return reading.New(envelopedAvg), nil
	} else {
		return reading.NewInvalid(), nil
	}
}

func (ema *EMA) applyEnvelope(input float64, minEnvelope, maxEnvelope reading.Reading) (float64, error) {
	minxMaxConstraints := constraints.NewMinMaxConstraints()
	if maxEnvelope.Valid {
		maxErr := minxMaxConstraints.SetMax(maxEnvelope.Value)
		if maxErr != nil {
			return 0, maxErr
		}
	}
	if minEnvelope.Valid {
		minErr := minxMaxConstraints.SetMin(minEnvelope.Value)
		if minErr != nil {
			return 0, minErr
		}
	}

	constrainedValue, constraintType := minxMaxConstraints.Constrain(input)
	var correctedConstrainedValue float64

	if constraintType == constraints.MinConstraint && ema.correctionFactorOnMinViolation != 1 {
		correctedConstrainedValue *= ema.correctionFactorOnMinViolation
	} else if constraintType == constraints.MaxConstraint && ema.correctionFactorOnMaxViolation != 1 {
		correctedConstrainedValue *= ema.correctionFactorOnMaxViolation
	} else {
		correctedConstrainedValue = constrainedValue
	}
	return correctedConstrainedValue, nil
}
