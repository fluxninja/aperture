package alerts

import (
	"strings"

	"github.com/fluxninja/aperture/v2/pkg/config"
)

// AlertsFxTag - name tag for alerter in fx.
var AlertsFxTag = config.NameTag("AlertsFx")

// Alerter is responsible for receiving alerts and propagating them to the channel
// returned by AlertsChan().
type Alerter interface {
	AddAlert(*Alert)
	AlertsChan() <-chan *Alert
	WithLabels(map[string]string) Alerter
}

// SimpleAlerter implements Alerter interface. It just simple propagates alerts
// to the channel.
type SimpleAlerter struct {
	alertsCh chan *Alert
}

// NewSimpleAlerter returns new instance of SimpleAlerter with channel of given size.
func NewSimpleAlerter(channelSize int) Alerter {
	return &SimpleAlerter{
		alertsCh: make(chan *Alert, channelSize),
	}
}

// AddAlert adds alert to the channel.
func (a *SimpleAlerter) AddAlert(alert *Alert) {
	a.alertsCh <- alert
}

// AlertsChan returns the alerts channel.
func (a *SimpleAlerter) AlertsChan() <-chan *Alert {
	return a.alertsCh
}

// WithLabels returns the alerter wrapper with specified labels.
func (a *SimpleAlerter) WithLabels(labels map[string]string) Alerter {
	return newAlerterWrapper(a, sanitizeKeysInLabels(labels))
}

type alerterWrapper struct {
	parentAlerter Alerter
	labels        map[string]string
}

func newAlerterWrapper(parent Alerter, labels map[string]string) Alerter {
	return &alerterWrapper{
		parentAlerter: parent,
		labels:        labels,
	}
}

// AddAlert adds alert to the channel with labels specified at wrapper creation.
func (aw *alerterWrapper) AddAlert(alert *Alert) {
	for key, val := range aw.labels {
		alert.SetLabel(key, val)
	}
	aw.parentAlerter.AddAlert(alert)
}

// AlertsChan returns the alerts channel.
func (aw *alerterWrapper) AlertsChan() <-chan *Alert {
	return aw.parentAlerter.AlertsChan()
}

// WithLabels returns the alerter with new labels added to parent labels.
func (aw *alerterWrapper) WithLabels(labels map[string]string) Alerter {
	mergedLabels := make(map[string]string)
	for k, v := range aw.labels {
		mergedLabels[k] = v
	}
	// overwrite older values with new ones
	for k, v := range labels {
		mergedLabels[k] = v
	}
	return newAlerterWrapper(aw, sanitizeKeysInLabels(mergedLabels))
}

func sanitizeKeysInLabels(labels map[string]string) map[string]string {
	// change '-' to '_' in key because alertmanager does not accept it
	fixedLabels := make(map[string]string)
	for key, val := range labels {
		fixedKey := strings.ReplaceAll(key, "-", "_")
		fixedLabels[fixedKey] = val
	}
	return fixedLabels
}
