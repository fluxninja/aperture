package controlpoints

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	flowcontrolpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
)

// Module provides preview handler and registers the service.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		grpcgateway.RegisterHandler{Handler: flowcontrolpointsv1.RegisterFlowControlPointsServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(Register),
	)
}

// Register registers the handler on grpc.Server.
func Register(handler *Handler, server *grpc.Server) error {
	flowcontrolpointsv1.RegisterFlowControlPointsServiceServer(server, handler)

	log.Info().Msg("FlowControl ControlPoints handler registered")
	return nil
}
