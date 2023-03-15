package heartbeats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/fluxninja/v1"
	"github.com/fluxninja/aperture/extensions/fluxninja/extconfig"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/discovery/entities"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/etcd/election"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/peers"
	autoscalediscovery "github.com/fluxninja/aperture/pkg/policies/autoscale/kubernetes/discovery"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/status"
)

// Module returns the module for heartbeats.
func Module() fx.Option {
	log.Info().Msg("Loading Heartbeats extension")
	return fx.Options(
		grpcgateway.RegisterHandler{Handler: heartbeatv1.RegisterControllerInfoServiceHandlerFromEndpoint}.Annotate(),
		grpcclient.ClientConstructor{Name: "heartbeats-grpc-client", ConfigKey: extconfig.ExtensionConfigKey + ".client.grpc"}.Annotate(),
		httpclient.ClientConstructor{Name: "heartbeats-http-client", ConfigKey: extconfig.ExtensionConfigKey + ".client.http"}.Annotate(),
		PeersWatcherModule(),
		jobs.JobGroupConstructor{Name: heartbeatsGroup}.Annotate(),
		fx.Provide(provide),
	)
}

// ConstructorIn injects dependencies into the Heartbeats constructor.
type ConstructorIn struct {
	fx.In

	Lifecycle                        fx.Lifecycle
	ExtensionConfig                  *extconfig.FluxNinjaExtensionConfig
	GRPCServer                       *grpc.Server
	JobGroup                         *jobs.JobGroup                     `name:"heartbeats-job-group"`
	GRPClientConnectionBuilder       grpcclient.ClientConnectionBuilder `name:"heartbeats-grpc-client"`
	HTTPClient                       *http.Client                       `name:"heartbeats-http-client"`
	StatusRegistry                   status.Registry
	Entities                         *entities.Entities   `optional:"true"`
	AgentInfo                        *agentinfo.AgentInfo `optional:"true"`
	PeersWatcher                     *peers.PeerDiscovery `name:"fluxninja-peers-watcher" optional:"true"`
	EtcdClient                       *etcdclient.Client
	Election                         *election.Election                        `optional:"true"`
	PolicyFactory                    *controlplane.PolicyFactory               `optional:"true"`
	FlowControlPoints                *cache.Cache[selectors.ControlPointID]    `optional:"true"`
	AutoscaleKubernetesControlPoints autoscalediscovery.AutoScaleControlPoints `optional:"true"`
}

// provide provides a new instance of Heartbeats.
func provide(in ConstructorIn) (*Heartbeats, error) {
	if in.ExtensionConfig.APIKey == "" {
		log.Info().Msg("Heartbeats API key not set, skipping")
		return nil, nil
	}

	installationMode := getInstallationMode()

	heartbeats := newHeartbeats(
		in.JobGroup,
		in.ExtensionConfig,
		in.StatusRegistry,
		in.Entities,
		in.AgentInfo,
		in.PeersWatcher,
		in.PolicyFactory,
		in.Election,
		in.FlowControlPoints,
		in.AutoscaleKubernetesControlPoints,
		installationMode,
	)

	heartbeatv1.RegisterControllerInfoServiceServer(in.GRPCServer, heartbeats)

	runCtx, cancel := context.WithCancel(context.Background())

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := heartbeats.setupControllerInfo(runCtx, in.EtcdClient)
			if err != nil {
				log.Error().Err(err).Msg("Could not read/create controller id in heartbeats")
				return err
			}

			err = heartbeats.start(runCtx, &in)
			if err != nil {
				log.Error().Err(err).Msg("Heartbeats start had an error")
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			heartbeats.stop()
			return nil
		},
	})

	return heartbeats, nil
}

func getInstallationMode() string {
	if isKubernetes() {
		if isKubernetesSidecar() {
			return "KUBERNETES_SIDECAR"
		} else if isKubernetesDaemonSet() {
			return "KUBERNETES_DAEMONSET"
		}
		return "KUBERNETES_POD"
	} else if isLinux() {
		return "LINUX_BARE_METAL"
	} else {
		return "UNKNOWN"
	}
}

func isLinux() bool {
	// Check if running on bare metal Linux
	if _, err := os.Stat("/proc/cpuinfo"); err == nil {
		return true
	}
	return false
}

func isKubernetes() bool {
	_, kubernetesServiceHostExists := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	_, kubernetesServicePortExists := os.LookupEnv("KUBERNETES_SERVICE_PORT")
	return kubernetesServiceHostExists && kubernetesServicePortExists
}

func isKubernetesSidecar() bool {
	_, podNameExists := os.LookupEnv("POD_NAME")
	_, podNamespaceExists := os.LookupEnv("POD_NAMESPACE")
	if podNameExists && podNamespaceExists {
		// Check if running as a sidecar container
		file, err := os.Open("/proc/1/cgroup")
		if err != nil {
			return false
		}
		defer file.Close()

		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil {
			return false
		}

		cgroup := string(buf[:n])
		if strings.Contains(cgroup, "pod") && !strings.Contains(cgroup, "sandbox") {
			return true
		}
	}

	return false
}

func isKubernetesDaemonSet() bool {
	podName, podNameExists := os.LookupEnv("POD_NAME")
	podNamespace, podNamespaceExists := os.LookupEnv("POD_NAMESPACE")
	if podNameExists && podNamespaceExists {
		// Get the pod's metadata
		path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", podNamespace, podName)
		data, err := os.ReadFile(fmt.Sprintf("/var/run/secrets/kubernetes.io/serviceaccount/%s", path))
		if err != nil {
			return false
		}

		// Check if the pod has the daemonset label
		var obj map[string]interface{}
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return false
		}
		metadata, ok := obj["metadata"].(map[string]interface{})
		if !ok {
			return false
		}
		labels, ok := metadata["labels"].(map[string]interface{})
		if !ok {
			return false
		}
		_, daemonsetLabelExists := labels["daemonset"]
		if daemonsetLabelExists {
			return true
		}
	}

	return false
}
