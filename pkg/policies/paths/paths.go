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
	// PoliciesAPIConfigPath is config path in etcd for policies via API.
	PoliciesAPIConfigPath = path.Join(ConfigPrefix, "api", "policies")
	// PoliciesAPIDynamicConfigPath is config path in etcd for  dynamic configuration of policies via API.
	PoliciesAPIDynamicConfigPath = path.Join(ConfigPrefix, "api", "dynamic-config-policies")
	// PoliciesConfigPath is config path in etcd for policies.
	PoliciesConfigPath = path.Join(ConfigPrefix, "policies")
	// PoliciesDynamicConfigPath is config path in etcd for dynamic configuration of policies.
	PoliciesDynamicConfigPath = path.Join(ConfigPrefix, "dynamic-config-policies")
	// ClassifiersPath is config path in etcd for classifiers.
	ClassifiersPath = path.Join(ConfigPrefix, "classifiers")
	// LoadSchedulerConfigPath is config path in etcd for load scheduler.
	LoadSchedulerConfigPath = path.Join(ConfigPrefix, "load_scheduler")
	// LoadSchedulerDecisionsPath is decision path in etcd for load decisions.
	LoadSchedulerDecisionsPath = path.Join(DecisionsPrefix, "load_scheduler")
	// RateLimiterConfigPath is config path in etcd for rate limiter.
	RateLimiterConfigPath = path.Join(ConfigPrefix, "rate_limiter")
	// RateLimiterDecisionsPath is decision path in etcd for rate limiter decisions.
	RateLimiterDecisionsPath = path.Join(DecisionsPrefix, "rate_limiter")
	// QuotaSchedulerConfigPath is config path in etcd for quota scheduler.
	QuotaSchedulerConfigPath = path.Join(ConfigPrefix, "quota_scheduler")
	// QuotaSchedulerDecisionsPath is decision path in etcd for quota scheduler decisions.
	QuotaSchedulerDecisionsPath = path.Join(DecisionsPrefix, "quota_scheduler")
	// FluxMeterConfigPath is config path in etcd for flux meters.
	FluxMeterConfigPath = path.Join(ConfigPrefix, "flux_meter")
	// InfraMeterConfigPath is config path in etcd for infra meter.
	InfraMeterConfigPath = path.Join(ConfigPrefix, "infra_meter")
	// PodScalerConfigPath is config path in etcd for pod scaler.
	PodScalerConfigPath = path.Join(ConfigPrefix, "pod_scaler")
	// PodScalerDecisionsPath is decision path in etcd for pod scaler decisions.
	PodScalerDecisionsPath = path.Join(DecisionsPrefix, "pod_scaler")
	// PodScalerStatusPath is decision path in etcd for pod scaler status.
	PodScalerStatusPath = path.Join(StatusPrefix, "pod_scaler")
	// SamplerConfigPath is config path in etcd for sampler.
	SamplerConfigPath = path.Join(ConfigPrefix, "sampler")
	// SamplerDecisionsPath is decision path in etcd for sampler decisions.
	SamplerDecisionsPath = path.Join(DecisionsPrefix, "sampler")
	// PeersPrefix is the prefix for peers path in etcd.
	PeersPrefix = path.Join("/peers")
	// AgentPeerPath is path in etcd for agent peers.
	AgentPeerPath = path.Join(PeersPrefix, "aperture-agent")
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

// InfraMeterKey returns the identifier for InfraMeter in etcd.
func InfraMeterKey(agentGroupName, policyName, infraMeterName string) string {
	return PolicyPrefix(agentGroupName, policyName) + "-infra_meter-" + infraMeterName
}

// ClassifierKey returns the identifier for a Classifier in etcd.
func ClassifierKey(agentGroupName, policyName string, classifierIndex int) string {
	return PolicyPrefix(agentGroupName, policyName) + "-classifier_index-" + strconv.Itoa(classifierIndex)
}
