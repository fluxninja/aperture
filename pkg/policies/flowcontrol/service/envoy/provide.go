package envoy

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/fluxninja/aperture/pkg/log"
)

// Module provides authz handler
//
// Authz handler is one of the APIs to classification and flowcontrol modules.
// Authz uses envoy's external authorization grpc API.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		fx.Invoke(Register),
	)
}

// Register registers the handler on grpc.Server
//
// To be used in fx.Invoke.
func Register(handler *Handler, server *grpc.Server, healthsrv *health.Server) {
	// If changing params to this function, keep RegisterAnnotated in sync.
	authv3.RegisterAuthorizationServer(server, handler)

	healthsrv.SetServingStatus("envoy.service.auth.v3.Authorization", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("Authz handler registered")
}
