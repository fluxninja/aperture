package status

import (
	"go.uber.org/fx"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
)

// Module is a fx module that provides a status Registry and registers status service handlers as grpc-gateway handlers.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideRegistry),
		grpcgateway.RegisterHandler{Handler: statusv1.RegisterStatusServiceHandlerFromEndpoint}.Annotate(),
	)
}

func provideRegistry() Registry {
	return NewRegistry()
}
