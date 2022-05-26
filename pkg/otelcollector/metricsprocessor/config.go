package metricsprocessor

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/config"

	"aperture.tech/aperture/pkg/policies/dataplane/iface"
)

// Config holds configuration for the metrics processor.
type Config struct {
	promRegistry             *prometheus.Registry
	engine                   iface.EngineAPI
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	// The lowest bucket in latency histogram
	LatencyBucketStartMS float64 `mapstructure:"latency_bucket_start_ms" validate:"gte=0" default:"20"`
	// The bucket width in latency histogram
	LatencyBucketWidthMS float64 `mapstructure:"latency_bucket_width_ms" validate:"gte=0" default:"20"`
	// The number of buckets in latency histogram
	LatencyBucketCount int `mapstructure:"latency_bucket_count" validate:"gte=0" default:"100"`
}
