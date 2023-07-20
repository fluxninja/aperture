package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	namespacev3 "go.etcd.io/etcd/client/v3/namespace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/etcd"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Module is a fx module that provides etcd client.
func Module() fx.Option {
	return fx.Provide(
		ProvideClient,
		ProvideSession,
		ProvideSessionScopedKV,
	)
}

// ConfigOverride can be provided by an extension to provide parts of etcd client config directly.
type ConfigOverride struct {
	Namespace         string                        // required
	Endpoints         []string                      // required
	PerRPCCredentials credentials.PerRPCCredentials // optional
	OverriderName     string                        // who is providing the override, for logs
}

const (
	// swagger:operation POST /etcd common-configuration Etcd
	// ---
	// x-fn-config-env: true
	// parameters:
	// - in: body
	//   schema:
	//     $ref: "#/definitions/EtcdConfig"
	defaultClientConfigKey = "etcd"
)

// ClientIn holds parameters for ProvideClient.
type ClientIn struct {
	fx.In

	Unmarshaller   config.Unmarshaller
	Lifecycle      fx.Lifecycle
	Logger         *log.Logger
	ConfigOverride *ConfigOverride `optional:"true"`
}

// Client is a wrapper around etcd client v3. It provides interfaces rooted by a namespace in etcd.
//
// Client.Client is nil before OnStart.
type Client struct {
	*clientv3.Client
	KVWrapper KVWrapper // wraps the same KV as Client
	// hack: This field is here only so that it can be propagated from config
	// to Session.
	leaseTTL config.Duration
}

// KVWrapper wraps clientv3.KV, can be used when wanting to depend on clientv3.KV
// already before OnStart.
//
// KVWrapper.KV is nil before OnStart.
//
// Note: This is not named just KV not to break .KV field access.
type KVWrapper struct {
	clientv3.KV
}

// ProvideClient creates a new Etcd Client and provides it via Fx.
func ProvideClient(in ClientIn) (*Client, error) {
	var config etcd.EtcdConfig

	if err := in.Unmarshaller.UnmarshalKey(defaultClientConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize etcd client configuration!")
		return nil, err
	}

	if in.ConfigOverride != nil {
		if config.Namespace != "aperture" {
			log.Warn().Msg("ignoring etcd.namespace")
		}
		config.Namespace = in.ConfigOverride.Namespace

		if len(config.Endpoints) != 0 {
			log.Warn().Msg("ignoring etcd.endpoints")
		}
		config.Endpoints = in.ConfigOverride.Endpoints

		log.Info().Msgf("etcd endpoints and namespace set by %s", in.ConfigOverride.OverriderName)
	}

	if len(config.Endpoints) == 0 {
		return nil, fmt.Errorf("no etcd endpoints provided")
	}

	ctx, cancel := context.WithCancel(context.Background())

	etcdClient := Client{
		leaseTTL: config.LeaseTTL,
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(startCtx context.Context) error {
			tlsConfig, tlsErr := config.ClientTLSConfig.GetTLSConfig()
			if tlsErr != nil {
				log.Error().Err(tlsErr).Msg("Failed to get TLS config")
				cancel()
				return tlsErr
			}

			var dialOptions []grpc.DialOption

			if in.ConfigOverride != nil && in.ConfigOverride.PerRPCCredentials != nil {
				dialOptions = append(
					dialOptions,
					grpc.WithPerRPCCredentials(in.ConfigOverride.PerRPCCredentials),
				)
			}

			log.Info().Msg("Initializing etcd client")
			cli, err := clientv3.New(clientv3.Config{
				Endpoints:   config.Endpoints,
				Context:     ctx,
				TLS:         tlsConfig,
				Username:    config.Username,
				Password:    config.Password,
				Logger:      zap.New(log.NewZapAdapter(in.Logger, "etcd-client"), zap.IncreaseLevel(zap.WarnLevel)),
				DialOptions: dialOptions,
			})
			if err != nil {
				log.Error().Err(err).Msg("Unable to initialize etcd client")
				cancel()
				return err
			}

			if cli.Username != "" && cli.Password != "" {
				if _, err = cli.AuthEnable(startCtx); err != nil {
					log.Error().Err(err).Msg("Unable to enable auth of the etcd cluster")
					cancel()
					return err
				}
			}

			if config.Namespace != "" {
				// namespace the client
				cli.Lease = namespacev3.NewLease(cli.Lease, config.Namespace)
				cli.KV = namespacev3.NewKV(cli.KV, config.Namespace)
				cli.Watcher = namespacev3.NewWatcher(cli.Watcher, config.Namespace)
			}

			etcdClient.Client = cli
			etcdClient.KVWrapper.KV = cli.KV

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info().Msg("Closing etcd connections")
			cancel()
			return etcdClient.Client.Close()
		},
	})
	return &etcdClient, nil
}
