//go:generate swagger generate spec --scan-models --include="github.com/FluxNinja*" --include-tag=common-configuration -o ../../docs/gen/config/aperture-agent/config-swagger.yaml

// Aperture Agent
//   BasePath: /aperture-agent
// swagger:meta
package main

import (
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"

	"github.com/FluxNinja/aperture/cmd/aperture-agent/agent"
	"github.com/FluxNinja/aperture/pkg/agentinfo"
	"github.com/FluxNinja/aperture/pkg/authz"
	"github.com/FluxNinja/aperture/pkg/classification"
	"github.com/FluxNinja/aperture/pkg/discovery"
	"github.com/FluxNinja/aperture/pkg/distcache"
	"github.com/FluxNinja/aperture/pkg/entitycache"
	"github.com/FluxNinja/aperture/pkg/flowcontrol"
	"github.com/FluxNinja/aperture/pkg/k8s"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/net/grpc"
	"github.com/FluxNinja/aperture/pkg/net/http"
	"github.com/FluxNinja/aperture/pkg/notifiers"
	"github.com/FluxNinja/aperture/pkg/otel"
	"github.com/FluxNinja/aperture/pkg/otelcollector"
	"github.com/FluxNinja/aperture/pkg/platform"
	"github.com/FluxNinja/aperture/pkg/policies/dataplane"
	"github.com/FluxNinja/aperture/pkg/prometheus"
)

func main() {
	app := platform.New(
		platform.Config{}.Module(),
		http.ClientConstructor{Name: "k8s-http-client", Key: "kubernetes.http_client"}.Annotate(),
		notifiers.TrackersConstructor{Name: "entity_trackers"}.Annotate(),
		prometheus.Module(),
		otel.ProvideAnnotatedAgentConfig(),
		fx.Provide(
			agentinfo.ProvideAgentInfo,
			k8s.Providek8sClient,
			clockwork.NewRealClock,
			entitycache.ProvideEntityCache,
			otel.AgentOTELComponents,
			agent.ProvidePeersPrefix,
		),
		flowcontrol.Module,
		classification.Module,
		authz.Module,
		otelcollector.Module(),
		distcache.Module(),
		dataplane.PolicyModule(),
		discovery.Module(),
		fx.Invoke(
			authz.Register,
			flowcontrol.Register,
		),
		grpc.ClientConstructor{Name: "flowcontrol-grpc-client", Key: "flowcontrol.client.grpc"}.Annotate(),
	)

	if err := app.Err(); err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Fatal().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("aperture-agent app created")
	platform.Run(app)
}
