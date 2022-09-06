package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	namespacev3 "go.etcd.io/etcd/client/v3/namespace"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panichandler"
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

// EtcdConfig holds configuration for etcd client.
// swagger:model
// +kubebuilder:object:generate=true
type EtcdConfig struct {
	// Lease time-to-live
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="60s"
	LeaseTTL config.Duration `json:"lease_ttl" validate:"gte=1s" default:"60s"`
	// List of Etcd server endpoints
	//+kubebuilder:validation:Optional
	//+kubebuilder:validate:MinItems=1
	Endpoints []string `json:"endpoints" validate:"gt=0,dive,hostname_port|url|fqdn"`
	// TODO: add auth params
}

// ClientIn holds parameters for ProvideClient.
type ClientIn struct {
	fx.In

	Unmarshaller config.Unmarshaller
	Lifecycle    fx.Lifecycle
	Shutdowner   fx.Shutdowner
}

// Client is a wrapper around etcd client v3.
type Client struct {
	// raw client
	Client *clientv3.Client
	// interfaces rooted by namespace -- use these for all operations instead of the raw client
	KV      clientv3.KV
	Watcher clientv3.Watcher
	Lease   clientv3.Lease
	LeaseID clientv3.LeaseID
}

// ProvideClient creates a new Etcd Client and provides it via Fx.
func ProvideClient(in ClientIn) (*Client, error) {
	var config EtcdConfig

	if err := in.Unmarshaller.UnmarshalKey(defaultClientConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize etcd client configuration!")
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	etcdClient := &Client{}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info().Msg("Initializing etcd client")

			cli, err := clientv3.New(clientv3.Config{
				Endpoints: config.Endpoints,
				Context:   ctx,
			})
			if err != nil {
				log.Error().Err(err).Msg("Unable to initialize etcd client")
				cancel()
				return err
			}
			etcdClient.Client = cli

			etcdClient.Lease = namespacev3.NewLease(etcdClient.Client, namespace)
			etcdClient.KV = namespacev3.NewKV(etcdClient.Client, namespace)
			etcdClient.Watcher = namespacev3.NewWatcher(etcdClient.Client, namespace)

			// Create a lease with etcd for this client, exit app if lease maintenance fails
			resp, err := etcdClient.Lease.Grant(ctx, (int64)(config.LeaseTTL.AsDuration().Seconds()))
			if err != nil {
				log.Error().Err(err).Msg("Unable to grant a lease")
				cancel()
				return err
			}
			// save the lease id
			etcdClient.LeaseID = resp.ID

			// try to keep the lease alive
			keepAlive, err := etcdClient.Lease.KeepAlive(ctx, etcdClient.LeaseID)
			if err != nil || keepAlive == nil {
				log.Error().Err(err).Msg("Unable to keep alive the lease")
			}

			panichandler.Go(func() {
				for ka := range keepAlive {
					if ka != nil {
						continue
					}
					log.Error().Msg("Lease failed, TTL is null")
					break
				}
				select {
				case <-ctx.Done():
					// regular shutdown
				default:
					log.Error().Msg("Request shutdown on lease failure")
					err := in.Shutdowner.Shutdown()
					if err != nil {
						log.Error().Err(err).Msg("Error on invoking shutdown")
					}
				}
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info().Msg("Closing etcd connections")
			cancel()
			// revoke the lease
			_, err := etcdClient.Lease.Revoke(context.Background(), etcdClient.LeaseID)
			if err != nil {
				log.Error().Err(err).Msg("Unable to revoke lease")
			}
			err = etcdClient.Client.Close()
			if err != nil {
				return err
			}
			return nil
		},
	})

	return etcdClient, nil
}
