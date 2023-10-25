package etcd

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lukejoshuapark/infchan"
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
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// Module is a fx module that provides etcd client.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(ProvideClient),
	)
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

	// ElectionPathFxTag is the fx tag for the election path.
	ElectionPathFxTag = "etcd.election-path"
	// EnforceLeaderOnlyFxTag is the fx tag for the enforce leader only flag.
	EnforceLeaderOnlyFxTag = "etcd.enforce-leader-only"
)

// ClientIn holds parameters for ProvideClient.
type ClientIn struct {
	fx.In
	Unmarshaller      config.Unmarshaller
	Lifecycle         fx.Lifecycle
	Shutdowner        fx.Shutdowner
	Logger            *log.Logger
	ConfigOverride    *ConfigOverride `optional:"true"`
	ElectionPath      string          `name:"etcd.election-path"`
	EnforceLeaderOnly bool            `name:"etcd.enforce-leader-only"`
}

const (
	put       = 0
	del       = 1
	delPrefix = 2
	bootstrap = 3
)

type operation struct {
	key        string
	value      string
	opts       []clientv3.OpOption
	opType     int
	withLease  bool
	withExpiry bool
	ttl        int64
}

// ElectionWatcher is used for tracking changes to election.
type ElectionWatcher interface {
	OnLeaderStart()
	OnLeaderStop()
}

// Client is a wrapper around etcd client.
type Client struct {
	kv                   clientv3.KV
	watcher              clientv3.Watcher
	lease                clientv3.Lease
	opChannel            infchan.Channel[operation]
	client               *clientv3.Client
	readyChannel         chan bool
	electionWatchers     map[ElectionWatcher]struct{}
	cache                map[string]string
	leaseID              clientv3.LeaseID
	electionWatcherMutex sync.Mutex
	cacheMutex           sync.Mutex
	isLeader             atomic.Bool
	bootstrapPending     atomic.Bool
	enforceLeaderOnly    bool
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

	var dialOptions []grpc.DialOption

	if in.ConfigOverride != nil && in.ConfigOverride.PerRPCCredentials != nil {
		dialOptions = append(
			dialOptions,
			grpc.WithPerRPCCredentials(in.ConfigOverride.PerRPCCredentials),
		)
	}

	tlsConfig, tlsErr := config.ClientTLSConfig.GetTLSConfig()
	if tlsErr != nil {
		log.Error().Err(tlsErr).Msg("Failed to get TLS config")
		return nil, tlsErr
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

	ctx, cancel := context.WithCancel(context.Background())

	clientConfig := clientv3.Config{
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
	}

	etcdClient := Client{
		readyChannel:      make(chan bool),
		opChannel:         infchan.NewChannel[operation](),
		electionWatchers:  make(map[ElectionWatcher]struct{}),
		cache:             make(map[string]string),
		enforceLeaderOnly: in.EnforceLeaderOnly,
	}
	etcdClient.bootstrapPending.Store(true)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// A goroutine keeps trying etcd connection in the background
			panichandler.Go(etcdClient.mainLoopFn(clientConfig, in.Shutdowner, ctx, cancel, config.Namespace, (int)(config.LeaseTTL.AsDuration().Seconds()), in.ElectionPath, wg))

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info().Msg("Closing etcd connections")
			cancel()
			wg.Wait()
			etcdClient.closeOpChannel()

			return etcdClient.client.Close()
		},
	})
	return &etcdClient, nil
}

func (etcdClient *Client) closeOpChannel() {
	// drain the opChannel.Out
	for {
		select {
		case <-etcdClient.opChannel.Out():
		default:
			etcdClient.opChannel.Close()
			return
		}
	}
}

func (etcdClient *Client) mainLoopFn(clientConfig clientv3.Config, shutdowner fx.Shutdowner, ctx context.Context, cancel context.CancelFunc, namespace string, ttl int, electionPath string, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()
		for {
			log.Info().Msg("Initializing etcd client")
			cli, err := clientv3.New(clientConfig)
			if err != nil {
				log.Error().Err(err).Msg("Unable to initialize etcd client, try again")
				// wait 5 seconds before retrying
				time.Sleep(5 * time.Second)
				continue
			}

			if cli.Username != "" && cli.Password != "" {
				if _, err = cli.AuthEnable(ctx); err != nil {
					log.Error().Err(err).Msg("Unable to enable auth of the etcd cluster")
					utils.Shutdown(shutdowner)
				}
			}

			etcdClient.client = cli
			// namespace the client
			cli.Lease = namespacev3.NewLease(cli.Lease, namespace)
			etcdClient.lease = cli.Lease
			cli.KV = namespacev3.NewKV(cli.KV, namespace)
			etcdClient.kv = cli.KV
			cli.Watcher = namespacev3.NewWatcher(cli.Watcher, namespace)
			etcdClient.watcher = cli.Watcher

			if !etcdClient.enforceLeaderOnly {
				// close the ready channel
				close(etcdClient.readyChannel)
			}
			break
		}

	SESSION_LOOP:
		for {
			// Create a new Session
			session, err := concurrencyv3.NewSession(etcdClient.client, concurrencyv3.WithTTL(ttl))
			if err != nil {
				log.Error().Err(err).Msg("Unable to create a new session")
				// wait 5 seconds before retrying
				time.Sleep(5 * time.Second)
				continue
			}
			// save the lease id
			etcdClient.leaseID = session.Lease()

			// write loop wait group and context
			writeLoopWaitGroup := sync.WaitGroup{}
			writeLoopCtx, writeLoopCancel := context.WithCancel(ctx)

			campaignCtx, campaignCancel := context.WithCancel(context.Background())
			// A goroutine keeps trying to campaign for leadership
			panichandler.Go(etcdClient.campaignLoopFn(campaignCtx, session, electionPath, &writeLoopWaitGroup, writeLoopCtx))

			if !etcdClient.enforceLeaderOnly {
				etcdClient.launchWriteLoopRoutine(&writeLoopWaitGroup, writeLoopCtx)
			}
			// wait for the context to be done or session to be closed
			select {
			case <-ctx.Done():
				// regular shutdown
				campaignCancel()
				writeLoopCancel()
				writeLoopWaitGroup.Wait()
				break SESSION_LOOP
			case <-session.Done():
				etcdClient.bootstrapPending.Store(true)
				campaignCancel()
				writeLoopCancel()
				writeLoopWaitGroup.Wait()
				if etcdClient.isLeader.Load() {
					etcdClient.informElectionWatcher(false)
				}
				if etcdClient.enforceLeaderOnly {
					// Shutdown
					log.Info().Msg("Etcd session is done, shutting down")
					utils.Shutdown(shutdowner)
					break SESSION_LOOP
				} else {
					log.Error().Msg("Etcd session is done, re-create it")
				}
			}
		}
	}
}

func (etcdClient *Client) campaignLoopFn(ctx context.Context, session *concurrencyv3.Session, electionPath string, writeLoopWaitGroup *sync.WaitGroup, writeLoopCtx context.Context) func() {
	return func() {
		for {
			// Create an election for this client
			election := concurrencyv3.NewElection(session, electionPath)
			// Campaign for leadership
			err := election.Campaign(ctx, info.GetHostInfo().Uuid)
			// Check if canceled
			if ctx.Err() != nil {
				return
			}
			if err != nil {
				log.Error().Err(err).Msg("Unable to elect a leader, try again")
				continue
			}
			// This is the leader
			etcdClient.informElectionWatcher(true)
			if etcdClient.enforceLeaderOnly {
				// close the ready channel
				close(etcdClient.readyChannel)
				etcdClient.launchWriteLoopRoutine(writeLoopWaitGroup, writeLoopCtx)
			}

			log.Info().Msg("Node is now the leader")
			break
		}
	}
}

func (etcdClient *Client) launchWriteLoopRoutine(writeLoopWaitGroup *sync.WaitGroup, writeLoopCtx context.Context) {
	// bootstrap writes
	etcdClient.bootstrapPending.Store(false)
	etcdClient.opChannel.In() <- operation{
		opType: bootstrap,
	}

	writeLoopWaitGroup.Add(1)
	// write loop
	panichandler.Go(etcdClient.writeLoopFn(writeLoopCtx, writeLoopWaitGroup))
}

func (etcdClient *Client) writeLoopFn(ctx context.Context, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()
		for {
			select {
			case op := <-etcdClient.opChannel.Out():
				if op.withLease && op.opType == put {
					op.opts = append(op.opts, clientv3.WithLease(etcdClient.leaseID))
				}
				if op.withExpiry && op.opType == put {
					lease, err := etcdClient.lease.Grant(ctx, op.ttl)
					if err != nil {
						log.Error().Err(err).Msg("failed to get lease for request")
						continue
					}
					op.opts = append(op.opts, clientv3.WithLease(lease.ID))
				}

				var err error
				switch op.opType {
				case bootstrap:
					etcdClient.cacheMutex.Lock()
					cacheCopy := make(map[string]string)
					for key, value := range etcdClient.cache {
						cacheCopy[key] = value
					}
					etcdClient.cacheMutex.Unlock()
					for key, value := range cacheCopy {
						_, err = etcdClient.kv.Put(clientv3.WithRequireLeader(ctx), key, value)
					}
				case put:
					_, err = etcdClient.kv.Put(clientv3.WithRequireLeader(ctx), op.key, op.value, op.opts...)
				case del:
					_, err = etcdClient.kv.Delete(clientv3.WithRequireLeader(ctx), op.key, op.opts...)
				case delPrefix:
					op.opts = append(op.opts, clientv3.WithPrefix())
					_, err = etcdClient.kv.Delete(clientv3.WithRequireLeader(ctx), op.key, op.opts...)
				}
				if err != nil {
					log.Error().Err(err).Msg("failed to write to etcd")
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

// Watch starts a watch with etcd and returns a channel that will receive events.
func (etcdClient *Client) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) (clientv3.WatchChan, error) {
	// wait for the client to be ready
	select {
	case <-etcdClient.readyChannel:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return etcdClient.watcher.Watch(ctx, key, opts...), nil
}

// Get gets the value for the given key.
func (etcdClient *Client) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	// wait for the client to be ready
	select {
	case <-etcdClient.readyChannel:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return etcdClient.kv.Get(ctx, key, opts...)
}

// Txn returns a transaction.
func (etcdClient *Client) Txn(ctx context.Context) (clientv3.Txn, error) {
	// wait for the client to be ready
	select {
	case <-etcdClient.readyChannel:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return etcdClient.client.Txn(ctx), nil
}

// PutWithExpiry puts the given value for the given key and lease duration.
func (etcdClient *Client) PutWithExpiry(key, val string, leaseTTL int, opts ...clientv3.OpOption) {
	// send the operation to the channel
	etcdClient.opChannel.In() <- operation{
		key:        key,
		value:      val,
		opType:     put,
		opts:       opts,
		withExpiry: true,
		ttl:        (int64)(leaseTTL),
	}
}

// Put puts the given value for the given key with lease.
func (etcdClient *Client) Put(key, val string, opts ...clientv3.OpOption) {
	etcdClient.cacheMutex.Lock()
	defer etcdClient.cacheMutex.Unlock()
	etcdClient.cache[key] = val
	if etcdClient.bootstrapPending.Load() {
		return
	}
	// send the operation to the channel
	etcdClient.opChannel.In() <- operation{
		key:       key,
		value:     val,
		opType:    put,
		opts:      opts,
		withLease: true,
	}
}

// Delete deletes the given key.
func (etcdClient *Client) Delete(key string, opts ...clientv3.OpOption) {
	etcdClient.cacheMutex.Lock()
	defer etcdClient.cacheMutex.Unlock()
	delete(etcdClient.cache, key)
	if etcdClient.bootstrapPending.Load() {
		return
	}
	// etcd robustness TODO: update the local tracker
	// send the operation to the channel
	etcdClient.opChannel.In() <- operation{
		key:    key,
		opType: del,
		opts:   opts,
	}
}

// DeletePrefix deletes all keys with the given prefix.
func (etcdClient *Client) DeletePrefix(prefix string, opts ...clientv3.OpOption) {
	// send the operation to the channel
	etcdClient.opChannel.In() <- operation{
		key:    prefix,
		opType: delPrefix,
		opts:   opts,
	}
}

// PutSync puts the given value for the given key.
func (etcdClient *Client) PutSync(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	// wait for the client to be ready
	select {
	case <-etcdClient.readyChannel:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return etcdClient.kv.Put(ctx, key, val, opts...)
}

// DeleteSync deletes the given key.
func (etcdClient *Client) DeleteSync(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	// wait for the client to be ready
	select {
	case <-etcdClient.readyChannel:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return etcdClient.kv.Delete(ctx, key, opts...)
}

// IsLeader returns true if the current node is the leader.
func (etcdClient *Client) IsLeader() bool {
	return etcdClient.isLeader.Load()
}

// AddElectionWatcher adds a watcher for election changes.
func (etcdClient *Client) AddElectionWatcher(electionWatcher ElectionWatcher) {
	etcdClient.electionWatcherMutex.Lock()
	defer etcdClient.electionWatcherMutex.Unlock()
	etcdClient.electionWatchers[electionWatcher] = struct{}{}
	if etcdClient.IsLeader() {
		electionWatcher.OnLeaderStart()
	}
}

// RemoveElectionWatcher removes a watcher for election changes.
func (etcdClient *Client) RemoveElectionWatcher(electionWatcher ElectionWatcher) {
	etcdClient.electionWatcherMutex.Lock()
	defer etcdClient.electionWatcherMutex.Unlock()
	delete(etcdClient.electionWatchers, electionWatcher)
}

// informElectionWatcher informs all watchers about the current leader status.
func (etcdClient *Client) informElectionWatcher(isLeader bool) {
	etcdClient.electionWatcherMutex.Lock()
	defer etcdClient.electionWatcherMutex.Unlock()
	etcdClient.isLeader.Store(isLeader)
	for watcher := range etcdClient.electionWatchers {
		if isLeader {
			watcher.OnLeaderStart()
		} else {
			watcher.OnLeaderStop()
		}
	}
}
