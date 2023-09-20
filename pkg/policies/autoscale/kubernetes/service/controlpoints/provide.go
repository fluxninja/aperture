package controlpoints

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	controlpointsv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
)

// Module returns an fx.Option that provides the Kubernetes discovery module.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewHandler),
		fx.Invoke(
			RegisterControlPointCacheService,
			RegisterControlPointsHandler,
		),
		grpcgateway.RegisterHandler{Handler: controlpointsv1.RegisterAutoScaleKubernetesControlPointsServiceHandlerFromEndpoint}.Annotate(),
	)
}

// RegisterControlPointCacheServiceIn bundles and annotates parameters.
type RegisterControlPointCacheServiceIn struct {
	fx.In
	Server  *grpc.Server `name:"default"`
	Handler *Handler
}

// RegisterControlPointCacheService registers the ControlPointCache service handler with the gRPC server.
func RegisterControlPointCacheService(in RegisterControlPointCacheServiceIn) {
	controlpointsv1.RegisterAutoScaleKubernetesControlPointsServiceServer(in.Server, in.Handler)
}
