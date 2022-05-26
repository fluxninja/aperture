package policyapi

import (
	"time"

	"github.com/prometheus/prometheus/model/labels"
)

const (
	// PoliciesStatusRoot - path in status registry for policies results.
	PoliciesStatusRoot = "policies"
)

// PolicyBaseAPI is for read only access to base policy info.
type PolicyBaseAPI interface {
	GetPolicyName() string
	GetPolicyHash() string
}

// PolicyReadAPI is for read only access to full policy state.
type PolicyReadAPI interface {
	PolicyBaseAPI
	ResolveMetricNames(query string) (string, error)
	GetEvaluationInterval() time.Duration
}

// MetricSubRegistry is for registering metric substitution patterns (used by FluxMeter).
type MetricSubRegistry interface {
	RegisterHistogramSub(metricNameOrig, metricNameSub string, labelMatchers []*labels.Matcher)
	RegisterMetricSub(metricsNameOrig, metricNameSub string, labelMatchers []*labels.Matcher)
}

// PolicyAPI is the global interface composed of all of the above APIs.
type PolicyAPI interface {
	PolicyReadAPI
	MetricSubRegistry
}
