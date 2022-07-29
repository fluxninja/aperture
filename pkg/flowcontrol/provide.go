package flowcontrol

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	flowcontrolv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/FluxNinja/aperture/pkg/entitycache"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/policies/dataplane/iface"
	"github.com/FluxNinja/aperture/pkg/webhooks/validation"
)

// Module is a set of default providers for flowcontrol components
//
// Note that the handler needs to be Registered for flowcontrol to be available
// externally.
var Module = fx.Provide(
	ProvideMetrics,
	ProvideHandler,
)

// ConstructorIn holds parameters for ProvideHandler.
type ConstructorIn struct {
	fx.In

	EntityCache *entitycache.EntityCache
	Metrics     Metrics
	EngineAPI   iface.EngineAPI
}

// ProvideHandler provides a Flow Control Handler.
func ProvideHandler(in ConstructorIn) (flowcontrolv1.FlowControlServiceServer, HandlerWithValues, error) {
	h := NewHandler(in.EntityCache, in.Metrics, in.EngineAPI)

	// Note: Returning the same handler twice as different interfaces – once as
	// a handler to be registered on grpc server and once for consumption by
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

// Register registers flowcontrol service on a grpc server.
func Register(server *grpc.Server, handler flowcontrolv1.FlowControlServiceServer, healthsrv *health.Server) {
	flowcontrolv1.RegisterFlowControlServiceServer(server, handler)

	healthsrv.SetServingStatus("aperture.flowcontrol.v1.FlowControlService", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("flowcontrol handler registered")
}

// ProvideCMFileValidator provides flowcontrol config map file validator
//
// Note: This validator must be registered to be accessible.
func ProvideCMFileValidator() *CMFileValidator {
	return &CMFileValidator{}
}

// RegisterCMFileValidator registers flowcontrol configmap validator as configmap file validator.
func RegisterCMFileValidator(validator *CMFileValidator, configMapValidator *validation.CMValidator) {
	configMapValidator.RegisterCMFileValidator(validator)
}
