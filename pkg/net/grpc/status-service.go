package grpc

import (
	"google.golang.org/grpc"

	statusv1 "aperture.tech/aperture/api/gen/proto/go/aperture/common/status/v1"
	"aperture.tech/aperture/pkg/status"
)

// StatusService is the implementation of the statusv1.StatusServiceServer interface.
type StatusService struct {
	statusv1.UnimplementedStatusServiceServer
}

// RegisterStatusService registers the StatusService implementation with the provided grpc server.
func RegisterStatusService(server *grpc.Server, reg *status.Registry) {
	statusv1.RegisterStatusServiceServer(server, reg)
}
