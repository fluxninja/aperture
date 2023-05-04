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
	"github.com/fluxninja/aperture/pkg/utils"
)

// Differentiator is a component that calculates rate of change per tick.
type Differentiator struct {
	// readings are saved in ring buffer
	readings  []runtime.Reading
	window    time.Duration
	oldestIdx int
	newestIdx int
	// capacity is calculated from window duration divided by tick interval
	capacity    int
	initialized bool
}

// Name implements runtime.Component.
func (*Differentiator) Name() string { return "Differentiator" }

// Type implements runtime.Component.
func (*Differentiator) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (d *Differentiator) ShortDescription() string { return fmt.Sprintf("win: %s", d.window) }

// IsActuator implements runtime.Component.
func (*Differentiator) IsActuator() bool { return false }

// NewDifferentiator creates a differentiator component.
func NewDifferentiator(diffProto *policylangv1.Differentiator) runtime.Component {
	diff := &Differentiator{
		window:      diffProto.Window.AsDuration(),
		oldestIdx:   0,
		newestIdx:   0,
		capacity:    0,
		initialized: false,
	}
	return diff
}

// NewDifferentiatorAndOptions creates a differentiator component and its fx options.
func NewDifferentiatorAndOptions(diffProto *policylangv1.Differentiator, _ runtime.ComponentID, _ iface.Policy) (runtime.Component, fx.Option, error) {
	return NewDifferentiator(diffProto), fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (d *Differentiator) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	if !d.initialized {
		d.init(tickInfo)
	}

	inputVal := inPortReadings.ReadSingleReadingPort("input")
	outputVal := runtime.InvalidReading()

	// add input to readings array
	d.readings[d.newestIdx] = inputVal

	if d.oldestIdx != d.newestIdx {
		oldest := d.readings[d.oldestIdx]
		newest := d.readings[d.newestIdx]
		oldestIdx := d.oldestIdx
		newestIdx := d.newestIdx

		if !oldest.Valid() {
			oldest, oldestIdx = d.firstValid(d.oldestIdx, true)
		}

		// find the newest valid reading for extrapolation
		if !newest.Valid() {
			found, foundIdx := d.firstValid(d.newestIdx, false)
			if found.Valid() && oldest.Valid() && foundIdx != oldestIdx {
				extrapolatedValue := d.extrapolate(oldest, found, oldestIdx, foundIdx)
				newest = runtime.NewReading(extrapolatedValue)
				newestIdx = foundIdx
			}
		}

		// calculate the derivative
		if oldest.Valid() && newest.Valid() && newestIdx != oldestIdx {
			diff := (newest.Value() - oldest.Value()) / float64(d.capacity)
			outputVal = runtime.NewReading(diff)
		}
	}

	// shift pointers for newest and oldest
	d.newestIdx = utils.Mod((d.newestIdx + 1), d.capacity)
	if d.newestIdx == d.oldestIdx {
		d.oldestIdx = utils.Mod((d.oldestIdx + 1), d.capacity)
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{outputVal},
	}, nil
}

// DynamicConfigUpdate is a no-op for Differentiator.
func (d *Differentiator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (d *Differentiator) extrapolate(firstVal, secondVal runtime.Reading, firstIdx, secondIdx int) float64 {
	extrapolatedIdx := utils.Mod(secondIdx+1, d.capacity)
	extValue := firstVal.Value() + float64(extrapolatedIdx-firstIdx)/
		float64(secondIdx-firstIdx)*
		(secondVal.Value()-firstVal.Value())
	return extValue
}

func (d *Differentiator) init(tickInfo runtime.TickInfo) {
	d.capacity = int(math.Ceil(float64(d.window) / float64(tickInfo.Interval())))
	d.readings = make([]runtime.Reading, d.capacity)
	for i := 0; i < d.capacity; i++ {
		d.readings[i] = runtime.InvalidReading()
	}
	d.initialized = true
}

func (d *Differentiator) firstValid(fromIdx int, addition bool) (runtime.Reading, int) {
	step := 1
	if !addition {
		step = -1
	}

	idx := utils.Mod((fromIdx + step), d.capacity)
	for i := 0; i < d.capacity; i++ {
		if d.readings[idx].Valid() {
			return d.readings[idx], idx
		}
		idx = utils.Mod((idx + step), d.capacity)
	}
	return runtime.InvalidReading(), idx
}
