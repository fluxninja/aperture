package checkhttp

import (
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/servicegetter"
)

// Module provides flowcontrol HTTP handler.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(ProvideHandler),
		grpcgateway.RegisterHandler{Handler: flowcontrolv1.RegisterFlowControlServiceHTTPHandlerFromEndpoint}.Annotate(),
		fx.Invoke(Register),
	)
}

// ConstructorIn holds parameters for ProvideHandler.
type ConstructorIn struct {
	fx.In

	ServiceGetter servicegetter.ServiceGetter
	Classifier    *classification.ClassificationEngine
	FCHandler     check.HandlerWithValues
}

// ProvideHandler provides a Flow Control Handler.
func ProvideHandler(
	in ConstructorIn,
) (flowcontrolv1.FlowControlServiceHTTPServer, error) {
	h := NewHandler(in.Classifier, in.ServiceGetter, in.FCHandler)

	return h, nil
}

// Register registers flowcontrol service on a gRPC server.
func Register(server *grpc.Server, handler flowcontrolv1.FlowControlServiceHTTPServer, healthsrv *health.Server) {
	flowcontrolv1.RegisterFlowControlServiceHTTPServer(server, handler)

	healthsrv.SetServingStatus("aperture.flowcontrol.v1.FlowControlServiceHTTP", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("flowcontrol http handler registered")
}
