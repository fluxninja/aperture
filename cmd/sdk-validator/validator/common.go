package validator

import (
	"context"

	"golang.org/x/exp/maps"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/flowcontrol/common"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/selectors"
)

// CommonHandler implements common.HandlerWithValues.
type CommonHandler struct {
	common.HandlerWithValues
}

// CheckWithValues is a dummy function for creating *flowcontrolv1.CheckResponse from given parameters.
func (c *CommonHandler) CheckWithValues(ctx context.Context, services []string, controlPoint selectors.ControlPoint, labels map[string]string) *flowcontrolv1.CheckResponse {
	resp := &flowcontrolv1.CheckResponse{
		DecisionType:     flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
		FlowLabelKeys:    maps.Keys(labels),
		Services:         services,
		ControlPointInfo: controlPoint.ToControlPointInfoProto(),
		RejectReason:     flowcontrolv1.CheckResponse_REJECT_REASON_NONE,
	}

	return resp
}
