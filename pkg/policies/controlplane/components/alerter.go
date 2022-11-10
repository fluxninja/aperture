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

// Alerter is a component that monitors signal value and creates alert on true value.
type Alerter struct {
	name           string
	severity       string
	resolveTimeout time.Duration
	alertChannels  []string
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
	}
	alerter.alertChannels = append(alerter.alertChannels, alerterProto.AlertChannels...)

	return alerter, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (alerter *Alerter) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	signalValue := inPortReadings.ReadSingleValuePort("alert")
	if !signalValue.Valid() {
		return nil, nil
	}

	return nil, nil
}

// DynamicConfigUpdate is a no-op for Alerter.
func (alerter *Alerter) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
