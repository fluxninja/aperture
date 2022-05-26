package fluxmeter

import (
	"context"
	"path"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"

	configv1 "aperture.tech/aperture/api/gen/proto/go/aperture/common/config/v1"
	policylangv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"aperture.tech/aperture/pkg/agentinfo"
	"aperture.tech/aperture/pkg/config"
	etcdclient "aperture.tech/aperture/pkg/etcd/client"
	etcdwatcher "aperture.tech/aperture/pkg/etcd/watcher"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/notifiers"
	"aperture.tech/aperture/pkg/paths"
	"aperture.tech/aperture/pkg/policies/dataplane/component"
	"aperture.tech/aperture/pkg/policies/dataplane/iface"
	"aperture.tech/aperture/pkg/status"
)

const (
	// The path in status registry for concurrency control status.
	fluxMeterStatusRoot = "concurrency_control"

	// Label Keys for FluxMeter Metrics.
	metricIDLabelKey = "metric_id"

	// FxNameTag is Flux Meter Watcher's Fx Tag.
	FxNameTag = "name:\"flux_meter\""
)

var engineAPI iface.EngineAPI

// fluxMeterModule returns the fx options for dataplane side pieces of concurrency control in the main fx app.
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
	etcdPath := path.Join(paths.FluxMeterConfigPath, paths.AgentGroupPrefix(agentGroupName))
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
	e iface.EngineAPI,
	sr *status.Registry,
	pr *prometheus.Registry,
) error {
	// save policy config api
	engineAPI = e

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: []notifiers.FxOptionsFunc{NewFluxMeterOptions},
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.NewProtobufUnmarshaller,
		},
		StatusPath:         fluxMeterStatusRoot,
		StatusRegistry:     sr,
		PrometheusRegistry: pr,
	}

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// FluxMeter describes single fluxmeter from policy.
type FluxMeter struct {
	component.ComponentAPI
	histMetric     prometheus.Histogram
	selector       *policylangv1.Selector
	fluxMeterProto *policylangv1.FluxMeter
	metricName     string
	metricID       string
	buckets        []float64
}

// NewFluxMeterOptions creates fluxmeter for usage in dataplane and also returns its fx options.
func NewFluxMeterOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry *status.Registry,
) (fx.Option, error) {
	registryPath := path.Join(fluxMeterStatusRoot, key.String())
	wrapperMessage := &configv1.ConfigPropertiesWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.Config == nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal flux meter config wrapper")
		return fx.Options(), err
	}
	fluxMeterProto := &policylangv1.FluxMeter{}
	err = wrapperMessage.Config.UnmarshalTo(fluxMeterProto)
	if err != nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal flux meter")
		return fx.Options(), err
	}

	fluxMeter := &FluxMeter{
		fluxMeterProto: fluxMeterProto,
		ComponentAPI:   wrapperMessage,
	}

	// Original metric name
	fluxMeter.metricName = fluxMeterProto.Name
	// Selector
	fluxMeter.selector = fluxMeterProto.GetSelector()
	// Buckets
	fluxMeter.buckets = fluxMeterProto.GetHistogramBuckets()
	// Metric ID
	fluxMeter.metricID = paths.MetricIDForFluxMeter(fluxMeter.ComponentAPI, fluxMeter.metricName)

	return fx.Options(
			fx.Invoke(fluxMeter.setup),
		),
		nil
}

// TODO (hasit): rename fluxmeter metric name to a static one 'flux_meter'
// with user-provided name as label 'flux_meter_name' with other labels 'policy_name', 'policy_hash'.
// metricID goes away

func (fluxMeter *FluxMeter) setup(lc fx.Lifecycle, prometheusRegistry *prometheus.Registry) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Initialize a prometheus histogram metric
			fluxMeter.histMetric = prometheus.NewHistogram(prometheus.HistogramOpts{
				Name:        fluxMeter.metricName,
				Buckets:     fluxMeter.buckets,
				ConstLabels: prometheus.Labels{metricIDLabelKey: fluxMeter.metricID},
			})
			// Register metric with Prometheus
			err := prometheusRegistry.Register(fluxMeter.histMetric)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to register metric %s with Prometheus registry", fluxMeter.metricName)
				return err
			}

			// Register metric with PCA
			err = engineAPI.RegisterFluxMeter(fluxMeter)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to register FluxMeter %s with PolicyConfigAPI", fluxMeter.metricName)
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			// Unregister metric with PCA
			err := engineAPI.UnregisterFluxMeter(fluxMeter)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to unregister FluxMeter %s with PolicyConfigAPI", fluxMeter.metricName)
			}
			// Unregister metric with Prometheus
			unregistered := prometheusRegistry.Unregister(fluxMeter.histMetric)
			if !unregistered {
				log.Error().Err(err).Msgf("Failed to unregister metric %s with Prometheus registry", fluxMeter.metricName)
			}

			return err
		},
	})
}

// GetSelector returns the selector.
func (fluxMeter *FluxMeter) GetSelector() *policylangv1.Selector {
	return fluxMeter.selector
}

// GetFluxMeterProto returns the flux meter proto.
func (fluxMeter *FluxMeter) GetFluxMeterProto() *policylangv1.FluxMeter {
	return fluxMeter.fluxMeterProto
}

// GetMetricName returns the metric name.
func (fluxMeter *FluxMeter) GetMetricName() string {
	return fluxMeter.metricName
}

// GetMetricID returns the metric ID.
func (fluxMeter *FluxMeter) GetMetricID() string {
	return fluxMeter.metricID
}

// GetHistogram returns the histogram.
func (fluxMeter *FluxMeter) GetHistogram() prometheus.Histogram {
	return fluxMeter.histMetric
}

// GetBuckets returns the buckets.
func (fluxMeter *FluxMeter) GetBuckets() []float64 {
	return fluxMeter.buckets
}
