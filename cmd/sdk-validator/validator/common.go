package validator

import (
	"context"
	"sync/atomic"

	"golang.org/x/exp/maps"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
)

// CommonHandler implements common.HandlerWithValues.
type CommonHandler struct {
	check.HandlerWithValues

	Rejects  int64
	Rejected int64
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

	if c.Rejected != c.Rejects {
		log.Trace().Msg("Rejecting call")
		resp.DecisionType = flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED
		resp.RejectReason = flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED
		atomic.AddInt64(&c.Rejected, 1)
	}

	return resp
}
