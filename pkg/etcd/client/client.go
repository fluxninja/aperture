// +kubebuilder:validation:Optional
package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	concurrencyv3 "go.etcd.io/etcd/client/v3/concurrency"
	namespacev3 "go.etcd.io/etcd/client/v3/namespace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/hashicorp/go-multierror"
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
	LeaseTTL config.Duration `json:"lease_ttl" validate:"gte=1s" default:"60s"`
	// Authentication
	Username string `json:"username"`
	Password string `json:"password"`
	// Client TLS configuration
	ClientTLSConfig tlsconfig.ClientTLSConfig `json:"tls"`
	// List of Etcd server endpoints
	Endpoints []string `json:"endpoints" validate:"gt=0,dive,hostname_port|url|fqdn"`
}

// ClientIn holds parameters for ProvideClient.
type ClientIn struct {
	fx.In

	Unmarshaller config.Unmarshaller
	Lifecycle    fx.Lifecycle
	Shutdowner   fx.Shutdowner
	Logger       *log.Logger
	AgentInfo    *agentinfo.AgentInfo
}

// Client is a wrapper around etcd client v3. It provides interfaces rooted by a namespace in etcd.
type Client struct {
	KV       clientv3.KV
	Watcher  clientv3.Watcher
	Lease    clientv3.Lease
	Client   *clientv3.Client
	Election *concurrencyv3.Election
	LeaseID  clientv3.LeaseID
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
			// save the lease id
			etcdClient.LeaseID = session.Lease()
			// Create an election for this client
			etcdClient.Election = concurrencyv3.NewElection(session, "/election/"+in.AgentInfo.GetAgentGroup())
			// A goroutine to do leader election
			panichandler.Go(func() {
				// try to elect a leader
				err := etcdClient.Election.Campaign(ctx, info.GetHostInfo().Hostname)
				if err != nil {
					log.Error().Err(err).Msg("Unable to elect a leader")
					shutdownErr := in.Shutdowner.Shutdown()
					if shutdownErr != nil {
						log.Error().Err(shutdownErr).Msg("Error on invoking shutdown")
					}
				}
				// wait for the context to be done or session to be closed
				select {
				case <-ctx.Done():
					// regular shutdown
				case <-session.Done():
					log.Error().Msg("Etcd session is done, request shutdown")
					shutdownErr := in.Shutdowner.Shutdown()
					if shutdownErr != nil {
						log.Error().Err(shutdownErr).Msg("Error on invoking shutdown")
					}
				}
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info().Msg("Closing etcd connections")
			cancel()
			var merr error
			// resign from the election if are the leader
			if etcdClient.Election.Key() != "" {
				stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
				err := etcdClient.Election.Resign(stopCtx)
				stopCancel()
				if err != nil {
					log.Error().Err(err).Msg("Unable to resign from the election")
					merr = multierror.Append(merr, err)
				}
			}
			err := etcdClient.Client.Close()
			if err != nil {
				merr = multierror.Append(merr, err)
			}
			return merr
		},
	})

	return etcdClient, nil
}

// IsLeader returns true if the current node is the leader.
func (c *Client) IsLeader() bool {
	return c.Election.Key() != ""
}
