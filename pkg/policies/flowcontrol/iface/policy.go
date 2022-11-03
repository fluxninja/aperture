package iface

// Policy is the interface that wraps the GetPolicyName, GetPolicyHash methods.
type Policy interface {
	GetPolicyName() string
	GetPolicyHash() string
}
