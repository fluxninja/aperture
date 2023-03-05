package controlpoints

import (
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	controlpointcachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
)

// Module returns an fx.Option that provides the Kubernetes discovery module.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		fx.Invoke(RegisterControlPointCacheService),
		grpcgateway.RegisterHandler{Handler: controlpointcachev1.RegisterControlPointsServiceHandlerFromEndpoint}.Annotate(),
	)
}

// RegisterControlPointCacheService registers the ControlPointCache service handler with the gRPC server.
func RegisterControlPointCacheService(handler *Handler, server *grpc.Server) {
	controlpointcachev1.RegisterControlPointsServiceServer(server, handler)
}
