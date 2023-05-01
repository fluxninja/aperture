package iface

import (
	"github.com/prometheus/client_golang/prometheus"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

//go:generate mockgen -source=flux-meter.go -destination=../../mocks/mock_flux_meter.go -package=mocks

// FluxMeterID is the ID of the FluxMeter.
type FluxMeterID struct {
	FluxMeterName string
}

// String function returns the FluxMeterID as a string.
func (fmID FluxMeterID) String() string {
	return "flux_meter_name-" + fmID.FluxMeterName
}

// FluxMeter in an interface for interacting with fluxmeters.
type FluxMeter interface {
	// GetSelectors returns the selectors
	GetSelectors() []*policylangv1.Selector

	// GetAttributeKey returns the attribute key
	GetAttributeKey() string

	// GetFluxMeterName returns the metric name
	GetFluxMeterName() string

	// GetFluxMeterID returns the flux meter ID
	GetFluxMeterID() FluxMeterID

	// GetHistogram returns the histogram observer for given labels.
	// It expects the following labels to be set:
	//  * metrics.DecisionTypeLabel,
	//  * metrics.ResponseStatusLabel,
	//  * metrics.StatusCodeLabel,
	//  * metrics.FeatureStatusLabel.
	GetHistogram(labels map[string]string) prometheus.Observer

	// GetInvalidFluxMeterTotal returns a counter metric for the total number of invalid flux meters with the specified labels.
	GetInvalidFluxMeterTotal(labels map[string]string) (prometheus.Counter, error)
}
