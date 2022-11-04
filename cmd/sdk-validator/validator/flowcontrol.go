package validator

import (
	"context"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/api/base"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

// FlowControlHandler implements FlowControlService.
type FlowControlHandler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer

	CommonHandler base.HandlerWithValues
	Rejects       int64
	Rejected      int64
}

// Check is a dummy Check handler.
func (f *FlowControlHandler) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	log.Trace().Msg("Received FlowControl Check request")

	services := []string{}
	rpcPeer, peerExists := peer.FromContext(ctx)
	if peerExists {
		services = append(services, rpcPeer.Addr.String())
	}

	start := time.Now()
	resp := f.CommonHandler.CheckWithValues(ctx, services, selectors.NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_FEATURE, req.Feature), req.Labels)
	// randomly reject requests based on rejectRatio
	if f.Rejected != f.Rejects {
		log.Trace().Msg("Rejecting call")
		resp.DecisionType = flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED
		resp.RejectReason = flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED
		atomic.AddInt64(&f.Rejected, 1)
	}
	end := time.Now()

	resp.Start = timestamppb.New(start)
	resp.End = timestamppb.New(end)

	return resp, nil
}
