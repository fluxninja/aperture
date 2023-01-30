package iface

// Component is the interface that wraps the GetPolicyName, GetPolicyHash, and GetComponentID methods.
type Component interface {
	Policy
	GetComponentId() string
}

// ComponentKey returns a unique Key for a component.
func ComponentKey(component Component) string {
	return "policy-" + component.GetPolicyName() + "-component_id-" + component.GetComponentId() + "-policy_hash-" + component.GetPolicyHash()
}
