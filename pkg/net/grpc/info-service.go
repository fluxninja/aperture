package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	infov1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/info/v1"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
)

// InfoService is the implementation of the infov1.InfoServiceServer interface.
type InfoService struct {
	infov1.UnimplementedInfoServiceServer
}

// Version returns the version of the service.
func (vh *InfoService) Version(ctx context.Context, req *emptypb.Empty) (*infov1.VersionInfo, error) {
	log.Trace().Msg("Received request on version info handler")
	resp := info.GetVersionInfo()
	return resp, nil
}

// Process returns the process info of the service.
func (vh *InfoService) Process(ctx context.Context, req *emptypb.Empty) (*infov1.ProcessInfo, error) {
	log.Trace().Msg("Received request on process info handler")
	resp := info.GetProcessInfo()
	return resp, nil
}

// Host returns the hostname of the service.
func (vh *InfoService) Host(ctx context.Context, req *emptypb.Empty) (*infov1.HostInfo, error) {
	log.Trace().Msg("Received request on host info handler")
	resp := info.GetHostInfo()
	return resp, nil
}

// RegisterInfoService registers the InfoService implementation with the provided grpc server.
func RegisterInfoService(server *grpc.Server) {
	vh := &InfoService{}
	infov1.RegisterInfoServiceServer(server, vh)
}
