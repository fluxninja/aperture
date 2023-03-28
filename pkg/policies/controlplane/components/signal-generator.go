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
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime/tristate"
)

// SignalGenerator generates a signal based on the steps specified.
type SignalGenerator struct {
	endBehavior  string
	steps        []*policylangv1.SignalGenerator_Parameters_Step
	currentStep  int
	tickCount    uint32
	stepDuration uint32
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
func NewSignalGeneratorAndOptions(generatorProto *policylangv1.SignalGenerator, _ string, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	evaluationPeriod := policyReadAPI.GetEvaluationInterval()

	signalGenerator := &SignalGenerator{
		steps:        generatorProto.Parameters.Steps,
		endBehavior:  generatorProto.Parameters.EndBehavior,
		currentStep:  0,
		tickCount:    0,
		stepDuration: uint32(math.Ceil(float64(generatorProto.Parameters.Steps[0].Duration.AsDuration()) / float64(evaluationPeriod))),
	}
	return signalGenerator, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (sg *SignalGenerator) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	resetVal := inPortReadings.ReadSingleReadingPort("reset")
	if tristate.FromReading(resetVal).IsTrue() {
		sg.currentStep = 0
		sg.tickCount = 0
		sg.stepDuration = uint32(math.Ceil(float64(sg.steps[sg.currentStep].Duration.AsDuration()) / float64(tickInfo.Interval())))
	} else {
		sg.tickCount++

		if sg.tickCount >= sg.stepDuration {
			sg.tickCount = 0

			sg.currentStep++

			if sg.currentStep >= len(sg.steps) {
				switch sg.endBehavior {
				case "loop":
					sg.currentStep = 0
				case "laststep":
					sg.currentStep = len(sg.steps) - 1
				default:
					return nil, fmt.Errorf("invalid end_behavior: %s", sg.endBehavior)
				}
			}
			sg.stepDuration = uint32(math.Ceil(float64(sg.steps[sg.currentStep].Duration.AsDuration()) / float64(tickInfo.Interval())))
		}
	}

	currentSignal := runtime.ConstantSignalFromProto(sg.steps[sg.currentStep].ConstantSignal)
	currentValue := currentSignal.Float()
	return runtime.PortToReading{
		"output": []runtime.Reading{runtime.NewReading(currentValue)},
	}, nil
}

// DynamicConfigUpdate is a no-op for SignalGenerator.
func (*SignalGenerator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
