package iface

import "strconv"

// Component is the interface that wraps the GetPolicyName, GetPolicyHash, and GetComponentIndex methods.
type Component interface {
	Policy
	GetComponentIndex() int64
}

// ComponentID returns the metric ID for a component.
func ComponentID(component Component) string {
	return ComponentIDExpanded(component.GetPolicyName(), component.GetComponentIndex(), component.GetPolicyHash())
}

// ComponentIDExpanded returns the metric ID for a component.
func ComponentIDExpanded(policyName string, componentIndex int64, policyHash string) string {
	return "policy-" + policyName + "-component_index-" + strconv.FormatInt(componentIndex, 10) + "-policy_hash-" + policyHash
}
