package components

import (
	"strings"
	"time"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Alerter is a component that monitors signal value and creates alert on true value.
type Alerter struct {
	name           string
	severity       string
	resolveTimeout time.Duration
	alertChannels  []string
	alerterIface   alerts.Alerter
	policyReadAPI  iface.Policy
}

// Make sure Alerter complies with Component interface.
var _ runtime.Component = (*Alerter)(nil)

// NewAlerterAndOptions creates alerter and its fx options.
func NewAlerterAndOptions(alerterProto *policylangv1.Alerter, _ int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	alerter := &Alerter{
		name:           alerterProto.AlertName,
		severity:       alerterProto.Severity,
		resolveTimeout: alerterProto.ResolveTimeout.AsDuration(),
		alertChannels:  make([]string, 0),
		policyReadAPI:  policyReadAPI,
	}
	alerter.alertChannels = append(alerter.alertChannels, alerterProto.AlertChannels...)

	return alerter, fx.Options(
		fx.Invoke(
			alerter.setup,
		)), nil
}

func (a *Alerter) setup(alerterIface *alerts.SimpleAlerter) {
	a.alerterIface = alerterIface
}

// Execute implements runtime.Component.Execute.
func (a *Alerter) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	signalValue := inPortReadings.ReadSingleValuePort("alert")
	if !signalValue.Valid() {
		return nil, nil
	}

	a.alerterIface.AddAlert(a.createAlert())

	return nil, nil
}

// DynamicConfigUpdate is a no-op for Alerter.
func (a *Alerter) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (a *Alerter) createAlert() *alerts.Alert {
	newAlert := alerts.NewAlert()
	newAlert.SetName(a.name)
	newAlert.SetSeverity(a.severity)
	newAlert.SetLabel("policy_name", a.policyReadAPI.GetPolicyName())
	newAlert.SetLabel("type", "alerter")
	newAlert.SetAnnotations(map[string]string{
		"alert_channels":  strings.Join(a.alertChannels, ","),
		"resolve_timeout": a.resolveTimeout.String(),
	})
	return newAlert
}
