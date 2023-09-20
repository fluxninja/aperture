package agentfunctions

import (
	"go.uber.org/fx"

	afconfig "github.com/fluxninja/aperture/v2/pkg/agent-functions/config"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	etcd "github.com/fluxninja/aperture/v2/pkg/etcd/writer"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	grpcclient "github.com/fluxninja/aperture/v2/pkg/net/grpc"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/rpc"
)

const (
	rpcEtcdWatcher = "rpc-etcd-watcher"
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
	etcdwatcher.Constructor{
		Name: rpcEtcdWatcher,
	}.Annotate(),
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
	EtcdWatcher  notifiers.Watcher
	EtcdClient   *etcdclient.Client
	ConnBuilder  grpcclient.ClientConnectionBuilder `name:"agent-functions"`
}

// RegisterClient registers a client which will allow calling agent functions from controller.
func RegisterClient(in RegisterClientIn) error {
	var config afconfig.AgentFunctionsConfig
	if err := in.Unmarshaller.UnmarshalKey(afconfig.Key, &config); err != nil {
		return err
	}

	etcdWriter := *etcd.NewWriter(&in.EtcdClient.KVWrapper)

	for _, addr := range config.Endpoints {
		rpc.RegisterEtcdClient(info.UUID, in.Lc, in.Handlers, in.EtcdWatcher, etcdWriter, addr)
		log.Info().Msgf("Rpc client started, server: %s", addr)
	}

	return nil
}
