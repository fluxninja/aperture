package fluxmeter

import (
	"context"
	"path"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
)

const (
	// FxNameTag is Flux Meter Watcher's Fx Tag.
	FxNameTag = "name:\"flux_meter\""
)

var engineAPI iface.Engine

// fluxMeterModule returns the fx options for flowcontrol side pieces of concurrency control in the main fx app.
func fluxMeterModule() fx.Option {
	return fx.Options(
		// Tag the watcher so that other modules can find it.
		fx.Provide(
			fx.Annotate(
				provideWatcher,
				fx.ResultTags(FxNameTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupFluxMeterModule,
				fx.ParamTags(FxNameTag),
			),
		),
	)
}

// provideWatcher provides pointer to flux meter watcher.
func provideWatcher(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	// Get Agent Group from host info gatherer
	agentGroupName := ai.GetAgentGroup()
	// Scope the sync to the agent group.
	etcdPath := path.Join(paths.FluxMeterConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}

// setupFluxMeterModule sets up the flux meter module in the main fx app.
func setupFluxMeterModule(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	sr status.Registry,
	pr *prometheus.Registry,
) error {
	// save policy config api
	engineAPI = e

	reg := sr.Child("resource", "flux_meters")

	fmf := &fluxMeterFactory{
		registry: reg,
	}

	fxDriver, err := notifiers.NewFxDriver(reg, pr,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{fmf.newFluxMeterOptions},
	)
	if err != nil {
		return err
	}

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// FluxMeter describes single fluxmeter.
type FluxMeter struct {
	registry      status.Registry
	flowSelector  *policylangv1.FlowSelector
	histMetricVec *prometheus.HistogramVec
	fluxMeterName string
	attributeKey  string
	buckets       []float64
}

type fluxMeterFactory struct {
	registry status.Registry
}

// NewFluxMeterOptions creates fluxmeter for usage in flowcontrol and also returns its fx options.
func (fluxMeterFactory *fluxMeterFactory) newFluxMeterOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := fluxMeterFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.FluxMeterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.FluxMeter == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal flux meter config wrapper")
		return fx.Options(), err
	}
	fluxMeterProto := wrapperMessage.FluxMeter

	buckets := make([]float64, 0)
	switch fluxMeterProto.GetHistogramBuckets().(type) {
	case *policylangv1.FluxMeter_LinearBuckets_:
		if linearBuckets := fluxMeterProto.GetLinearBuckets(); linearBuckets != nil {
			buckets = append(buckets, prometheus.LinearBuckets(
				linearBuckets.GetStart(), linearBuckets.GetWidth(), int(linearBuckets.GetCount()))...)
		}
	case *policylangv1.FluxMeter_ExponentialBuckets_:
		if exponentialBuckets := fluxMeterProto.GetExponentialBuckets(); exponentialBuckets != nil {
			buckets = append(buckets, prometheus.ExponentialBuckets(
				exponentialBuckets.GetStart(), exponentialBuckets.GetFactor(), int(exponentialBuckets.GetCount()))...)
		}
	case *policylangv1.FluxMeter_ExponentialBucketsRange_:
		if exponentialBucketsRange := fluxMeterProto.GetExponentialBucketsRange(); exponentialBucketsRange != nil {
			buckets = append(buckets, prometheus.ExponentialBucketsRange(
				exponentialBucketsRange.GetMin(), exponentialBucketsRange.GetMax(), int(exponentialBucketsRange.GetCount()))...)
		}
	default:
		if defaultBuckets := fluxMeterProto.GetStaticBuckets(); defaultBuckets != nil {
			buckets = append(buckets, defaultBuckets.Buckets...)
		}
	}

	fluxMeter := &FluxMeter{
		fluxMeterName: wrapperMessage.FluxMeterName,
		attributeKey:  fluxMeterProto.AttributeKey,
		flowSelector:  fluxMeterProto.GetFlowSelector(),
		buckets:       buckets,
		registry:      reg,
	}

	return fx.Options(
			fx.Invoke(fluxMeter.setup),
		),
		nil
}

func (fluxMeter *FluxMeter) setup(lc fx.Lifecycle, prometheusRegistry *prometheus.Registry) {
	logger := fluxMeter.registry.GetLogger()
	metricLabels := make(map[string]string)
	metricLabels[metrics.FluxMeterNameLabel] = fluxMeter.GetFluxMeterName()

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Initialize a prometheus histogram metric
			fluxMeter.histMetricVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name:        metrics.FluxMeterMetricName,
				Buckets:     fluxMeter.buckets,
				ConstLabels: prometheus.Labels{metrics.FluxMeterNameLabel: fluxMeter.fluxMeterName},
			}, []string{
				metrics.DecisionTypeLabel,
				metrics.StatusCodeLabel,
				metrics.FlowStatusLabel,
				metrics.ValidLabel,
			})
			// Register metric with Prometheus
			err := prometheusRegistry.Register(fluxMeter.histMetricVec)
			if err != nil {
				logger.Error().Err(err).Msgf("Failed to register metric %+v with Prometheus registry", fluxMeter.histMetricVec)
				return err
			}

			// Register metric with PCA
			err = engineAPI.RegisterFluxMeter(fluxMeter)
			if err != nil {
				logger.Error().Err(err).Msgf("Failed to register FluxMeter %s with EngineAPI", fluxMeter.fluxMeterName)
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			var errMulti error
			// Unregister metric with PCA
			err := engineAPI.UnregisterFluxMeter(fluxMeter)
			if err != nil {
				logger.Error().Err(err).Msgf("Failed to unregister FluxMeter %s with EngineAPI", fluxMeter.fluxMeterName)
				errMulti = multierr.Append(errMulti, err)
			}

			// Unregister metric with Prometheus
			unregistered := prometheusRegistry.Unregister(fluxMeter.histMetricVec)
			if !unregistered {
				logger.Error().Err(err).Msgf("Failed to unregister metric %+v with Prometheus registry", fluxMeter.histMetricVec)
			}

			return errMulti
		},
	})
}

// GetFlowSelector returns the selector.
func (fluxMeter *FluxMeter) GetFlowSelector() *policylangv1.FlowSelector {
	return fluxMeter.flowSelector
}

// GetFluxMeterName returns the metric name.
func (fluxMeter *FluxMeter) GetFluxMeterName() string {
	return fluxMeter.fluxMeterName
}

// GetAttributeKey returns the attribute key.
func (fluxMeter *FluxMeter) GetAttributeKey() string {
	return fluxMeter.attributeKey
}

// GetFluxMeterID returns the flux meter ID.
func (fluxMeter *FluxMeter) GetFluxMeterID() iface.FluxMeterID {
	return iface.FluxMeterID{
		FluxMeterName: fluxMeter.GetFluxMeterName(),
	}
}

// GetHistogram returns the histogram.
func (fluxMeter *FluxMeter) GetHistogram(labels map[string]string) prometheus.Observer {
	logger := fluxMeter.registry.GetLogger()
	fluxMeterHistogram, err := fluxMeter.histMetricVec.GetMetricWith(labels)
	if err != nil {
		logger.Warn().Err(err).Msg("Getting latency histogram")
		return nil
	}
	return fluxMeterHistogram
}
