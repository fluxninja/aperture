package metricsprocessor

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/fluxninja/aperture/v2/pkg/cache"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
)

// Config holds configuration for the metrics processor.
type Config struct {
	promRegistry         *prometheus.Registry
	engine               iface.Engine
	classificationEngine iface.ClassificationEngine
	controlPointCache    *cache.Cache[selectors.TypedControlPointID]
}
