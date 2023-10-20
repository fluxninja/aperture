package etcd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	clientv3 "go.etcd.io/etcd/client/v3"
	concurrencyv3 "go.etcd.io/etcd/client/v3/concurrency"
	namespacev3 "go.etcd.io/etcd/client/v3/namespace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/etcd"
	"github.com/fluxninja/aperture/v2/pkg/log"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// Module is a fx module that provides etcd client.
func Module() fx.Option {
	return fx.Provide(ProvideClient)
}

// ConfigOverride can be provided by an extension to provide parts of etcd client config directly.
type ConfigOverride struct {
	PerRPCCredentials credentials.PerRPCCredentials // optional
	Namespace         string                        // required
	OverriderName     string                        // who is providing the override, for logs
	Endpoints         []string                      // required
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
	Shutdowner     fx.Shutdowner
	Logger         *log.Logger
	ConfigOverride *ConfigOverride `optional:"true"`
}

// Client is a wrapper around etcd client.
type Client struct {
	KV      clientv3.KV
	Watcher clientv3.Watcher
	Lease   clientv3.Lease
	Client  *clientv3.Client
	Session *concurrencyv3.Session
	LeaseID clientv3.LeaseID
}

// ProvideClient creates a new etcd Client and provides it via Fx.
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

	zapLogLevel, err := zapcore.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid etcd log level")
	}

	ctx, cancel := context.WithCancel(context.Background())

	etcdClient := Client{}

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

			// Workaround for https://github.com/fluxninja/cloud/issues/10613
			logger := in.Logger.Hook(zerolog.HookFunc(
				func(e *zerolog.Event, _ zerolog.Level, msg string) {
					if msg == "lease keepalive response queue is full; dropping response send" {
						// This log pollutes the logs, but is harmless otherwise. The
						// queue isn't used in any meaningful way and this
						// warning doesn't harm the lease itself. Just ignore
						// it for now, until the root cause is fixed.
						e.Discard()
					}
				},
			))

			log.Info().Msg("Initializing etcd client")
			cli, err := clientv3.New(clientv3.Config{
				Endpoints: config.Endpoints,
				Context:   ctx,
				TLS:       tlsConfig,
				Username:  config.Username,
				Password:  config.Password,
				Logger: zap.New(
					log.NewZapAdapter(logger, "etcd-client"),
					zap.IncreaseLevel(zapLogLevel),
				),
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

			etcdClient.Client = cli

			// namespace the client
			cli.Lease = namespacev3.NewLease(cli.Lease, config.Namespace)
			etcdClient.Lease = cli.Lease
			cli.KV = namespacev3.NewKV(cli.KV, config.Namespace)
			etcdClient.KV = cli.KV
			cli.Watcher = namespacev3.NewWatcher(cli.Watcher, config.Namespace)
			etcdClient.Watcher = cli.Watcher

			// Create a new Session
			session, err := concurrencyv3.NewSession(etcdClient.Client, concurrencyv3.WithTTL((int)(config.LeaseTTL.AsDuration().Seconds())))
			if err != nil {
				log.Error().Err(err).Msg("Unable to create a new session")
				cancel()
				return err
			}
			etcdClient.Session = session
			// save the lease id
			etcdClient.LeaseID = session.Lease()
			// A goroutine to check if the session is expired
			panichandler.Go(func() {
				// wait for the context to be done or session to be closed
				select {
				case <-ctx.Done():
					// regular shutdown
				case <-session.Done():
					log.Error().Msg("Etcd session is done, request shutdown")
					utils.Shutdown(in.Shutdowner)
				}
			})

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
