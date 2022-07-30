package distcache

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/peers"
)

const (
	defaultKey                 = "dist_cache"
	olricMemberlistServiceName = "olric-memberlist"
)

// Module provides a new DistCache FX module.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(DistCacheConstructor{Key: defaultKey}.ProvideDistCache),
	)
}

// swagger:operation POST /dist_cache common-configuration DistCache
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/DistCacheConfig"

// DistCacheConfig configures distributed cache that holds per-label counters in distributed rate limiters.
// swagger:model
type DistCacheConfig struct {
	// BindAddr denotes the address that Olric will bind to for communication with other Olric nodes.
	BindAddr string `json:"bind_addr" default:":3320" validate:"hostname_port"`
	// ReplicaCount is 1 by default.
	ReplicaCount int `json:"replica_count" default:"1"`
	// Address to bind mememberlist server to.
	MemberlistConfigBindAddr string `json:"memberlist_config_bind_addr" default:":3322" validate:"hostname_port"`
	// Address of memberlist to advertise to other cluster members. Used for nat traversal if provided.
	MemberlistConfigAdvertiseAddr string `json:"memberlist_config_advertise_addr" validate:"omitempty,hostname_port"`
}

// DistCache wraps an Olric instance along with its config for further reference.
type DistCache struct {
	sync.Mutex
	Config *olricconfig.Config
	Olric  *olric.Olric
}

// AddDMapCustomConfig adds a named DMap config into DistCache's config.
// If a custom config with the name does not exist, it is added. If it already exists, it is overwritten.
func (dc *DistCache) AddDMapCustomConfig(name string, dmapConfig olricconfig.DMap) {
	dc.Config.DMaps.Custom[name] = dmapConfig
}

// RemoveDMapCustomConfig removes a named DMap config from DistCache's config.
func (dc *DistCache) RemoveDMapCustomConfig(name string) {
	delete(dc.Config.DMaps.Custom, name)
}

// DistCacheConstructorIn holds parameters of ProvideDistCache.
type DistCacheConstructorIn struct {
	fx.In
	PeerDiscovery *peers.PeerDiscovery
	Unmarshaller  config.Unmarshaller
	Lifecycle     fx.Lifecycle
	Shutdowner    fx.Shutdowner
}

// DistCacheConstructor holds fields to create an instance of *DistCache.
type DistCacheConstructor struct {
	Key           string
	DefaultConfig DistCacheConfig
}

// ProvideDistCache creates a new instance of Olric distributed cache in embded mode.
// It also hooks in the service discovery plugin.
func (constructor DistCacheConstructor) ProvideDistCache(in DistCacheConstructorIn) (*DistCache, error) {
	config := constructor.DefaultConfig
	if err := in.Unmarshaller.UnmarshalKey(constructor.Key, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize configuration of DistCache")
		return nil, err
	}

	dc := &DistCache{}

	memberlistEnv := "lan"
	oc := olricconfig.New(memberlistEnv)
	oc.ServiceDiscovery = map[string]interface{}{
		"plugin": &ServiceDiscovery{
			discovery: in.PeerDiscovery,
		},
	}
	oc.ReplicaCount = config.ReplicaCount
	oc.WriteQuorum = 1
	oc.ReadQuorum = 1
	oc.MemberCountQuorum = 1
	oc.DMaps.Custom = make(map[string]olricconfig.DMap)

	bindAddr, port, err := net.SplitHostPort(config.BindAddr)
	if err != nil {
		log.Error().Err(err).Msg("Unable to split bind_addr")
		return nil, err
	}
	bindPort, _ := strconv.Atoi(port)

	if bindAddr == "" {
		bindAddr = info.Hostname
	}
	oc.BindAddr = bindAddr
	oc.BindPort = bindPort

	memberlistBindAddr, p, err := net.SplitHostPort(config.MemberlistConfigBindAddr)
	if err != nil {
		log.Error().Err(err).Msg("Unable to split memberlist bind address")
		return nil, err
	}
	memberlistBindPort, _ := strconv.Atoi(p)

	if memberlistBindAddr == "" {
		memberlistBindAddr = info.Hostname
	}
	oc.MemberlistConfig.BindAddr = memberlistBindAddr
	oc.MemberlistConfig.BindPort = memberlistBindPort
	memberlistAddr := config.MemberlistConfigBindAddr

	if config.MemberlistConfigAdvertiseAddr != "" {
		advertiseAddr, p, e := net.SplitHostPort(config.MemberlistConfigAdvertiseAddr)
		if e != nil {
			log.Error().Err(e).Msg("Unable to split memberlist advertise address")
			return nil, e
		}
		advertisePort, _ := strconv.Atoi(p)
		oc.MemberlistConfig.AdvertiseAddr = advertiseAddr
		oc.MemberlistConfig.AdvertisePort = advertisePort
		memberlistAddr = config.MemberlistConfigAdvertiseAddr
	}

	in.PeerDiscovery.RegisterService(olricMemberlistServiceName, memberlistAddr)

	startChan := make(chan struct{})
	oc.Started = func() {
		log.Info().Msg("Olric started")
		startChan <- struct{}{}
	}

	dc.Config = oc

	o, err := olric.New(dc.Config)
	if err != nil {
		return nil, err
	}

	dc.Olric = o

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msg("Starting OTEL Collector")
			panichandler.Go(func() {
				err := dc.Olric.Start()
				if err != nil {
					log.Error().Err(err).Msg("Failed to start olric")
				}
				_ = in.Shutdowner.Shutdown()
			})
			// wait for olric to start by waiting on startChan until ctx is canceled
			select {
			case <-ctx.Done():
				return errors.New("olric failed to start")
			case <-startChan:
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := dc.Olric.Shutdown(ctx)
			if err != nil {
				return err
			}
			return nil
		},
	})

	return dc, nil
}
