package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/status"
)

// StatusService is the implementation of the statusv1.StatusServiceServer interface.
type StatusService struct {
	statusv1.UnimplementedStatusServiceServer
	registry *status.Registry
}

// RegisterStatusService registers the StatusService implementation with the provided grpc server.
func RegisterStatusService(server *grpc.Server, reg *status.Registry) {
	svc := &StatusService{
		registry: reg,
	}
	statusv1.RegisterStatusServiceServer(server, svc)
}

// GetGroupStatus returns the group status for the requested group in the Registry.
func (svc *StatusService) GetGroupStatus(ctx context.Context, req *statusv1.GroupStatusRequest) (*statusv1.GroupStatus, error) {
	log.Trace().Str("group", req.Group).Msg("Received request on GetGroupStatus handler")

	status := svc.registry.At(req.Group).Get()

	return status, nil
}

// GetGroups returns the groups from the keys in the Registry.
func (svc *StatusService) GetGroups(ctx context.Context, req *emptypb.Empty) (*statusv1.Groups, error) {
	log.Trace().Msg("Received request on GetGroups handler")

	response := &statusv1.Groups{
		Groups: svc.registry.Keys(),
	}

	return response, nil
}
