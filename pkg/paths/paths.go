package paths

import (
	"path"
	"strconv"

	"github.com/fluxninja/aperture/pkg/policies/dataplane/component"
)

var (
	// ConfigPrefix is key prefix in etcd for config.
	ConfigPrefix = path.Join("/config")
	// DecisionsPrefix is key prefix in etcd for decisions.
	DecisionsPrefix = path.Join("/decisions")
	// Classifiers is config path in etcd for classifiers.
	Classifiers = path.Join(ConfigPrefix, "classifiers")
	// Policies is config path in etcd for policies.
	Policies = path.Join(ConfigPrefix, "policies")
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

// IdentifierForComponent returns the identifier for a component.
func IdentifierForComponent(agentGroupName, policyName string, componentIndex int64) string {
	return PolicyPrefix(agentGroupName, policyName) + "-component_index-" + strconv.FormatInt(componentIndex, 10)
}

// MetricIDForComponent returns the metric ID for a component.
func MetricIDForComponent(componentAPI component.ComponentAPI) string {
	return MetricIDForComponentExpanded(componentAPI.GetAgentGroup(), componentAPI.GetPolicyName(), componentAPI.GetComponentIndex(), componentAPI.GetPolicyHash())
}

// MetricIDForComponentExpanded returns the metric ID for a component.
func MetricIDForComponentExpanded(agentGroupName, policyName string, componentIndex int64, policyHash string) string {
	return IdentifierForComponent(agentGroupName, policyName, componentIndex) + "-policy_hash-" + policyHash
}

// IdentifierForFluxMeter returns the identifier for FluxMeter.
func IdentifierForFluxMeter(agentGroupName, policyName, fluxMeterName string) string {
	return PolicyPrefix(agentGroupName, policyName) + "-flux_meter-" + fluxMeterName
}

// MetricIDForFluxMeter returns the metric ID for a flux meter.
func MetricIDForFluxMeter(componentAPI component.ComponentAPI, fluxMeterName string) string {
	return MetricIDForFluxMeterExpanded(componentAPI.GetAgentGroup(), componentAPI.GetPolicyName(), fluxMeterName, componentAPI.GetPolicyHash())
}

// MetricIDForFluxMeterExpanded returns the metric ID for a flux meter.
func MetricIDForFluxMeterExpanded(agentGroupName, policyName, fluxMeterName, policyHash string) string {
	return IdentifierForFluxMeter(agentGroupName, policyName, fluxMeterName) + "-policy_hash-" + policyHash
}
