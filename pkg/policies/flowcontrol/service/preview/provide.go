// +kubebuilder:validation:Optional
package preview

import (
	flowpreviewv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	cfg "github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Module provides preview handler and registers the service.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		grpcgateway.RegisterHandler{Handler: flowpreviewv1.RegisterFlowPreviewServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(Register),
	)
}

// FlowPreviewConfig is the configuration for the flow control preview service.
// swagger:model
// +kubebuilder:object:generate=true
type FlowPreviewConfig struct {
	// Enables the flow preview service.
	Enabled bool `json:"enabled" default:"true"`
}

// Register registers the handler on grpc.Server.
func Register(handler *Handler,
	server *grpc.Server,
	healthsrv *health.Server,
	unmarshaller cfg.Unmarshaller,
) error {
	var config FlowPreviewConfig
	if err := unmarshaller.UnmarshalKey("flow_control.preview_service", &config); err != nil {
		return err
	}

	if !config.Enabled {
		log.Info().Msg("flow preview service disabled")
		return nil
	}
	flowpreviewv1.RegisterFlowPreviewServiceServer(server, handler)

	healthsrv.SetServingStatus("aperture.flowcontrol.v1.FlowPreviewService", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("Preview handler registered")
	return nil
}
