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
)

// PulseGenerator is a component that accumulates sum of signal every tick.
type PulseGenerator struct {
	trueFor     time.Duration
	falseFor    time.Duration
	trueWindow  uint32
	falseWindow uint32
	windowCount uint32
	emitTrue    bool
}

// Name implements runtime.Component.
func (*PulseGenerator) Name() string { return "PulseGenerator" }

// Type implements runtime.Component.
func (*PulseGenerator) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// ShortDescription implements runtime.Component.
func (pg *PulseGenerator) ShortDescription() string {
	return fmt.Sprintf("true for: %d, false for: %d", pg.trueFor, pg.falseFor)
}

// IsActuator implements runtime.Component.
func (*PulseGenerator) IsActuator() bool { return false }

// NewPulseGeneratorAndOptions creates an pulse generator component and its fx options.
func NewPulseGeneratorAndOptions(generatorProto *policylangv1.PulseGenerator, _ runtime.ComponentID, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	evaluationPeriod := policyReadAPI.GetEvaluationInterval()
	trueFor := math.Ceil(float64(generatorProto.TrueFor.AsDuration()) / float64(evaluationPeriod))
	falseFor := math.Ceil(float64(generatorProto.FalseFor.AsDuration()) / float64(evaluationPeriod))

	pulseGenerator := &PulseGenerator{
		trueFor:     generatorProto.TrueFor.AsDuration(),
		falseFor:    generatorProto.FalseFor.AsDuration(),
		trueWindow:  uint32(trueFor),
		falseWindow: uint32(falseFor),
		windowCount: 0,
		emitTrue:    true,
	}
	return pulseGenerator, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (pg *PulseGenerator) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	outputValue := 1.0
	pg.windowCount++

	if pg.emitTrue {
		if pg.windowCount >= pg.trueWindow {
			pg.emitTrue = false
			pg.windowCount = 0
		}
	} else {
		if pg.windowCount >= pg.falseWindow {
			pg.emitTrue = true
			pg.windowCount = 0
		}
		outputValue = 0.0
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{runtime.NewReading(outputValue)},
	}, nil
}

// DynamicConfigUpdate is a no-op for PulseGenerator.
func (*PulseGenerator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
