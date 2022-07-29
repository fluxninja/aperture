package component

import (
	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
)

// ComponentAPI is the interface that wraps the GetPolicyName, GetPolicyID, and GetComponentIndex methods.
type ComponentAPI interface {
	GetAgentGroupName() string
	GetPolicyName() string
	GetPolicyHash() string
	GetComponentIndex() int64
}

// Component is a struct that provides ComponentAPI based on wrapper proto message that contains a policy name, hash, and component ID.
type Component struct {
	ComponentAPI
}

// Make sure Component implements ComponentAPI.
var _ ComponentAPI = &Component{}

// NewComponent returns a new Component.
func NewComponent(wrapperMessage *configv1.ConfigPropertiesWrapper) *Component {
	if wrapperMessage == nil {
		return nil
	}

	return &Component{
		ComponentAPI: wrapperMessage,
	}
}
