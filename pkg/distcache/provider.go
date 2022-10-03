// +kubebuilder:validation:Optional
package distcache

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"strconv"
	"sync"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"github.com/prometheus/client_golang/prometheus"
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
		fx.Provide(DistCacheConstructor{ConfigKey: defaultKey}.ProvideDistCache),
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
// +kubebuilder:object:generate=true
type DistCacheConfig struct {
	// BindAddr denotes the address that DistCache will bind to for communication with other peer nodes.
	BindAddr string `json:"bind_addr" default:":3320" validate:"hostname_port"`
	// ReplicaCount is 1 by default.
	ReplicaCount int `json:"replica_count" default:"1"`
	// Address to bind mememberlist server to.
	MemberlistBindAddr string `json:"memberlist_bind_addr" default:":3322" validate:"hostname_port"`
	// Address of memberlist to advertise to other cluster members. Used for nat traversal if provided.
	MemberlistAdvertiseAddr string `json:"memberlist_advertise_addr" validate:"omitempty,hostname_port"`
}

// DistCache is a peer to peer distributed cache.
type DistCache struct {
	sync.Mutex
	Config  *olricconfig.Config
	Olric   *olric.Olric
	Metrics *OlricMetrics
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

func (dc *DistCache) registerMetrics(prometheusRegistry *prometheus.Registry) error {
	for _, m := range dc.Metrics.allMetrics() {
		err := prometheusRegistry.Register(m)
		if err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return fmt.Errorf("unable to register Olric metrics: %w", err)
			}
		}
	}
	return nil
}

func (dc *DistCache) unregisterMetrics(prometheusRegistry *prometheus.Registry) error {
	for _, m := range dc.Metrics.allMetrics() {
		if !prometheusRegistry.Unregister(m) {
			return errors.New("unable to unregister Olric Metrics")
		}
	}
	return nil
}

func (dc *DistCache) scrapeMetrics() {
	stats, err := dc.Olric.Stats()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to scrape Olric statistics")
		return
	}
	dc.Metrics.EntriesTotal.Set(float64(stats.DMaps.EntriesTotal))
	dc.Metrics.DeleteHits.Set(float64(stats.DMaps.DeleteHits))
	dc.Metrics.DeleteMisses.Set(float64(stats.DMaps.DeleteMisses))
	dc.Metrics.GetMisses.Set(float64(stats.DMaps.GetMisses))
	dc.Metrics.GetHits.Set(float64(stats.DMaps.GetHits))
	dc.Metrics.EvictedTotal.Set(float64(stats.DMaps.EvictedTotal))
}

// DistCacheConstructorIn holds parameters of ProvideDistCache.
type DistCacheConstructorIn struct {
	fx.In
	PeerDiscovery      *peers.PeerDiscovery
	Unmarshaller       config.Unmarshaller
	Lifecycle          fx.Lifecycle
	Shutdowner         fx.Shutdowner
	Logger             *log.Logger
	PrometheusRegistry *prometheus.Registry
}

// DistCacheConstructor holds fields to create an instance of *DistCache.
type DistCacheConstructor struct {
	ConfigKey     string
	DefaultConfig DistCacheConfig
}

// ProvideDistCache creates a new instance of distributed cache.
// It also hooks in the service discovery plugin.
func (constructor DistCacheConstructor) ProvideDistCache(in DistCacheConstructorIn) (*DistCache, error) {
	config := constructor.DefaultConfig
	if err := in.Unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
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
	oc.Logger = stdlog.New(&OlricLogWriter{Logger: in.Logger}, "", 0)

	bindAddr, port, err := net.SplitHostPort(config.BindAddr)
	if err != nil {
		log.Error().Err(err).Msg("Unable to split bind_addr")
		return nil, err
	}
	bindPort, _ := strconv.Atoi(port)

	if bindAddr == "" {
		bindAddr = info.LocalIP
	}
	oc.BindAddr = bindAddr
	oc.BindPort = bindPort

	memberlistBindAddr, p, err := net.SplitHostPort(config.MemberlistBindAddr)
	if err != nil {
		log.Error().Err(err).Msg("Unable to split memberlist bind address")
		return nil, err
	}
	memberlistBindPort, _ := strconv.Atoi(p)

	if memberlistBindAddr == "" {
		memberlistBindAddr = info.LocalIP
	}
	oc.MemberlistConfig.BindAddr = memberlistBindAddr
	oc.MemberlistConfig.BindPort = memberlistBindPort
	memberlistAddr := oc.MemberlistConfig.BindAddr + ":" + strconv.Itoa(oc.MemberlistConfig.BindPort)

	if config.MemberlistAdvertiseAddr != "" {
		advertiseAddr, p, e := net.SplitHostPort(config.MemberlistAdvertiseAddr)
		if e != nil {
			log.Error().Err(e).Msg("Unable to split memberlist advertise address")
			return nil, e
		}
		advertisePort, _ := strconv.Atoi(p)
		oc.MemberlistConfig.AdvertiseAddr = advertiseAddr
		oc.MemberlistConfig.AdvertisePort = advertisePort
		memberlistAddr = config.MemberlistAdvertiseAddr
	}

	in.PeerDiscovery.RegisterService(olricMemberlistServiceName, memberlistAddr)

	startChan := make(chan struct{})
	oc.Started = func() {
		log.Info().Msg("DistCache started")
		startChan <- struct{}{}
	}

	dc.Config = oc

	o, err := olric.New(dc.Config)
	if err != nil {
		return nil, err
	}

	dc.Olric = o
	dc.Metrics = newOlricMetrics()

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_ = dc.registerMetrics(in.PrometheusRegistry)
			log.Info().Msg("Starting OTEL Collector")
			panichandler.Go(func() {
				err := dc.Olric.Start()
				if err != nil {
					log.Error().Err(err).Msg("Failed to start olric")
				}
				// TODO: Periodically scrape metrics from DMaps struct
				dc.scrapeMetrics()
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
			_ = dc.unregisterMetrics(in.PrometheusRegistry)
			err := dc.Olric.Shutdown(ctx)
			if err != nil {
				return err
			}
			return nil
		},
	})

	return dc, nil
}
