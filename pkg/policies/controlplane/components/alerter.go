package components

import (
	"fmt"
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime/tristate"
)

// Alerter is a component that monitors signal value and creates alert on true value.
type Alerter struct {
	alerterIface   alerts.Alerter
	policyReadAPI  iface.Policy
	labels         map[string]string
	name           string
	severity       string
	componentID    string
	alertChannels  []string
	resolveTimeout time.Duration
}

// Name implements runtime.Component.
func (*Alerter) Name() string { return "Alerter" }

// Type implements runtime.Component.
func (*Alerter) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (a *Alerter) ShortDescription() string { return fmt.Sprintf("%s/%s", a.name, a.severity) }

// IsActuator implements runtime.Component.
func (*Alerter) IsActuator() bool { return false }

// Make sure Alerter complies with Component interface.
var _ runtime.Component = (*Alerter)(nil)

// NewAlerterAndOptions creates alerter and its fx options.
func NewAlerterAndOptions(alerterProto *policylangv1.Alerter, componentID runtime.ComponentID, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	parameters := alerterProto.Parameters
	alerter := &Alerter{
		name:           parameters.AlertName,
		severity:       parameters.Severity,
		resolveTimeout: parameters.ResolveTimeout.AsDuration(),
		alertChannels:  make([]string, 0),
		policyReadAPI:  policyReadAPI,
		componentID:    componentID.String(),
		labels:         make(map[string]string),
	}
	alerter.alertChannels = append(alerter.alertChannels, parameters.AlertChannels...)
	if alerterProto.Parameters.Labels != nil {
		alerter.labels = alerterProto.Parameters.Labels
	}

	return alerter, fx.Options(
		fx.Invoke(
			alerter.setup,
		)), nil
}

func (a *Alerter) setup(alerterIface *alerts.SimpleAlerter) {
	a.alerterIface = alerterIface
}

// Execute implements runtime.Component.Execute.
func (a *Alerter) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	signalValue := inPortReadings.ReadSingleReadingPort("signal")

	if tristate.FromReading(signalValue).IsTrue() {
		a.alerterIface.AddAlert(a.createAlert())
	}

	return nil, nil
}

// DynamicConfigUpdate is a no-op for Alerter.
func (a *Alerter) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (a *Alerter) createAlert() *alerts.Alert {
	newAlert := alerts.NewAlert(
		alerts.WithName(a.name),
		alerts.WithSeverity(alerts.ParseSeverity(a.severity)),
		alerts.WithAlertChannels(a.alertChannels),
		alerts.WithLabel("policy_name", a.policyReadAPI.GetPolicyName()),
		alerts.WithLabel("type", "alerter"),
		alerts.WithLabel("component_id", a.componentID),
		alerts.WithResolveTimeout(a.resolveTimeout),
		alerts.WithGeneratorURL(
			fmt.Sprintf("http://%s/%s/%s", info.GetHostInfo().Hostname, a.policyReadAPI.GetPolicyName(), a.name),
		),
	)

	for key, val := range a.labels {
		newAlert.SetLabel(key, val)
	}

	return newAlert
}
