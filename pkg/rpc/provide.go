package rpc

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	rpcv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/rpc/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

// ServerModule are components needed for server-side of rpc.
var ServerModule = fx.Options(
	fx.Provide(NewClients),
	fx.Provide(NewStreamServer),
	fx.Invoke(RegisterStreamServer),
)

// ClientModule are components needed for client-side of rpc
//
// Note: Not providing StreamClient, as this package is generic and doesn't
// know what to connect to.
var ClientModule = fx.Provide(NewHandlerRegistry)

// RegisterStreamServer registers the handler on grpc.Server
//
// To be used in fx.Invoke.
func RegisterStreamServer(handler *StreamServer, server *grpc.Server, healthsrv *health.Server) {
	rpcv1.RegisterCoordinatorServer(server, handler)

	healthsrv.SetServingStatus(
		"aperture.rpc.v1.Coordinator",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)
	log.Info().Msg("Coordinator handler registered")
}
