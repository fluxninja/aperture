package validator

import (
	"context"
	"sync/atomic"
	"time"

	"golang.org/x/exp/maps"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

// FlowControlHandler implements FlowControlService.
type FlowControlHandler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer

	Rejects  int64
	Rejected int64
}

// Check is a dummy Check handler.
func (f *FlowControlHandler) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	log.Trace().Msg("Received Check request")

	services := []string{}
	rpcPeer, peerExists := peer.FromContext(ctx)
	if peerExists {
		services = append(services, rpcPeer.Addr.String())
	}

	start := time.Now()
	resp := f.check(ctx, req.Feature, req.Labels, services)
	end := time.Now()

	resp.Start = timestamppb.New(start)
	resp.End = timestamppb.New(end)

	return resp, nil
}

func (f *FlowControlHandler) check(ctx context.Context, feature string, labels map[string]string, services []string) *flowcontrolv1.CheckResponse {
	resp := &flowcontrolv1.CheckResponse{
		DecisionType:  flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
		FlowLabelKeys: maps.Keys(labels),
		Services:      services,
		ControlPointInfo: &flowcontrolv1.ControlPointInfo{
			Feature: feature,
			Type:    flowcontrolv1.ControlPointInfo_TYPE_FEATURE,
		},
		RejectReason: flowcontrolv1.CheckResponse_REJECT_REASON_NONE,
	}

	// randomly reject requests based on rejectRatio
	// nolint:gosec
	if f.Rejected != f.Rejects {
		log.Trace().Msg("Rejecting call")
		resp.DecisionType = flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED
		resp.RejectReason = flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED
		atomic.AddInt64(&f.Rejected, 1)
	}

	return resp
}
