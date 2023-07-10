package validator

import (
	"context"
	"strings"
	"sync/atomic"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/check"
)

// CommonHandler implements common.HandlerWithValues.
type CommonHandler struct {
	check.HandlerWithValues

	Rejects  int64
	Rejected int64
}

const targetLabelMissing = "UNKNOWN"

// CheckRequest is a dummy function for creating *flowcontrolv1.CheckResponse from given parameters.
func (c *CommonHandler) CheckRequest(ctx context.Context,
	requestContext iface.RequestContext,
) *flowcontrolv1.CheckResponse {
	labels := requestContext.FlowLabels
	controlPoint := requestContext.ControlPoint
	services := requestContext.Services
	var path string
	var found bool
	if path, found = labels.Get("http.target"); !found {
		// traffic control points will have this label set
		log.Trace().Msg("Missing request path label")
		path = targetLabelMissing
	}
	log.Trace().Msgf("Received FlowControl Check request from path %v", path)

	resp := &flowcontrolv1.CheckResponse{
		DecisionType:  flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
		FlowLabelKeys: labels.SortedKeys(),
		Services:      services,
		ControlPoint:  controlPoint,
		RejectReason:  flowcontrolv1.CheckResponse_REJECT_REASON_NONE,
	}

	if c.Rejected != c.Rejects && shouldBeTested(path) {
		log.Trace().Msg("Rejecting call")
		resp.DecisionType = flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED
		resp.RejectReason = flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED
		atomic.AddInt64(&c.Rejected, 1)
	}

	return resp
}

func shouldBeTested(path string) bool {
	if path == targetLabelMissing {
		// handle feature control points
		return true
	}
	return strings.Contains(path, "super")
}
