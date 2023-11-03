package envoy

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Module provides authz handler
//
// Authz handler is one of the APIs to classification and flowcontrol modules.
// Authz uses envoy's external authorization gRPC API.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		fx.Invoke(Register),
	)
}

// RegisterIn bundles and annotates parameters.
type RegisterIn struct {
	fx.In
	Server       *grpc.Server `name:"default"`
	Handler      *Handler
	HealthServer *health.Server
}

// Register registers the handler on grpc.Server
//
// To be used in fx.Invoke.
func Register(in RegisterIn) {
	// If changing params to this function, keep RegisterAnnotated in sync.
	authv3.RegisterAuthorizationServer(in.Server, in.Handler)

	in.HealthServer.SetServingStatus("envoy.service.auth.v3.Authorization", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("Authz handler registered")
}
