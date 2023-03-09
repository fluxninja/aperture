// +kubebuilder:validation:Optional
package agentfunctions

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/rpc"
)

// ConfigKey is the key for agentfunctions configuration.
const ConfigKey = "agent_functions"

// Module provides rpc client for agent functions.
var Module = fx.Options(
	// FIXME(krdln) Do we actually need a separate grpc client for each module?
	grpcclient.ClientConstructor{
		Name:      "agent-functions",
		ConfigKey: ConfigKey + ".client.grpc",
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

// Config is configuration for agent functions.
// swagger:model
// +kubebuilder:object:generate=true
type Config struct {
	// RPC servers to connect to (which will be able to call agent functions)
	Endpoints []string `json:"endpoints,omitempty" validate:"omitempty,dive,omitempty"`

	// Network client configuration
	ClientConfig ClientConfig `json:"client"`
}

// ClientConfig is configuration for network clients used by agent-functions.
// swagger:model
// +kubebuilder:object:generate=true
type ClientConfig struct {
	// GRPC client settings.
	GRPCClient grpcclient.GRPCClientConfig `json:"grpc"`
}

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
	var config Config
	if err := in.Unmarshaller.UnmarshalKey(ConfigKey, &config); err != nil {
		return err
	}

	for _, addr := range config.Endpoints {
		rpc.RegisterStreamClient(info.UUID, in.Lc, in.Handlers, in.ConnBuilder.Build(), addr)
		log.Info().Msgf("Rpc client started, server: %s", addr)
	}

	return nil
}
