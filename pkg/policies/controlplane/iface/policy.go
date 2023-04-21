package iface

import (
	"fmt"
	"time"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/status"
)

const (
	// PoliciesRoot - path in config and status registry for policies results.
	PoliciesRoot = "policies"
)

// FxOptionsFuncTag allows sub-modules to provide their options to per policy apps independently.
var FxOptionsFuncTag = config.GroupTag("policy-fx-funcs")

// PolicyBase is for read only access to base policy info.
type PolicyBase interface {
	GetPolicyName() string
	GetPolicyHash() string
}

// Policy is for read only access to full policy state.
type Policy interface {
	PolicyBase
	GetEvaluationInterval() time.Duration
	GetStatusRegistry() status.Registry
}

// GetServiceShortDescription returns a short description of the service selector.
func GetServiceShortDescription(flowSelector *policylangv1.FlowSelector) string {
	agentGroup := flowSelector.GetAgentGroup()
	var prefix string
	if agentGroup == "default" {
		prefix = ""
	} else {
		prefix = fmt.Sprintf("%s/", agentGroup)
	}
	return prefix
}
