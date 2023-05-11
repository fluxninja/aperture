package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	concurrencyv3 "go.etcd.io/etcd/client/v3/concurrency"
	namespacev3 "go.etcd.io/etcd/client/v3/namespace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/etcd"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/panichandler"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// Module is a fx module that provides etcd client.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(ProvideClient),
	)
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
	namespace              = "aperture"
)

// ClientIn holds parameters for ProvideClient.
type ClientIn struct {
	fx.In

	Unmarshaller config.Unmarshaller
	Lifecycle    fx.Lifecycle
	Shutdowner   fx.Shutdowner
	Logger       *log.Logger
}

// Client is a wrapper around etcd client v3. It provides interfaces rooted by a namespace in etcd.
type Client struct {
	KV      clientv3.KV
	Watcher clientv3.Watcher
	Lease   clientv3.Lease
	Client  *clientv3.Client
	Session *concurrencyv3.Session
	LeaseID clientv3.LeaseID
}

// ProvideClient creates a new Etcd Client and provides it via Fx.
func ProvideClient(in ClientIn) (*Client, error) {
	var config etcd.EtcdConfig

	if err := in.Unmarshaller.UnmarshalKey(defaultClientConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize etcd client configuration!")
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	etcdClient := &Client{}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			tlsConfig, tlsErr := config.ClientTLSConfig.GetTLSConfig()
			if tlsErr != nil {
				log.Error().Err(tlsErr).Msg("Failed to get TLS config")
				cancel()
				return tlsErr
			}
			log.Info().Msg("Initializing etcd client")
			cli, err := clientv3.New(clientv3.Config{
				Endpoints: config.Endpoints,
				Context:   ctx,
				TLS:       tlsConfig,
				Username:  config.Username,
				Password:  config.Password,
				Logger:    zap.New(log.NewZapAdapter(in.Logger, "etcd-client"), zap.IncreaseLevel(zap.WarnLevel)),
			})
			if err != nil {
				log.Error().Err(err).Msg("Unable to initialize etcd client")
				cancel()
				return err
			}

			if cli.Username != "" && cli.Password != "" {
				if _, err = cli.AuthEnable(ctx); err != nil {
					log.Error().Err(err).Msg("Unable to enable auth of the etcd cluster")
					cancel()
					return err
				}
			}
			etcdClient.Client = cli

			// namespace the client
			cli.Lease = namespacev3.NewLease(cli.Lease, namespace)
			etcdClient.Lease = cli.Lease
			cli.KV = namespacev3.NewKV(cli.KV, namespace)
			etcdClient.KV = cli.KV
			cli.Watcher = namespacev3.NewWatcher(cli.Watcher, namespace)
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

	return etcdClient, nil
}
