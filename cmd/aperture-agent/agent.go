//go:generate swagger generate spec --scan-models --include="github.com/fluxninja*" --include-tag=common-configuration --include-tag=agent-configuration -o ../../docs/gen/config/agent/config-swagger.yaml
//go:generate go run ../../docs/tools/swagger/process-go-tags.go ../../docs/gen/config/agent/config-swagger.yaml

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

	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent"
	agentfunctions "github.com/fluxninja/aperture/v2/pkg/agent-functions"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/discovery"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	"github.com/fluxninja/aperture/v2/pkg/etcd/election"
	"github.com/fluxninja/aperture/v2/pkg/k8s"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector"
	"github.com/fluxninja/aperture/v2/pkg/peers"
	"github.com/fluxninja/aperture/v2/pkg/platform"
	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol"
	"github.com/fluxninja/aperture/v2/pkg/prometheus"
	"github.com/fluxninja/aperture/v2/pkg/rpc"
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
		distcache.Module(),
		flowcontrol.Module(),
		autoscale.Module(),
		agent.ModuleForAgentOTel(),
		discovery.Module(),
		election.Module(),
		rpc.ClientModule,
		agentfunctions.Module,
		Module(),
		// Start collector after all extensions started, so it won't
		// immediately reload when extensions add their config.
		otelcollector.Module(),
	)

	defer log.WaitFlush()

	if err := app.Err(); err != nil {
		v, verr := fx.VisualizeError(err)
		if verr != nil {
			log.Error().Err(verr).Msg("Failed to visualize fx error")
			return
		}
		log.Error().Err(err).Str("visualize", v).Msg("Failed to run create an initialize platform")
		return
	}

	log.Info().Msg("aperture-agent app created")
	platform.Run(app)
}
