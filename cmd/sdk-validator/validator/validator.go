package validator

import (
	"context"
	"time"

	"golang.org/x/exp/maps"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

type flowcontrolHandler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer
}

// Check is a dummy Check handler.
func (f *flowcontrolHandler) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	log.Info().Msg("Received Check request")

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

func (f *flowcontrolHandler) check(ctx context.Context, feature string, labels map[string]string, services []string) *flowcontrolv1.CheckResponse {
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

	return resp
}
