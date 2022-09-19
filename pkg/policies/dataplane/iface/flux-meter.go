package iface

import (
	"github.com/prometheus/client_golang/prometheus"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
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
	// GetSelector returns the selector
	GetSelector() *selectorv1.Selector

	// GetAttributeKey returns the attribute key
	GetAttributeKey() string

	// GetFluxMeterName returns the metric name
	GetFluxMeterName() string

	// GetFluxMeterID returns the flux meter ID
	GetFluxMeterID() FluxMeterID

	// GetHistogram returns the histogram observer for the flowcontrolv1.DecisionType
	GetHistogram(decisionType flowcontrolv1.CheckResponse_DecisionType, statusCode string, featureStatus string) prometheus.Observer
}
