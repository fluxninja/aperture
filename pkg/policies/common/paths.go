package common

import (
	"path"
	"strconv"
)

var (
	// ConfigPrefix is key prefix in etcd for config.
	ConfigPrefix = path.Join("/config")
	// DecisionsPrefix is key prefix in etcd for decisions.
	DecisionsPrefix = path.Join("/decisions")
	// ClassifiersConfigPath is config path in etcd for classifiers.
	ClassifiersConfigPath = path.Join(ConfigPrefix, "classifiers")
	// PoliciesConfigPath is config path in etcd for policies.
	PoliciesConfigPath = path.Join(ConfigPrefix, "policies")
	// LoadShedDecisionsPath is decision path in etcd for load shed decisions.
	LoadShedDecisionsPath = path.Join(DecisionsPrefix, "load_shed")
	// AutoTokenResultsPath is config path in etcd for query tokens.
	AutoTokenResultsPath = path.Join(ConfigPrefix, "tokens")
	// ConcurrencyLimiterConfigPath is config path in etcd for concurrency limiter.
	ConcurrencyLimiterConfigPath = path.Join(ConfigPrefix, "concurrency_limiter")
	// RateLimiterConfigPath is config path in etcd for concurrency limiter.
	RateLimiterConfigPath = path.Join(ConfigPrefix, "rate_limiter")
	// RateLimiterDecisionsPath is decision path in etcd for rate limiter decisions.
	RateLimiterDecisionsPath = path.Join(DecisionsPrefix, "rate_limiter")
	// FluxMeterConfigPath is config path in etcd for flux meters.
	FluxMeterConfigPath = path.Join(ConfigPrefix, "flux_meter")
)

// AgentGroupPrefix returns the prefix for an agent group.
func AgentGroupPrefix(agentGroupName string) string {
	return "agent_group-" + agentGroupName
}

// PolicyPrefix returns the prefix for a policy.
func PolicyPrefix(agentGroupName, policyName string) string {
	return AgentGroupPrefix(agentGroupName) + "-policy-" + policyName
}

// DataplaneComponentKey returns the identifier for a Component in etcd.
func DataplaneComponentKey(agentGroupName, policyName string, componentIndex int64) string {
	return PolicyPrefix(agentGroupName, policyName) + "-component_index-" + strconv.FormatInt(componentIndex, 10)
}

// FluxMeterKey returns the identifier for FluxMeter in etcd.
func FluxMeterKey(agentGroupName, fluxMeterName string) string {
	return AgentGroupPrefix(agentGroupName) + "-flux_meter-" + fluxMeterName
}

// ClassifierKey returns the identifier for a Classifier in etcd.
func ClassifierKey(agentGroupName, policyName string, classifierIndex int64) string {
	return PolicyPrefix(agentGroupName, policyName) + "-classifier_index-" + strconv.FormatInt(classifierIndex, 10)
}
