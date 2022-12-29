package metricsprocessor

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/fluxninja/aperture/pkg/controlpointcache"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
)

// Config holds configuration for the metrics processor.
type Config struct {
	promRegistry         *prometheus.Registry
	engine               iface.Engine
	classificationEngine iface.ClassificationEngine
	controlPointCache    *controlpointcache.ControlPointCache
}
