//go:generate swagger generate spec --scan-models --include="github.com/fluxninja*" --include-tag=common-configuration -o ../../docs/gen/config/aperture-agent/config-swagger.yaml

// Aperture Agent
//
//	BasePath: /aperture-agent
//
// swagger:meta
package main

import (
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/cmd/aperture-agent/agent"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/discovery"
	"github.com/fluxninja/aperture/pkg/distcache"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/flowcontrol"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/otel"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/peers"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/policies/dataplane"
	"github.com/fluxninja/aperture/pkg/prometheus"
)

func main() {
	app := platform.New(
		platform.Config{}.Module(),
		notifiers.TrackersConstructor{Name: "entity_trackers"}.Annotate(),
		prometheus.Module(),
		k8s.Module(),
		otel.OTELConfigConstructor{Type: otel.AgentType}.Annotate(),
		peers.Constructor{}.Module(),
		fx.Provide(
			agentinfo.ProvideAgentInfo,
			clockwork.NewRealClock,
			otel.AgentOTELComponents,
			agent.ProvidePeersPrefix,
		),
		fx.Invoke(
			agent.AddAgentInfoAttribute,
		),
		entitycache.Module(),
		flowcontrol.Module,
		distcache.Module(),
		dataplane.PolicyModule(),
		discovery.Module(),
		grpc.ClientConstructor{Name: "flowcontrol-grpc-client", ConfigKey: "flowcontrol.client.grpc"}.Annotate(),
		otelcollector.Module(),
	)

	if err := app.Err(); err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Panic().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("aperture-agent app created")
	platform.Run(app)
}
