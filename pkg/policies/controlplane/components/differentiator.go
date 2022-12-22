package components

import (
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Differentiator is a component that accumulates sum of signal every tick.
type Differentiator struct {
	window time.Duration
	// readings are saved in ring buffer
	readings  []point
	oldestIdx int
	newestIdx int
	// capacity is calculated from window duration divided by tick interval
	capacity    int
	initialized bool
}

type point struct {
	reading   runtime.Reading
	timestamp time.Time
}

// Name implements runtime.Component.
func (*Differentiator) Name() string { return "Differentiator" }

// Type implements runtime.Component.
func (*Differentiator) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

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

// NewDifferentiatorAndOptions creates an integrator component and its fx options.
func NewDifferentiatorAndOptions(diffProto *policylangv1.Differentiator, _ int, _ iface.Policy) (runtime.Component, fx.Option, error) {
	return NewDifferentiator(diffProto), fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (d *Differentiator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	if !d.initialized {
		d.init(tickInfo)
	}

	inputVal := inPortReadings.ReadSingleValuePort("input")
	outputVal := runtime.NewReading(0)

	// add input to readings array
	d.readings[d.newestIdx] = point{reading: inputVal, timestamp: tickInfo.Timestamp()}

	if d.oldestIdx != d.newestIdx {
		oldest := d.readings[d.oldestIdx]
		newest := d.readings[d.newestIdx]

		if !oldest.reading.Valid() {
			oldest = d.firstValid(d.oldestIdx, true)
		}

		// find the newest valid reading for extrapolation
		if !newest.reading.Valid() {
			found := d.firstValid(d.newestIdx, false)
			if found.reading.Valid() && oldest.reading.Valid() {
				extrapolatedValue := d.extrapolate(oldest, found, tickInfo.Timestamp())
				newest = point{reading: runtime.NewReading(extrapolatedValue), timestamp: tickInfo.Timestamp()}
			}
		}

		// calculate the derivative
		if oldest.reading.Valid() && newest.reading.Valid() {
			diff := (newest.reading.Value() - oldest.reading.Value()) /
				float64(newest.timestamp.UnixMilli()-oldest.timestamp.UnixMilli())
			outputVal = runtime.NewReading(diff)
		}
	}

	// shift pointers for newest and oldest
	d.newestIdx = (d.newestIdx + 1) % d.capacity
	if d.newestIdx == d.oldestIdx {
		d.oldestIdx = (d.oldestIdx + 1) % d.capacity
	}

	return runtime.PortToValue{
		"output": []runtime.Reading{outputVal},
	}, nil
}

// DynamicConfigUpdate is a no-op for Differentiator.
func (d *Differentiator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (d *Differentiator) extrapolate(firstVal, secondVal point, currentTime time.Time) float64 {
	extValue := firstVal.reading.Value() + float64(currentTime.UnixMilli()-firstVal.timestamp.UnixMilli())/
		float64(secondVal.timestamp.UnixMilli()-firstVal.timestamp.UnixMilli())*
		(secondVal.reading.Value()-firstVal.reading.Value())
	return extValue
}

func (d *Differentiator) init(tickInfo runtime.TickInfo) {
	d.capacity = int(d.window / tickInfo.Interval())
	d.readings = make([]point, d.capacity)
	for i := 0; i < d.capacity; i++ {
		d.readings[i].reading = runtime.InvalidReading()
		d.readings[i].timestamp = time.Now()
	}
	d.initialized = true
}

func (d *Differentiator) firstValid(fromIdx int, addition bool) point {
	step := 1
	if !addition {
		step = -1
	}

	idx := (fromIdx + step) % d.capacity
	if (fromIdx + step) < 0 {
		// ring buffer wrap up
		idx = d.capacity - 1
	}

	for idx != fromIdx {
		if d.readings[idx].reading.Valid() {
			return d.readings[idx]
		}
		idx = (idx + step) % d.capacity
		if (idx + step) < 0 {
			// ring buffer wrap up
			idx = d.capacity - 1
		}
	}
	return point{reading: runtime.InvalidReading(), timestamp: time.Now()}
}
