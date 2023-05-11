package status

import (
	"go.uber.org/fx"

	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
)

// Module is a fx module that provides a status Registry and registers status service handlers as grpc-gateway handlers.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideRegistry,
				fx.ParamTags(alerts.AlertsFxTag),
			),
		),
		grpcgateway.RegisterHandler{Handler: statusv1.RegisterStatusServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(RegisterStatusService),
	)
}

func provideRegistry(alerter alerts.Alerter, logger *log.Logger) Registry {
	return NewRegistry(logger, alerter)
}
