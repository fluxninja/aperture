package status

import (
	"go.uber.org/fx"

	statusv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/common/status/v1"
	"github.com/FluxNinja/aperture/pkg/net/grpcgateway"
)

// Module is a fx module that provides a status registry and registers status service handlers as grpcgateway handlers.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideRegistry),
		grpcgateway.RegisterHandler{Handler: statusv1.RegisterStatusServiceHandlerFromEndpoint}.Annotate(),
	)
}

func provideRegistry() *Registry {
	return NewRegistry(".")
}
