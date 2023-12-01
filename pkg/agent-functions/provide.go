package agentfunctions

import (
	"go.uber.org/fx"

	afconfig "github.com/fluxninja/aperture/v2/pkg/agent-functions/config"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/etcd/transport"
	"github.com/fluxninja/aperture/v2/pkg/info"
	grpcclient "github.com/fluxninja/aperture/v2/pkg/net/grpc"
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
		NewCacheHandler,
	),
	fx.Invoke(
		RegisterEtcdTransport,
		RegisterControlPointsHandler,
		RegisterPreviewHandler,
		RegisterCacheHandlers,
	),
)

// RegisterClientIn are parameters for InvokeClient function.
type RegisterClientIn struct {
	fx.In
	Lc                  fx.Lifecycle
	Unmarshaller        config.Unmarshaller
	ConnBuilder         grpcclient.ClientConnectionBuilder `name:"agent-functions"`
	EtcdTransportServer *transport.EtcdTransportServer
	EtcdClient          *etcdclient.Client
}

// RegisterEtcdTransport registers a server on the etcd transport.
func RegisterEtcdTransport(in RegisterClientIn) {
	transport.RegisterWatcher(in.Lc, in.EtcdTransportServer, info.Hostname)
}
