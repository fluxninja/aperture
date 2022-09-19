package fluxmeter

import (
	"context"
	"path"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/status"
)

const (
	// FxNameTag is Flux Meter Watcher's Fx Tag.
	FxNameTag = "name:\"flux_meter\""
)

var engineAPI iface.Engine

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
	e iface.Engine,
	sr status.Registry,
	pr *prometheus.Registry,
) error {
	// save policy config api
	engineAPI = e

	reg := sr.Child("flux_meters")

	fmf := &fluxMeterFactory{
		statusRegistry: reg,
	}

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: []notifiers.FxOptionsFunc{fmf.newFluxMeterOptions},
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.NewProtobufUnmarshaller,
		},
		StatusRegistry:     reg,
		PrometheusRegistry: pr,
	}

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// FluxMeter describes single fluxmeter.
type FluxMeter struct {
	selector      *selectorv1.Selector
	histMetricVec *prometheus.HistogramVec
	fluxMeterName string
	attributeKey  string
	buckets       []float64
}

type fluxMeterFactory struct {
	statusRegistry status.Registry
}

// NewFluxMeterOptions creates fluxmeter for usage in dataplane and also returns its fx options.
func (fluxMeterFactory *fluxMeterFactory) newFluxMeterOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	wrapperMessage := &wrappersv1.FluxMeterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.FluxMeter == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		log.Warn().Err(err).Msg("Failed to unmarshal flux meter config wrapper")
		return fx.Options(), err
	}
	fluxMeterProto := wrapperMessage.FluxMeter

	fluxMeter := &FluxMeter{
		fluxMeterName: wrapperMessage.FluxMeterName,
		attributeKey:  fluxMeterProto.AttributeKey,
		selector:      fluxMeterProto.GetSelector(),
		buckets:       fluxMeterProto.GetHistogramBuckets(),
	}

	return fx.Options(
			fx.Invoke(fluxMeter.setup),
		),
		nil
}

func (fluxMeter *FluxMeter) setup(lc fx.Lifecycle, prometheusRegistry *prometheus.Registry) {
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
				metrics.FeatureStatusLabel,
			})
			// Register metric with Prometheus
			err := prometheusRegistry.Register(fluxMeter.histMetricVec)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to register metric %+v with Prometheus registry", fluxMeter.histMetricVec)
				return err
			}

			// Register metric with PCA
			err = engineAPI.RegisterFluxMeter(fluxMeter)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to register FluxMeter %s with EngineAPI", fluxMeter.fluxMeterName)
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			var errMulti error
			// Unregister metric with PCA
			err := engineAPI.UnregisterFluxMeter(fluxMeter)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to unregister FluxMeter %s with EngineAPI", fluxMeter.fluxMeterName)
				errMulti = multierr.Append(errMulti, err)
			}

			// Unregister metric with Prometheus
			unregistered := prometheusRegistry.Unregister(fluxMeter.histMetricVec)
			if !unregistered {
				log.Error().Err(err).Msgf("Failed to unregister metric %+v with Prometheus registry", fluxMeter.histMetricVec)
			}

			return errMulti
		},
	})
}

// GetSelector returns the selector.
func (fluxMeter *FluxMeter) GetSelector() *selectorv1.Selector {
	return fluxMeter.selector
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
func (fluxMeter *FluxMeter) GetHistogram(decisionType flowcontrolv1.CheckResponse_DecisionType, statusCode string, featureStatus string) prometheus.Observer {
	labels := make(map[string]string)
	labels[metrics.DecisionTypeLabel] = decisionType.String()
	labels[metrics.StatusCodeLabel] = statusCode
	labels[metrics.FeatureStatusLabel] = featureStatus

	fluxMeterHistogram, err := fluxMeter.histMetricVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return nil
	}

	return fluxMeterHistogram
}
