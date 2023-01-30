package paths

import (
	"path"
	"strconv"
)

var (
	// ConfigPrefix is key prefix in etcd for config.
	ConfigPrefix = path.Join("/config")
	// DecisionsPrefix is key prefix in etcd for decisions.
	DecisionsPrefix = path.Join("/decisions")
	// StatusPrefix is key prefix in etcd for status.
	StatusPrefix = path.Join("/status")
	// ClassifiersPath is config path in etcd for classifiers.
	ClassifiersPath = path.Join(ConfigPrefix, "classifiers")
	// LoadActuatorDecisionsPath is decision path in etcd for load decisions.
	LoadActuatorDecisionsPath = path.Join(DecisionsPrefix, "load")
	// AutoTokenResultsPath is config path in etcd for query tokens.
	AutoTokenResultsPath = path.Join(ConfigPrefix, "tokens")
	// ConcurrencyLimiterConfigPath is config path in etcd for concurrency limiter.
	ConcurrencyLimiterConfigPath = path.Join(ConfigPrefix, "concurrency_limiter")
	// RateLimiterConfigPath is config path in etcd for rate limiter.
	RateLimiterConfigPath = path.Join(ConfigPrefix, "rate_limiter")
	// RateLimiterDecisionsPath is decision path in etcd for rate limiter decisions.
	RateLimiterDecisionsPath = path.Join(DecisionsPrefix, "rate_limiter")
	// RateLimiterDynamicConfigPath is config path in etcd for dynamic config of rate limiter.
	RateLimiterDynamicConfigPath = path.Join(ConfigPrefix, "rate_limiter_dynamic_config")
	// FluxMeterConfigPath is config path in etcd for flux meters.
	FluxMeterConfigPath = path.Join(ConfigPrefix, "flux_meter")
	// HorizontalPodScalerConfigPath is config path in etcd for kubernetes horizontal pod scaler.
	HorizontalPodScalerConfigPath = path.Join(ConfigPrefix, "horizontal_pod_scaler")
	// HorizontalPodScalerDecisionsPath is decision path in etcd for kubernetes horizontal pod scaler decisions.
	HorizontalPodScalerDecisionsPath = path.Join(DecisionsPrefix, "horizontal_pod_scaler")
	// HorizontalPodScalerStatusPath is decision path in etcd for kubernetes horizontal pod scaler status.
	HorizontalPodScalerStatusPath = path.Join(StatusPrefix, "horizontal_pod_scaler")
	// HorizontalPodScalerDynamicConfigPath is config path in etcd for dynamic config of kubernetes horizontal pod scaler.
	HorizontalPodScalerDynamicConfigPath = path.Join(ConfigPrefix, "horizontal_pod_scaler_dynamic_config")
)

// AgentGroupPrefix returns the prefix for an agent group.
func AgentGroupPrefix(agentGroupName string) string {
	return "agent_group-" + agentGroupName
}

// PolicyPrefix returns the prefix for a policy.
func PolicyPrefix(agentGroupName, policyName string) string {
	return AgentGroupPrefix(agentGroupName) + "-policy-" + policyName
}

// AgentComponentKey returns the identifier for a Component in etcd.
func AgentComponentKey(agentGroupName, policyName string, componentID string) string {
	return PolicyPrefix(agentGroupName, policyName) + "-component_id-" + componentID
}

// FluxMeterKey returns the identifier for FluxMeter in etcd.
func FluxMeterKey(agentGroupName, fluxMeterName string) string {
	return AgentGroupPrefix(agentGroupName) + "-flux_meter-" + fluxMeterName
}

// ClassifierKey returns the identifier for a Classifier in etcd.
func ClassifierKey(agentGroupName, policyName string, classifierIndex int64) string {
	return PolicyPrefix(agentGroupName, policyName) + "-classifier_index-" + strconv.FormatInt(classifierIndex, 10)
}
