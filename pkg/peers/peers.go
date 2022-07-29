package peers

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"os"
	"path"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v2"

	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/peers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/status"
)

const (
	// swagger:operation POST /peer_discovery common-configuration PeerDiscovery
	// ---
	// parameters:
	// - in: body
	//   schema:
	//     "$ref": "#/definitions/PeerDiscoveryConfig"
	defaultKey   = "peer_discovery"
	watcherFxTag = "peer-discovery-watcher"
)

var (
	peerDiscoverySyncPath = path.Join(config.DefaultTempDirectory, "peers")
	etcdPath              = path.Join("/peers")
)

// PeerDiscoveryConfig holds configuration for Agent Peer Discovery.
// swagger:model
type PeerDiscoveryConfig struct {
	// Network address of aperture server to advertise to peers - this address should be reachable from other agents. Used for nat traversal when provided.
	AdvertisementAddr string `json:"advertisement_addr" validate:"omitempty,hostname_port"`
}

// Constructor holds fields to create and configure PeerDiscovery.
type Constructor struct {
	Key           string
	DefaultConfig PeerDiscoveryConfig
	Service       string
}

// Module is a fx module that creates peer directory and provides peer discovery.
func (constructor Constructor) Module() fx.Option {
	_ = os.MkdirAll(peerDiscoverySyncPath, fs.ModePerm)
	return fx.Options(
		fx.Provide(constructor.providePeerDiscovery),
	)
}

// PeerDiscoveryPrefix is the prefix for peer discovery service.
type PeerDiscoveryPrefix string

// PeerDiscoveryIn holds parameters for newPeerDiscovery.
type PeerDiscoveryIn struct {
	fx.In
	Lifecycle      fx.Lifecycle
	Unmarshaller   config.Unmarshaller
	Client         *etcdclient.Client
	Listener       *listener.Listener
	StatusRegistry *status.Registry
	Prefix         PeerDiscoveryPrefix
	Watchers       PeerWatchers `group:"peer-watchers"`
}

func (constructor Constructor) providePeerDiscovery(in PeerDiscoveryIn) (*PeerDiscovery, error) {
	var configKey string
	if constructor.Key == "" {
		configKey = defaultKey
	} else {
		configKey = constructor.Key
	}

	var cfg PeerDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		return nil, err
	}

	pd, err := NewPeerDiscovery(string(in.Prefix), in.Client, in.Watchers)
	if err != nil {
		return nil, err
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			hostname := info.Hostname
			var advertiseAddr string
			if cfg.AdvertisementAddr != "" {
				advertiseAddr = cfg.AdvertisementAddr
			} else {
				// Must be called in start stage
				addr := in.Listener.GetListener().Addr().String()
				_, port, err := net.SplitHostPort(addr)
				if err != nil {
					return err
				}
				advertiseAddr = hostname + ":" + port
			}
			log.Debug().Str("advertise_addr", advertiseAddr).Msg("advertise addr")

			err := pd.Start()
			if err != nil {
				return err
			}
			err = pd.registerSelf(ctx, advertiseAddr)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, e error
			e = pd.deregisterSelf(ctx)
			if e != nil {
				merr = multierr.Combine(merr, e)
			}

			e = pd.Stop()
			if e != nil {
				merr = multierr.Combine(merr, e)
			}
			return merr
		},
	})

	return pd, nil
}

// PeerDiscovery holds fields to manage peer discovery.
type PeerDiscovery struct {
	lock         sync.RWMutex
	peers        map[string]*peersv1.PeerInfo
	selfPeer     *peersv1.PeerInfo
	client       *etcdclient.Client
	etcdWatcher  notifiers.Watcher
	selfKey      string
	etcdPath     string
	peerNotifier notifiers.PrefixNotifier
	watchers     PeerWatchers
}

// NewPeerDiscovery creates a new PeerDiscovery.
func NewPeerDiscovery(prefix string,
	client *etcdclient.Client,
	watchers PeerWatchers,
) (*PeerDiscovery, error) {
	var err error
	pd := &PeerDiscovery{
		peers:    make(map[string]*peersv1.PeerInfo),
		watchers: watchers,
		etcdPath: path.Join(etcdPath, prefix),
		client:   client,
		selfPeer: &peersv1.PeerInfo{
			Services: make(map[string]string),
		},
	}

	// create and start etcdwatcher to track peers and sync them to disk
	pd.etcdWatcher, err = etcdwatcher.NewWatcher(client, pd.etcdPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to create etcd watcher")
		return nil, err
	}

	pd.peerNotifier = &notifiers.UnmarshalPrefixNotifier{
		GetUnmarshallerFunc: config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
		UnmarshalNotifyFunc: pd.updatePeer,
	}

	return pd, nil
}

// registerSelf registers self to etcd.
func (pd *PeerDiscovery) registerSelf(ctx context.Context, advertiseAddr string) error {
	var err error
	hostname := info.Hostname

	pd.selfPeer.Address = advertiseAddr
	pd.selfPeer.Hostname = hostname

	pd.selfKey = path.Join(pd.etcdPath, pd.selfPeer.Hostname)

	// register
	log.Debug().Str("key", pd.selfKey).Msg("self registering in peer discovery table")
	b, err := yaml.Marshal(pd.selfPeer)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal peer info")
		return err
	}

	_, err = pd.client.KV.Put(clientv3.WithRequireLeader(ctx),
		pd.selfKey, string(b), clientv3.WithLease(pd.client.LeaseID))

	return err
}

// deregisterSelf deregisters self from etcd.
func (pd *PeerDiscovery) deregisterSelf(ctx context.Context) error {
	_, err := pd.client.KV.Delete(clientv3.WithRequireLeader(ctx), pd.selfKey)
	return err
}

// Start starts peer discovery.
func (pd *PeerDiscovery) Start() error {
	if err := pd.etcdWatcher.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start etcd watcher")
		return err
	}

	if err := pd.etcdWatcher.AddPrefixNotifier(pd.peerNotifier); err != nil {
		log.Error().Err(err).Msg("failed to add directory notifier")
		return err
	}

	return nil
}

// Stop stops peer discovery.
func (pd *PeerDiscovery) Stop() error {
	var merr, err error
	err = pd.etcdWatcher.RemovePrefixNotifier(pd.peerNotifier)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove prefix notifier")
		merr = multierr.Combine(merr, err)
	}

	err = pd.etcdWatcher.Stop()
	if err != nil {
		log.Error().Err(err).Msg("failed to stop etcd watcher")
		merr = multierr.Combine(merr, err)
	}

	return merr
}

// GetPeers returns all the peer info that are added to PeerDiscovery.
func (pd *PeerDiscovery) GetPeers() []*peersv1.PeerInfo {
	pd.lock.RLock()
	defer pd.lock.RUnlock()

	peers := make([]*peersv1.PeerInfo, 0)
	for _, peer := range pd.peers {
		peers = append(peers, peer)
	}

	return peers
}

// RegisterService accepts a name, full address (host:port) and adds to the list of services in PeerDiscovery.
func (pd *PeerDiscovery) RegisterService(name string, address string) {
	pd.lock.Lock()
	defer pd.lock.Unlock()

	pd.selfPeer.Services[name] = address
}

func (pd *PeerDiscovery) addPeer(peer *peersv1.PeerInfo) {
	defer pd.watchers.OnPeerAdded(peer)
	pd.lock.Lock()
	defer pd.lock.Unlock()

	pd.peers[peer.Address] = peer
}

// GetPeer returns the peer info in the PeerDiscovery with the given address.
func (pd *PeerDiscovery) GetPeer(address string) (*peersv1.PeerInfo, error) {
	pd.lock.RLock()
	defer pd.lock.RUnlock()

	peer, ok := pd.peers[address]
	if !ok {
		return nil, errors.New("peer not found")
	}

	return peer, nil
}

// GetPeerKeys returns all the peer keys that are added to PeerDiscovery.
func (pd *PeerDiscovery) GetPeerKeys() []string {
	pd.lock.RLock()
	defer pd.lock.RUnlock()

	keys := make([]string, 0)
	for key := range pd.peers {
		keys = append(keys, key)
	}

	return keys
}

func (pd *PeerDiscovery) removePeer(address string) {
	var peer *peersv1.PeerInfo
	defer func() {
		if peer != nil {
			pd.watchers.OnPeerRemoved(peer)
		}
	}()

	pd.lock.Lock()
	defer pd.lock.Unlock()

	peer = pd.peers[address]
	delete(pd.peers, address)
}

func (pd *PeerDiscovery) updatePeer(event notifiers.Event, unmarshaller config.Unmarshaller) {
	log.Debug().Str("event", event.String()).Msg("Updating peer")
	if event.Type == notifiers.Write {
		var peer peersv1.PeerInfo
		if err := unmarshaller.UnmarshalKey("", &peer); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal peer info")
			return
		}
		pd.addPeer(&peer)
	} else if event.Type == notifiers.Remove {
		key := string(event.Key)
		addr := path.Base(key)
		pd.removePeer(addr)
	}
}
