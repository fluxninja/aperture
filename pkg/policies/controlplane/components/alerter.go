package components

import (
	"strconv"
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
	conditionFunc  AlerterConditionFunc
	labels         map[string]string
	annotations    map[string]string
}

// Make sure Alerter complies with Component interface.
var _ runtime.Component = (*Alerter)(nil)

// AlerterConditionFunc is a function that can be set in alerter to provide additional checks on signal values.
type AlerterConditionFunc func(runtime.Reading) bool

// NewAlerterAndOptions creates alerter and its fx options.
func NewAlerterAndOptions(alerterProto *policylangv1.Alerter, componentIdx int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	alerter := &Alerter{
		name:           alerterProto.AlertName,
		severity:       alerterProto.Severity,
		resolveTimeout: alerterProto.ResolveTimeout.AsDuration(),
		alertChannels:  make([]string, 0),
		policyReadAPI:  policyReadAPI,
		conditionFunc:  nil,
		labels: map[string]string{
			"type":            "alerter",
			"component_index": strconv.Itoa(componentIdx),
			"policy_name":     policyReadAPI.GetPolicyName(),
		},
		annotations: map[string]string{
			"resolve_timeout": alerterProto.ResolveTimeout.AsDuration().String(),
		},
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

	conditionRes := true
	if a.conditionFunc != nil {
		conditionRes = a.conditionFunc(signalValue)
	}

	if conditionRes {
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
		alerts.WithSeverity(a.severity),
		alerts.WithAnnotation("alert_channels", strings.Join(a.alertChannels, ",")),
	)

	for labelKey, labelVal := range a.labels {
		newAlert.SetLabel(labelKey, labelVal)
	}
	for annKey, annVal := range a.annotations {
		newAlert.SetAnnotation(annKey, annVal)
	}

	return newAlert
}

// SetConditionFunc allows to set condition function for alerter to check before creating alert.
func (a *Alerter) SetConditionFunc(cond AlerterConditionFunc) {
	a.conditionFunc = cond
}

// AddLabels allows to specify additional labels.
func (a *Alerter) AddLabels(newLabels map[string]string) {
	for labelKey, labelVal := range newLabels {
		a.labels[labelKey] = labelVal
	}
}

// AddAnnotations allows to specify additional annotations.
func (a *Alerter) AddAnnotations(newAnnotations map[string]string) {
	for annKey, annVal := range newAnnotations {
		a.annotations[annKey] = annVal
	}
}
