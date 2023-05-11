package cmd

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Module is a module for running cmd.v1.Controller service.
var Module = fx.Options(
	fx.Provide(NewHandler),
	fx.Invoke(RegisterControllerServer),
)

// RegisterControllerServer registers handler for cmd.v1.Controller service.
func RegisterControllerServer(handler *Handler, server *grpc.Server, healthsrv *health.Server) {
	cmdv1.RegisterControllerServer(server, handler)

	healthsrv.SetServingStatus(
		"aperture.cmd.v1.Controller",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)
	log.Info().Msg("Controller handler registered")
}
