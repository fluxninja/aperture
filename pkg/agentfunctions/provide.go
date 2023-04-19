package agentfunctions

import (
	"go.uber.org/fx"

	afconfig "github.com/fluxninja/aperture/pkg/agentfunctions/config"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/rpc"
)

// Module provides rpc client for agent functions.
var Module = fx.Options(
	// FIXME(krdln) Do we actually need a separate grpc client for each module?
	grpcclient.ClientConstructor{
		Name:      "agent-functions",
		ConfigKey: afconfig.Key + ".client.grpc",
	}.Annotate(),
	fx.Provide(
		NewFlowControlControlPointsHandler,
		ProvidePreviewHandler,
	),
	fx.Invoke(
		RegisterClient,
		RegisterControlPointsHandler,
		RegisterPreviewHandler,
	),
)

// RegisterClientIn are parameters for InvokeClient function.
type RegisterClientIn struct {
	fx.In
	Lc           fx.Lifecycle
	Unmarshaller config.Unmarshaller
	Handlers     *rpc.HandlerRegistry
	ConnBuilder  grpcclient.ClientConnectionBuilder `name:"agent-functions"`
}

// RegisterClient registers a client which will allow calling agent functions from controller.
func RegisterClient(in RegisterClientIn) error {
	var config afconfig.AgentFunctionsConfig
	if err := in.Unmarshaller.UnmarshalKey(afconfig.Key, &config); err != nil {
		return err
	}

	for _, addr := range config.Endpoints {
		rpc.RegisterStreamClient(info.UUID, in.Lc, in.Handlers, in.ConnBuilder.Build(), addr)
		log.Info().Msgf("Rpc client started, server: %s", addr)
	}

	return nil
}
