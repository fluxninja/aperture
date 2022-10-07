// +kubebuilder:validation:Optional
package distcache

import (
	"context"
	"errors"
	stdlog "log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/peers"
)

const (
	defaultKey                 = "distcache"
	olricMemberlistServiceName = "olric-memberlist"
	distCacheMetricsJobName    = "scrape-metrics"
)

// Module provides a new DistCache FX module.
func Module() fx.Option {
	return fx.Options(
		jobs.JobGroupConstructor{Name: "distcache", Key: defaultKey}.Annotate(),
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
	Config     *olricconfig.Config
	Olric      *olric.Olric
	Metrics    *DistCacheMetrics
	jobGroup   *jobs.JobGroup
	metricsJob *jobs.MultiJob
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

func (dc *DistCache) scrapeMetrics() error {
	stats, err := dc.Olric.Stats()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to scrape Olric statistics")
		return err
	}

	memberID := stats.Member.ID
	memberName := stats.Member.Name
	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.DistCacheMemberIDLabel] = strconv.FormatUint(memberID, 10)
	metricLabels[metrics.DistCacheMemberNameLabel] = memberName

	if dc.Metrics == nil {
		dc.Metrics = newDistCacheMetrics()
	}
	entriesTotalGauge, err := dc.Metrics.EntriesTotal.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract entries total gauge metric from olric instance: %v", err)
	} else {
		entriesTotalGauge.Set(float64(stats.DMaps.EntriesTotal))
	}

	deleteHitsGauge, err := dc.Metrics.DeleteHits.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract delete hits gauge metric from olric instance: %v", err)
	} else {
		deleteHitsGauge.Set(float64(stats.DMaps.DeleteHits))
	}

	deleteMissesGague, err := dc.Metrics.DeleteMisses.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract delete misses gauge metric from olric instance: %v", err)
	} else {
		deleteMissesGague.Set(float64(stats.DMaps.DeleteMisses))
	}

	getMissesGague, err := dc.Metrics.GetMisses.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract get misses gauge metric from olric instance: %v", err)
	} else {
		getMissesGague.Set(float64(stats.DMaps.GetMisses))
	}

	getHitsGague, err := dc.Metrics.GetHits.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract get hits gauge metric from olric instance: %v", err)
	} else {
		getHitsGague.Set(float64(stats.DMaps.GetHits))
	}

	evictedTotalGague, err := dc.Metrics.EvictedTotal.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract evicted total gauge metric from olric instance: %v", err)
	} else {
		evictedTotalGague.Set(float64(stats.DMaps.EvictedTotal))
	}
	return nil
}

// DistCacheConstructorIn holds parameters of ProvideDistCache.
type DistCacheConstructorIn struct {
	fx.In
	PeerDiscovery      *peers.PeerDiscovery
	Unmarshaller       config.Unmarshaller
	Lifecycle          fx.Lifecycle
	JobGroup           *jobs.JobGroup `name:"distcache"`
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
	defaultConfig := constructor.DefaultConfig
	if err := in.Unmarshaller.UnmarshalKey(constructor.ConfigKey, &defaultConfig); err != nil {
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
	oc.ReplicaCount = defaultConfig.ReplicaCount
	oc.WriteQuorum = 1
	oc.ReadQuorum = 1
	oc.MemberCountQuorum = 1
	oc.DMaps.Custom = make(map[string]olricconfig.DMap)
	oc.Logger = stdlog.New(&OlricLogWriter{Logger: in.Logger}, "", 0)

	bindAddr, port, err := net.SplitHostPort(defaultConfig.BindAddr)
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

	memberlistBindAddr, p, err := net.SplitHostPort(defaultConfig.MemberlistBindAddr)
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

	if defaultConfig.MemberlistAdvertiseAddr != "" {
		advertiseAddr, p, e := net.SplitHostPort(defaultConfig.MemberlistAdvertiseAddr)
		if e != nil {
			log.Error().Err(e).Msg("Unable to split memberlist advertise address")
			return nil, e
		}
		advertisePort, _ := strconv.Atoi(p)
		oc.MemberlistConfig.AdvertiseAddr = advertiseAddr
		oc.MemberlistConfig.AdvertisePort = advertisePort
		memberlistAddr = defaultConfig.MemberlistAdvertiseAddr
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
	dc.Metrics = newDistCacheMetrics()
	dc.jobGroup = in.JobGroup
	multiJob := jobs.NewMultiJob(dc.jobGroup.GetStatusRegistry().Child("distcache"), nil, nil)
	dc.metricsJob = multiJob

	job := &jobs.BasicJob{
		JobBase: jobs.JobBase{
			JobName: distCacheMetricsJobName,
		},
		JobFunc: func(ctx context.Context) (proto.Message, error) {
			select {
			case <-ctx.Done():
				return &emptypb.Empty{}, nil
			default:
				err = dc.scrapeMetrics()
				return &emptypb.Empty{}, err
			}
		},
	}
	err = dc.metricsJob.RegisterJob(job)
	if err != nil {
		return nil, err
	}

	jobConfig := jobs.JobConfig{
		InitialDelay:     config.MakeDuration(time.Millisecond * 100),
		ExecutionPeriod:  config.MakeDuration(time.Millisecond * 500),
		ExecutionTimeout: config.MakeDuration(time.Millisecond * 1000),
		InitiallyHealthy: false,
	}
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := dc.Metrics.registerMetrics(in.PrometheusRegistry)
			if err != nil {
				return err
			}
			err = dc.jobGroup.RegisterJob(job, jobConfig)
			if err != nil {
				return err
			}
			log.Info().Msg("Starting OTEL Collector")
			panichandler.Go(func() {
				log.Info().Msg("Started OTEL Collector")
				err = dc.Olric.Start()
				dc.jobGroup.TriggerJob(distCacheMetricsJobName)
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
			var multiErr error
			err := dc.Metrics.unregisterMetrics(in.PrometheusRegistry)
			if err != nil {
				multiErr = multierr.Append(multiErr, err)
			}

			err = dc.Olric.Shutdown(ctx)
			if err != nil {
				multiErr = multierr.Append(multiErr, err)
			}
			return multiErr
		},
	})

	return dc, nil
}
