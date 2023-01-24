package components

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Holder is a component that holds the last valid signal value
// for the specified duration then waits for next valid value to hold.
type Holder struct {
	holdWindow     uint32
	windowCount    uint32
	holdPhase      bool
	currentReading runtime.Reading
}

// Name implements runtime.Component.
func (*Holder) Name() string { return "Holder" }

// Type implements runtime.Component.
func (*Holder) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// NewHolderAndOptions creates a holder component and its fx options.
func NewHolderAndOptions(holderProto *policylangv1.Holder, _ int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	evaluationPeriod := policyReadAPI.GetEvaluationInterval()
	holdFor := math.Ceil(float64(holderProto.HoldFor.AsDuration()) / float64(evaluationPeriod))

	holder := &Holder{
		holdWindow:  uint32(holdFor),
		windowCount: 0,
		holdPhase:   false,
	}
	return holder, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (h *Holder) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	input := inPortReadings.ReadSingleValuePort("input")
	output := runtime.InvalidReading()

	// starting with holdPhase = false
	if !h.holdPhase {
		if input.Valid() {
			h.holdNewValue(input)
			output = input
		}
	} else {
		h.windowCount++
		// hold_for is finished
		if h.windowCount >= h.holdWindow {
			// immediately hold new value
			if input.Valid() {
				h.holdNewValue(input)
				output = input
			} else {
				h.holdPhase = false
				h.windowCount = 0
			}
		} else {
			output = h.currentReading
		}
	}

	return runtime.PortToValue{
		"output": []runtime.Reading{output},
	}, nil
}

func (h *Holder) holdNewValue(newValue runtime.Reading) {
	h.currentReading = newValue
	h.holdPhase = true
	h.windowCount = 0
}

// DynamicConfigUpdate is a no-op for Holder.
func (*Holder) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
