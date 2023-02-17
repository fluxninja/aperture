package servicegetter

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/fluxninja/aperture/pkg/metrics"
)

// Metrics is used for collecting metrics about servicegetter
//
// nil value of Metrics should be always usable.
type Metrics struct {
	okTotal     prometheus.Counter
	errorsTotal prometheus.Counter
}

// NewMetrics creates new Metrics, registering counters in given registry.
func NewMetrics(registry *prometheus.Registry) (*Metrics, error) {
	total := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.ServiceLookupsMetricName,
			Help: "Number of IP to services lookups",
		},
		[]string{metrics.ServiceLookupsStatusLabel},
	)
	if err := registry.Register(total); err != nil {
		return nil, err
	}
	return &Metrics{
		okTotal:     total.WithLabelValues(metrics.ServiceLookupsStatusOK),
		errorsTotal: total.WithLabelValues(metrics.ServiceLookupsStatusError),
	}, nil
}

func (m *Metrics) inc(ok bool) {
	if m == nil {
		return
	}
	if ok {
		m.okTotal.Inc()
	} else {
		m.errorsTotal.Inc()
	}
}
