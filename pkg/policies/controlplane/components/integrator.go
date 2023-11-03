package components

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime/tristate"
)

// Integrator is a component that accumulates sum of signal every tick.
type Integrator struct {
	sum               float64
	doExecute         bool
	ticksPerExecution int
	cpID              string
}

// Make sure Integrator complies with Component interface.
var _ runtime.Component = (*Integrator)(nil)

// Make sure Integrator implements background job.
var _ runtime.BackgroundJob = (*Integrator)(nil)

// Name implements runtime.Component.
func (*Integrator) Name() string { return "Integrator" }

// Type implements runtime.Component.
func (*Integrator) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (in *Integrator) ShortDescription() string {
	return ""
}

// IsActuator implements runtime.Component.
func (*Integrator) IsActuator() bool { return false }

// NewIntegrator creates an integrator component.

// NewIntegratorAndOptions creates an integrator component and its fx options.
func NewIntegratorAndOptions(integratorProto *policylangv1.Integrator, componentID runtime.ComponentID, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	initialValue := integratorProto.GetInitialValue()
	integrator := &Integrator{
		sum:               initialValue,
		cpID:              componentID.String(),
		ticksPerExecution: policyReadAPI.TicksInDurationPb(integratorProto.EvaluationInterval),
	}

	return integrator, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (in *Integrator) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	circuitAPI.ScheduleConditionalBackgroundJob(in, in.ticksPerExecution)

	if !in.doExecute {
		return runtime.PortToReading{
			"output": []runtime.Reading{runtime.NewReading(in.sum)},
		}, nil
	}
	in.doExecute = false
	inputVal := inPortReadings.ReadSingleReadingPort("input")
	resetVal := inPortReadings.ReadSingleReadingPort("reset")
	if tristate.FromReading(resetVal).IsTrue() {
		in.sum = 0
	} else if inputVal.Valid() {
		in.sum += inputVal.Value()

		maxVal := inPortReadings.ReadSingleReadingPort("max")
		if maxVal.Valid() {
			in.sum = math.Min(in.sum, maxVal.Value())
		}

		minVal := inPortReadings.ReadSingleReadingPort("min")
		if minVal.Valid() {
			in.sum = math.Max(in.sum, minVal.Value())
		}
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{runtime.NewReading(in.sum)},
	}, nil
}

// DynamicConfigUpdate is a no-op for Integrator.
func (in *Integrator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}

// GetJob implements runtime.BackgroundJob.GetJob.
func (in *Integrator) GetJob() jobs.Job {
	return jobs.NewNoOpJob(in.cpID)
}

// NotifyCompletion implements runtime.BackgroundJob.NotifyCompletion.
func (in *Integrator) NotifyCompletion() {
	in.doExecute = true
}
