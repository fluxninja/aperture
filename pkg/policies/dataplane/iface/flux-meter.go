package iface

import (
	"github.com/prometheus/client_golang/prometheus"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

//go:generate mockgen -source=flux-meter.go -destination=../../mocks/mock_flux_meter.go -package=mocks

// FluxMeterID is the ID of the FluxMeter.
type FluxMeterID struct {
	PolicyName    string
	FluxMeterName string
	PolicyHash    string
}

// String function returns the FluxMeterID as a string.
func (fmID FluxMeterID) String() string {
	return "policy_name-" + fmID.PolicyName + "-flux_meter_name-" + fmID.FluxMeterName + "-policy_hash-" + fmID.PolicyHash
}

// FluxMeter in an interface for interacting with fluxmeters.
type FluxMeter interface {
	// Policy
	GetPolicyName() string
	GetPolicyHash() string

	// GetSelector returns the selector
	GetSelector() *selectorv1.Selector

	// GetFluxMeterProto returns the flux meter proto
	GetFluxMeterProto() *policylangv1.FluxMeter

	// GetFluxMeterName returns the metric name
	GetFluxMeterName() string

	// GetFluxMeterID returns the flux meter ID
	GetFluxMeterID() FluxMeterID

	// GetBuckets returns the buckets
	GetBuckets() []float64

	// GetHistogram returns the histogram for the flowcontrolv1.DecisionType
	GetHistogram(flowcontrolv1.DecisionType) prometheus.Histogram
}
