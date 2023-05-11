package status

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
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
	log.Trace().Str("path", req.Path).Msg("Received request on GetGroupStatus handler")

	// extract pairs of key-value from the path, separated by /
	pathKeys := strings.Split(req.Path, "/")

	registry := svc.registry
	if req.Path == "" || len(pathKeys) == 0 {
		return registry.GetGroupStatus(), nil
	}

	for i := 0; i < len(pathKeys); i += 2 {
		if i >= len(pathKeys) || i+1 >= len(pathKeys) {
			log.Warn().Str("path", req.Path).Msg("Incorrect status path")
			_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
			return &statusv1.GroupStatus{}, nil
		}

		key := pathKeys[i]
		val := pathKeys[i+1]
		registry = registry.ChildIfExists(key, val)
		if registry == nil {
			_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
			return &statusv1.GroupStatus{}, nil
		}
		if registry.HasError() {
			_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "503"))
		}
	}

	return registry.GetGroupStatus(), nil
}
