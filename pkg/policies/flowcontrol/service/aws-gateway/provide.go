package awsgateway

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
)

// Module provides awsgateway flowcontrol handler.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		grpcgateway.RegisterHandler{Handler: flowcontrolv1.RegisterAWSGatewayFlowControlServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(Register),
	)
}

// Register registers the handler on grpc.Server
//
// To be used in fx.Invoke.
func Register(server *grpc.Server, handler *Handler, healthsrv *health.Server) {
	// If changing params to this function, keep RegisterAnnotated in sync.
	flowcontrolv1.RegisterAWSGatewayFlowControlServiceServer(server, handler)

	healthsrv.SetServingStatus("aperture.flowcontrol.v1.AWSGatewayFlowControlService", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("Flowcontrol handler for AWS Gateway registered")
}
