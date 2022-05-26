//go:generate swagger generate spec --scan-models --include="aperture.tech*" --include-tag=common-configuration -o ../../docs/gen/config/aperture-agent/config-swagger.yaml

// Aperture Agent
//   BasePath: /aperture-agent
// swagger:meta
package main

import (
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"

	"aperture.tech/aperture/cmd/aperture-agent/agent"
	"aperture.tech/aperture/pkg/agentinfo"
	"aperture.tech/aperture/pkg/authz"
	"aperture.tech/aperture/pkg/classification"
	"aperture.tech/aperture/pkg/discovery"
	"aperture.tech/aperture/pkg/distcache"
	"aperture.tech/aperture/pkg/entitycache"
	"aperture.tech/aperture/pkg/flowcontrol"
	"aperture.tech/aperture/pkg/k8s"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/net/grpc"
	"aperture.tech/aperture/pkg/net/http"
	"aperture.tech/aperture/pkg/notifiers"
	"aperture.tech/aperture/pkg/otel"
	"aperture.tech/aperture/pkg/otelcollector"
	"aperture.tech/aperture/pkg/platform"
	"aperture.tech/aperture/pkg/policies/dataplane"
	"aperture.tech/aperture/pkg/prometheus"
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
