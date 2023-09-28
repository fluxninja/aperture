package iface

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/status"
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
	TicksInDurationPb(*durationpb.Duration) int
	TicksInDuration(time.Duration) int
}

// GetSelectorsShortDescription returns a short description of the selectors.
func GetSelectorsShortDescription(selectors []*policylangv1.Selector) string {
	return fmt.Sprintf("%d selectors", len(selectors))
}
