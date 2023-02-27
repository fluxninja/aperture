package components

import (
	"fmt"
	"math"
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime/tristate"
)

// Holder is a component that holds the last valid signal value
// for the specified duration then waits for next valid value to hold.
type Holder struct {
	currentReading runtime.Reading
	holdFor        time.Duration
	holdWindow     uint32
	windowCount    uint32
	holdPhase      bool
}

// Name implements runtime.Component.
func (*Holder) Name() string { return "Holder" }

// Type implements runtime.Component.
func (*Holder) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (h *Holder) ShortDescription() string { return fmt.Sprintf("for: %s", h.holdFor) }

// NewHolderAndOptions creates a holder component and its fx options.
func NewHolderAndOptions(holderProto *policylangv1.Holder, _ string, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	evaluationPeriod := policyReadAPI.GetEvaluationInterval()
	holdWindow := uint32(math.Ceil(float64(holderProto.HoldFor.AsDuration()) / float64(evaluationPeriod)))

	holder := &Holder{
		holdFor:     holderProto.HoldFor.AsDuration(),
		holdWindow:  holdWindow,
		windowCount: 0,
		holdPhase:   false,
	}
	return holder, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (h *Holder) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	input := inPortReadings.ReadSingleReadingPort("input")
	reset := inPortReadings.ReadSingleReadingPort("reset")
	output := runtime.InvalidReading()

	if tristate.FromReading(reset).IsTrue() {
		h.reset()
	}

	if h.holdPhase {
		h.windowCount++
		// hold_for is finished
		if h.windowCount >= h.holdWindow {
			h.reset()
		} else {
			output = h.currentReading
		}
	}

	if !h.holdPhase {
		if input.Valid() {
			h.currentReading = input
			h.holdPhase = true
			h.windowCount = 0
			output = input
		}
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

func (h *Holder) reset() {
	h.currentReading = runtime.InvalidReading()
	h.holdPhase = false
	h.windowCount = 0
}

// DynamicConfigUpdate is a no-op for Holder.
func (*Holder) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
