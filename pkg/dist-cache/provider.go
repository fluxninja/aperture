package distcache

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"strconv"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	distcachev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/distcache/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	dcconfig "github.com/fluxninja/aperture/v2/pkg/dist-cache/config"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/aperture/v2/pkg/peers"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

const (
	defaultKey                 = "dist_cache"
	olricMemberlistServiceName = "olric-memberlist"
	distCacheMetricsJobName    = "dist-cache-scrape-metrics"
)

// Module provides a new DistCache FX module.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(DistCacheConstructor{ConfigKey: defaultKey}.ProvideDistCache),
		grpcgateway.RegisterHandler{Handler: distcachev1.RegisterDistCacheServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(RegisterDistCacheService),
	)
}

// DistCacheConstructorIn holds parameters of ProvideDistCache.
type DistCacheConstructorIn struct {
	fx.In
	PeerDiscovery      *peers.PeerDiscovery
	PrometheusRegistry *prometheus.Registry
	LivenessMultiJob   *jobs.MultiJob `name:"liveness.service"`
	Unmarshaller       config.Unmarshaller
	Lifecycle          fx.Lifecycle
	Shutdowner         fx.Shutdowner
	Logger             *log.Logger
}

// DistCacheConstructor holds fields to create an instance of DistCache.
type DistCacheConstructor struct {
	ConfigKey     string
	DefaultConfig dcconfig.DistCacheConfig
}

// ProvideDistCache creates a new instance of distributed cache.
// It also hooks in the service discovery plugin.
func (constructor DistCacheConstructor) ProvideDistCache(in DistCacheConstructorIn) (*DistCache, error) {
	defaultConfig := constructor.DefaultConfig
	if err := in.Unmarshaller.UnmarshalKey(constructor.ConfigKey, &defaultConfig); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize configuration of DistCache")
		return nil, err
	}

	memberlistEnv := "lan"
	oc := olricconfig.New(memberlistEnv)

	oc.WriteQuorum = 1
	oc.ReadQuorum = 1
	oc.MemberCountQuorum = 1
	oc.ReadRepair = false

	oc.ReplicaCount = defaultConfig.ReplicaCount
	if defaultConfig.SyncReplication {
		oc.ReplicationMode = olricconfig.SyncReplicationMode
	} else {
		oc.ReplicationMode = olricconfig.AsyncReplicationMode
	}

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

	serviceName := fmt.Sprintf("%s-%s", olricMemberlistServiceName, info.GetVersionInfo().Version)
	oc.ServiceDiscovery = map[string]interface{}{
		"plugin": &ServiceDiscovery{
			discovery:   in.PeerDiscovery,
			addr:        memberlistAddr,
			serviceName: serviceName,
		},
	}

	startChan := make(chan struct{})
	oc.Started = func() {
		log.Info().Msg("DistCache started")
		startChan <- struct{}{}
	}

	o, err := olric.New(oc)
	if err != nil {
		return nil, err
	}

	dc := NewDistCache(oc, o, newDistCacheMetrics(), in.Shutdowner)

	job := jobs.NewBasicJob(distCacheMetricsJobName, dc.scrapeMetrics)
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Register metrics with Prometheus.
			err := dc.metrics.registerMetrics(in.PrometheusRegistry)
			if err != nil {
				log.Error().Err(err).Msg("Failed to register distcache metrics with Prometheus registry")
				return err
			}

			panichandler.Go(func() {
				defer func() {
					utils.Shutdown(in.Shutdowner)
				}()

				startErr := dc.olric.Start()
				if startErr != nil {
					log.Error().Err(startErr).Msg("Failed to start distcache")
				}
			})

			// wait for olric to start by waiting on startChan or until ctx is canceled
			select {
			case <-ctx.Done():
				return errors.New("olric failed to start")
			case <-startChan:
			}

			err = in.LivenessMultiJob.RegisterJob(job)
			if err != nil {
				log.Error().Err(err).Msg("Failed to register distcache scrape metrics job with jobGroup")
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := in.LivenessMultiJob.DeregisterJob(distCacheMetricsJobName)
			if err != nil {
				log.Error().Err(err).Msg("Failed to deregister distcache scrape metrics job with jobGroup")
				return err
			}

			err = dc.olric.Shutdown(ctx)
			if err != nil {
				return err
			}

			// Unregister metrics with Prometheus.
			err = dc.metrics.unregisterMetrics(in.PrometheusRegistry)
			if err != nil {
				log.Error().Err(err).Msg("Failed to unregister distcache metrics with Prometheus registry")
				return err
			}
			return nil
		},
	})

	return dc, nil
}

// RegisterDistCacheService registers the handler on grpc.Server.
func RegisterDistCacheService(handler *DistCache, server *grpc.Server, healthsrv *health.Server) error {
	distcachev1.RegisterDistCacheServiceServer(server, handler)

	healthsrv.SetServingStatus("aperture.distcache.v1.DistCacheService", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info().Msg("DistCache Stats handler registered")
	return nil
}
