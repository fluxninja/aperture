package iface

import (
	"time"

	"github.com/prometheus/prometheus/model/labels"
)

// swagger:operation POST /policies common-configuration PoliciesConfig
// ---
// x-fn-config-env: true
// parameters:
// - name: promql_jobs_scheduler
//   in: body
//   schema:
//     "$ref": "#/definitions/JobGroupConfig"

const (
	// PoliciesRoot - path in config and status registry for policies results.
	PoliciesRoot = "policies"
)

// PolicyBase is for read only access to base policy info.
type PolicyBase interface {
	GetPolicyName() string
	GetPolicyHash() string
}

// PolicyRead is for read only access to full policy state.
type PolicyRead interface {
	PolicyBase
	ResolveMetricNames(query string) (string, error)
	GetEvaluationInterval() time.Duration
}

// MetricSubRegistry is for registering metric substitution patterns (used by FluxMeter).
type MetricSubRegistry interface {
	RegisterHistogramSub(metricNameOrig, metricNameSub string, labelMatchers []*labels.Matcher)
	RegisterMetricSub(metricsNameOrig, metricNameSub string, labelMatchers []*labels.Matcher)
}

// Policy is the global interface composed of all of the above APIs.
type Policy interface {
	PolicyRead
	MetricSubRegistry
}
