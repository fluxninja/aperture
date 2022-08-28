package iface

import (
	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
)

//go:generate mockgen -source=metrics-registry.go -destination=../../mocks/mock_metrics-registry.go -package=mocks

// ResponseMetricsAPI is an interface for getting response metrics.
type ResponseMetricsAPI interface {
	GetFluxMeterHistogram(fluxmeterID, statusCode, featureStatus string, decisionType flowcontrolv1.CheckResponse_DecisionType) (prometheus.Observer, error)
	GetTokenLatencyHistogram(labels map[string]string) (prometheus.Observer, error)

	DeleteFluxmeterHistogram(fluxmeterID string) bool
	DeleteTokenLatencyHistogram(labels map[string]string) bool
}
