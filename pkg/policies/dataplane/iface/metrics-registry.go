package iface

import (
	"github.com/prometheus/client_golang/prometheus"
)

//go:generate mockgen -source=metrics-registry.go -destination=../../mocks/mock_metrics-registry.go -package=mocks

// ResponseMetricsAPI is an interface for getting response metrics.
type ResponseMetricsAPI interface {
	GetTokenLatencyHistogram(labels map[string]string) (prometheus.Observer, error)

	DeleteTokenLatencyHistogram(labels map[string]string) bool
}
