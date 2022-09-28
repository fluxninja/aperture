package status

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

// StatusService is the implementation of the statusv1.StatusServiceServer interface.
type StatusService struct {
	statusv1.UnimplementedStatusServiceServer
	registry Registry
}

// RegisterStatusService registers the StatusService implementation with the provided grpc server.
func RegisterStatusService(server *grpc.Server, reg Registry) {
	svc := &StatusService{
		registry: reg,
	}
	statusv1.RegisterStatusServiceServer(server, svc)
}

// GetGroupStatus returns the group status for the requested group in the Registry.
func (svc *StatusService) GetGroupStatus(ctx context.Context, req *statusv1.GroupStatusRequest) (*statusv1.GroupStatus, error) {
	log.Trace().Interface("path", req.Path).Msg("Received request on GetGroupStatus handler")

	// extract keys from the path, separated by /
	keys := strings.Split(req.Path, "/")

	registry := svc.registry
	for _, key := range keys {
		if key == "" {
			continue
		}
		registry = registry.ChildIfExists(key)
		if registry == nil {
			return &statusv1.GroupStatus{}, nil
		}
	}
	return registry.GetGroupStatus(), nil
}
