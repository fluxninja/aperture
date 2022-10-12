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
	// ClassifiersPath is config path in etcd for classifiers.
	ClassifiersPath = path.Join(ConfigPrefix, "classifiers")
	// LoadShedDecisionsPath is decision path in etcd for load shed decisions.
	LoadShedDecisionsPath = path.Join(DecisionsPrefix, "load_shed")
	// ConcurrencyDemandDecisionsPath is decision path in etcd for concurrency demand decisions.
	ConcurrencyDemandDecisionsPath = path.Join(DecisionsPrefix, "concurrency_demand")
	// ConcurrencyDemandDefaultDecisionsPath is decision path in etcd for concurrency demand default decisions.
	ConcurrencyDemandDefaultDecisionsPath = path.Join(DecisionsPrefix, "concurrency_demand_default")
	// ConcurrencyMultiplierDecisionsPath is decision path in etcd for concurrency multiplier decisions.
	ConcurrencyMultiplierDecisionsPath = path.Join(DecisionsPrefix, "concurrency_multiplier")
	// AutoTokenResultsPath is config path in etcd for query tokens.
	AutoTokenResultsPath = path.Join(ConfigPrefix, "tokens")
	// ConcurrencyLimiterConfigPath is config path in etcd for concurrency limiter.
	ConcurrencyLimiterConfigPath = path.Join(ConfigPrefix, "concurrency_limiter")
	// RateLimiterConfigPath is config path in etcd for concurrency limiter.
	RateLimiterConfigPath = path.Join(ConfigPrefix, "rate_limiter")
	// RateLimiterDecisionsPath is decision path in etcd for rate limiter decisions.
	RateLimiterDecisionsPath = path.Join(DecisionsPrefix, "rate_limiter")
	// RateLimiterDynamicConfigPath is config path in etcd for dynamic config of rate limiter.
	RateLimiterDynamicConfigPath = path.Join(ConfigPrefix, "rate_limiter_dynamic_config")
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

// AgentPrefix returns the prefix for an agent.
func AgentPrefix(agentID string) string {
	return "agent-" + agentID
}

// DataplaneComponentKey returns the identifier for a Component in etcd.
func DataplaneComponentKey(agentGroupName, policyName string, componentIndex int64) string {
	return PolicyPrefix(agentGroupName, policyName) + "-component_index-" + strconv.FormatInt(componentIndex, 10)
}

// DataplaneComponentAgentKey returns the identifier for a Component for a particular Agent in etcd.
func DataplaneComponentAgentKey(agentID, agentGroupName, policyName string, componentIndex int64) string {
	return AgentPrefix(agentID) + "-" + DataplaneComponentKey(agentGroupName, policyName, componentIndex)
}

// FluxMeterKey returns the identifier for FluxMeter in etcd.
func FluxMeterKey(agentGroupName, fluxMeterName string) string {
	return AgentGroupPrefix(agentGroupName) + "-flux_meter-" + fluxMeterName
}

// ClassifierKey returns the identifier for a Classifier in etcd.
func ClassifierKey(agentGroupName, policyName string, classifierIndex int64) string {
	return PolicyPrefix(agentGroupName, policyName) + "-classifier_index-" + strconv.FormatInt(classifierIndex, 10)
}
