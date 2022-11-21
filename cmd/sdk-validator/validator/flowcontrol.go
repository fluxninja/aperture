package validator

import (
	"context"
	"time"

	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
)

// FlowControlHandler implements FlowControlService.
type FlowControlHandler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer

	CommonHandler check.HandlerWithValues
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
	resp := f.CommonHandler.CheckWithValues(ctx, services, req.ControlPoint, req.Labels)
	end := time.Now()

	resp.Start = timestamppb.New(start)
	resp.End = timestamppb.New(end)

	return resp, nil
}
