package check

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	servicegetter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service-getter"
)

// Module is a set of default providers for flowcontrol components
//
// Note that the handler needs to be Registered for flowcontrol to be available
// externally.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			ProvideMetrics,
			ProvideHandler,
		),
		fx.Invoke(Register),
	)
}

// ConstructorIn holds parameters for ProvideHandler.
type ConstructorIn struct {
	fx.In

	ServiceGetter servicegetter.ServiceGetter
	Metrics       Metrics
	EngineAPI     iface.Engine
	Cache         iface.Cache `optional:"true"`
}

// ProvideHandler provides a Flow Control Handler.
func ProvideHandler(
	in ConstructorIn,
) (flowcontrolv1.FlowControlServiceServer, HandlerWithValues, error) {
	h := NewHandler(in.ServiceGetter, in.Metrics, in.EngineAPI, in.Cache)

	// Note: Returning the same handler twice as different interfaces – once as
	// a handler to be registered on gRPC server and once for consumption by
	// authz
	return h, h, nil
}

// ProvideDummyHandler provides an empty Flow Control Handler.
var ProvideDummyHandler = fx.Annotate(NewHandler, fx.As(new(HandlerWithValues)))

// ProvideMetrics provides flowcontrol metrics that hook to prometheus registry.
func ProvideMetrics(promRegistry *prometheus.Registry) (Metrics, error) {
	metrics, err := NewPrometheusMetrics(promRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed creating Prometheus collector: %v", err)
	}
	return metrics, nil
}

// ProvideNopMetrics provides disabled flowcontrol metrics.
func ProvideNopMetrics() Metrics { return NopMetrics{} }

// RegisterIn bundles and annotates parameters.
type RegisterIn struct {
	fx.In
	Server       *grpc.Server `name:"default"`
	Handler      flowcontrolv1.FlowControlServiceServer
	HealthServer *health.Server
}

// Register registers flowcontrol service on a gRPC server.
func Register(in RegisterIn) {
	flowcontrolv1.RegisterFlowControlServiceServer(in.Server, in.Handler)

	in.HealthServer.SetServingStatus("aperture.flowcontrol.v1.FlowControlService", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("flowcontrol handler registered")
}
