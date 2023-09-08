package preview

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	flowpreviewv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	cfg "github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
	previewconfig "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/preview/config"
)

// Module provides preview handler and registers the service.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		grpcgateway.RegisterHandler{Handler: flowpreviewv1.RegisterFlowPreviewServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(Register),
	)
}

// RegisterIn bundles and annotates parameters.
type RegisterIn struct {
	fx.In
	Server       *grpc.Server `name:"default"`
	Handler      *Handler
	HealthServer *health.Server
	Unmarshaller cfg.Unmarshaller
}

// Register registers the handler on grpc.Server.
func Register(in RegisterIn) error {
	var config previewconfig.FlowPreviewConfig
	if err := in.Unmarshaller.UnmarshalKey(previewconfig.Key, &config); err != nil {
		return err
	}

	if !config.Enabled {
		log.Info().Msg("flow preview service disabled")
		return nil
	}
	flowpreviewv1.RegisterFlowPreviewServiceServer(in.Server, in.Handler)

	in.HealthServer.SetServingStatus("aperture.flowcontrol.v1.FlowPreviewService", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("Preview handler registered")
	return nil
}
