//go:generate swagger generate spec --scan-models --include="github.com/fluxninja*" --include-tag=common-configuration -o ../../docs/gen/config/agent/config-swagger.yaml

// Package main Agent
//
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
	"github.com/fluxninja/aperture/pkg/etcd/election"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/peers"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol"
	"github.com/fluxninja/aperture/pkg/policies/infra"
	"github.com/fluxninja/aperture/pkg/prometheus"
)

func main() {
	app := platform.New(
		platform.Config{}.Module(),
		prometheus.Module(),
		k8s.Module(),
		peers.Constructor{}.Module(),
		fx.Provide(
			agentinfo.ProvideAgentInfo,
			clockwork.NewRealClock,
			agent.ProvidePeersPrefix,
		),
		fx.Invoke(
			agent.AddAgentInfoAttribute,
		),
		entitycache.Module(),
		distcache.Module(),
		flowcontrol.Module(),
		infra.Module(),
		otelcollector.Module(),
		agent.ModuleForAgentOTEL(),
		discovery.Module(),
		election.Module(),
	)

	if err := app.Err(); err != nil {
		visualize, viserr := fx.VisualizeError(err)
		if viserr != nil {
			log.Panic().Err(viserr).Msgf("Failed to visualize fx error: %s", viserr)
		}
		log.Panic().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("aperture-agent app created")
	platform.Run(app)
}
