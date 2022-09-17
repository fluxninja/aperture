package envoy

import (
	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

// logSampled provides log sampling for envoy handler.
var logSampled log.Logger = log.Sample(&zerolog.BasicSampler{N: 1000})

// Module provides authz handler
//
// Authz handler is one of the APIs to classification and flowcontrol modules.
// Authz uses envoy's external authorization grpc API.
//
// Note: Register function is not bundled inside this module and should be
// invoked explicitly.
var Module = fx.Options(
	fx.Provide(ProvideHandler),
	fx.Invoke(Register),
)

// ProvideHandler provides an authz handler
//
// See NewHandler for more docs.
var ProvideHandler = NewHandler

// Register registers the handler on grpc.Server
//
// To be used in fx.Invoke.
func Register(handler *Handler, server *grpc.Server, healthsrv *health.Server) {
	// If changing params to this function, keep RegisterAnnotated in sync.
	ext_authz.RegisterAuthorizationServer(server, handler)

	healthsrv.SetServingStatus("envoy.service.auth.v3.Authorization", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("Authz handler registered")
}

// OnNamedServer returns a register function that will register authz handler on
// *named* grpc.Server
//
// Usage:
//
//	fx.Invoke(authz.OnNamedServer("foo").Register)
func OnNamedServer(serverName string) Invocations {
	return Invocations{
		Register: fx.Annotate(
			Register,
			fx.ParamTags(``, config.NameTag(serverName)),
		),
	}
}

// Invocations is a set of register functions to be used in fx.Invoke.
type Invocations struct {
	Register interface{}
}
