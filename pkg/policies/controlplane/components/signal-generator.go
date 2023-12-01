package components

import (
	"fmt"
	"math"
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime/tristate"
)

// SignalGenerator generates a signal based on the steps specified.
type SignalGenerator struct {
	steps            []*policylangv1.SignalGenerator_Parameters_Step
	currentStep      int
	tickCount        int32
	atStart          bool
	atEnd            bool
	evaluationPeriod time.Duration
}

// Name implements runtime.Component.
func (*SignalGenerator) Name() string { return "SignalGenerator" }

// Type implements runtime.Component.
func (*SignalGenerator) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// ShortDescription implements runtime.Component.
func (sg *SignalGenerator) ShortDescription() string {
	return fmt.Sprintf("SignalGenerator with %d steps", len(sg.steps))
}

// IsActuator implements runtime.Component.
func (*SignalGenerator) IsActuator() bool { return false }

// NewSignalGeneratorAndOptions creates a signal generator component and its fx options.
func NewSignalGeneratorAndOptions(generatorProto *policylangv1.SignalGenerator, _ runtime.ComponentID, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	evaluationPeriod := policyReadAPI.GetEvaluationInterval()

	signalGenerator := &SignalGenerator{
		steps:            generatorProto.Parameters.Steps,
		currentStep:      0,
		tickCount:        0,
		atStart:          true,
		atEnd:            false,
		evaluationPeriod: evaluationPeriod,
	}
	return signalGenerator, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (sg *SignalGenerator) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	resetVal := inPortReadings.ReadSingleReadingPort("reset")
	forwardVal := inPortReadings.ReadSingleReadingPort("forward")
	backwardVal := inPortReadings.ReadSingleReadingPort("backward")

	direction := 0
	forward := tristate.FromReading(forwardVal).IsTrue()
	backward := tristate.FromReading(backwardVal).IsTrue()

	if (forward && backward) ||
		(!forward && !backward) {
		direction = 0
	} else if forward {
		direction = 1
	} else if backward {
		direction = -1
	}

	if tristate.FromReading(resetVal).IsTrue() {
		sg.currentStep = 0
		sg.tickCount = 0
		sg.atStart = true
		sg.atEnd = false
	} else if direction == 1 {
		if sg.atStart {
			sg.atStart = false
		}

		if !sg.atEnd {
			sg.tickCount++
			// Process tick updates currentStep and atEnd if necessary
			sg.processForwardTick()
		}
	} else if direction == -1 {
		if sg.atEnd {
			sg.atEnd = false
		}

		if !sg.atStart {
			sg.tickCount--
			// Process tick updates currentStep and atStart if necessary
			sg.processBackwardTick()
		}
	}

	currentSignal := runtime.ConstantSignalFromProto(sg.steps[sg.currentStep].TargetOutput)
	currentValue := currentSignal.Float()

	if isInterpolable(currentSignal) {
		if sg.currentStep > 0 {
			tickFraction := float64(sg.tickCount) / float64(sg.getStepDuration())
			previousSignal := runtime.ConstantSignalFromProto(sg.steps[sg.currentStep-1].TargetOutput)
			if isInterpolable(previousSignal) {
				previousValue := previousSignal.Float()
				// If the previous and current values are not equal, interpolate between them. Avoid interpolation if the values are equal to prevent floating point errors.
				if currentValue != previousValue {
					// Interpolate between previous and current value
					currentValue = previousValue + (currentValue-previousValue)*tickFraction
				}
			}
		}
	}

	return runtime.PortToReading{
		"output":   []runtime.Reading{runtime.NewReading(currentValue)},
		"at_start": []runtime.Reading{runtime.NewBoolReading(sg.atStart)},
		"at_end":   []runtime.Reading{runtime.NewBoolReading(sg.atEnd)},
	}, nil
}

func isInterpolable(constantSignal *runtime.ConstantSignal) bool {
	return !constantSignal.IsSpecial()
}

func (sg *SignalGenerator) processForwardTick() {
	if sg.tickCount == sg.getStepDuration() {
		sg.currentStep++
		if sg.currentStep >= len(sg.steps) {
			sg.currentStep = len(sg.steps) - 1
			sg.atEnd = true
			sg.tickCount = sg.getStepDuration() - 1
		} else {
			sg.tickCount = 0
		}
	}
}

func (sg *SignalGenerator) processBackwardTick() {
	if sg.tickCount == -1 {
		sg.currentStep--
		if sg.currentStep < 0 {
			sg.currentStep = 0
			sg.atStart = true
			sg.tickCount = 0
		} else {
			sg.tickCount = sg.getStepDuration() - 1
		}
	}
}

func (sg *SignalGenerator) getStepDuration() int32 {
	stepDuration := sg.steps[sg.currentStep].Duration.AsDuration()
	if stepDuration == 0 {
		return 1
	}
	return int32(math.Ceil(float64(stepDuration) / float64(sg.evaluationPeriod)))
}

// DynamicConfigUpdate is a no-op for SignalGenerator.
func (*SignalGenerator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
