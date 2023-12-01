package controlpoints

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	flowcontrolpointsv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
)

// Module provides preview handler and registers the service.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		grpcgateway.RegisterHandler{Handler: flowcontrolpointsv1.RegisterFlowControlPointsServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(Register),
	)
}

// RegisterIn bundles and annotates parameters.
type RegisterIn struct {
	fx.In
	Server  *grpc.Server `name:"default"`
	Handler *Handler
}

// Register registers the handler on grpc.Server.
func Register(in RegisterIn) error {
	flowcontrolpointsv1.RegisterFlowControlPointsServiceServer(in.Server, in.Handler)

	log.Info().Msg("FlowControl ControlPoints handler registered")
	return nil
}
