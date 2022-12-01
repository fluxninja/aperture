package concurrency

import (
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// LoadActuatorAlerter is a component that extends basic alerter.
type LoadActuatorAlerter struct {
	alerter runtime.Component
}

// Make sure LoadActuatorAlerter complies with Component interface.
var _ runtime.Component = (*LoadActuatorAlerter)(nil)

// NewLoadActuatorAlerterAndOptions creates load actuator alerter and its fx options.
func NewLoadActuatorAlerterAndOptions(
	loadActuatorAlerterProto *policylangv1.LoadActuatorAlerter,
	componentIndex int,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	alerter, options, err := components.NewAlerterAndOptions(loadActuatorAlerterProto.BaseAlerter, componentIndex, policyReadAPI)
	if err != nil {
		log.Error().Msg("Could not create load actuator alerter")
		return nil, nil, err
	}

	// TODO get controllerID

	alerterObj, ok := alerter.(*components.Alerter)
	if ok {
		labels := map[string]string{
			"type":       "concurrency_limiter",
			"alert_name": "Load Shed Event",
		}

		durationProto := loadActuatorAlerterProto.BaseAlerter.ResolveTimeout.AsDuration()
		timeout, _ := time.ParseDuration("5s")
		if durationProto > timeout {
			timeout = durationProto
		}
		annotations := map[string]string{
			"resolve_timeout": timeout.String(),
		}

		alerterObj.SetConditionFunc(loadActuatorCondition)
		alerterObj.AddLabels(labels)
		alerterObj.AddAnnotations(annotations)
	}

	laAlerter := &LoadActuatorAlerter{
		alerter: alerter,
	}

	return laAlerter, fx.Options(options), nil
}

// Execute implements runtime.Component.Execute.
func (a *LoadActuatorAlerter) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	return a.alerter.Execute(inPortReadings, tickInfo)
}

// DynamicConfigUpdate is a no-op for LoadActuatorAlerter.
func (a *LoadActuatorAlerter) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func loadActuatorCondition(signalReading runtime.Reading) bool {
	return signalReading.Value() < 1
}
