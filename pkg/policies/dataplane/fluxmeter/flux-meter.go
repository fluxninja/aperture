package fluxmeter

import (
	"context"
	"path"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"

	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
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
	// The path in status registry for concurrency control status.
	fluxMeterStatusRoot = "concurrency_control"

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
	iface.Policy
	histMetrics    map[flowcontrolv1.DecisionType]prometheus.Histogram
	selector       *selectorv1.Selector
	fluxMeterProto *policylangv1.FluxMeter
	fluxMeterName  string
	buckets        []float64
}

// NewFluxMeterOptions creates fluxmeter for usage in dataplane and also returns its fx options.
func NewFluxMeterOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry *status.Registry,
) (fx.Option, error) {
	registryPath := path.Join(fluxMeterStatusRoot, key.String())
	wrapperMessage := &configv1.FluxMeterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.FluxMeter == nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal flux meter config wrapper")
		return fx.Options(), err
	}
	fluxMeterProto := wrapperMessage.FluxMeter

	fluxMeter := &FluxMeter{
		fluxMeterProto: fluxMeterProto,
		Policy:         wrapperMessage,
		histMetrics:    make(map[flowcontrolv1.DecisionType]prometheus.Histogram),
		fluxMeterName:  wrapperMessage.FluxmeterName,
		selector:       fluxMeterProto.GetSelector(),
		buckets:        fluxMeterProto.GetHistogramBuckets(),
	}

	return fx.Options(
			fx.Invoke(fluxMeter.setup),
		),
		nil
}

func (fluxMeter *FluxMeter) setup(lc fx.Lifecycle, prometheusRegistry *prometheus.Registry) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			decisionTypes := []flowcontrolv1.DecisionType{
				flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED,
				flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED,
				flowcontrolv1.DecisionType_DECISION_TYPE_UNSPECIFIED,
			}
			for _, decisionType := range decisionTypes {
				// Initialize a prometheus histogram metric
				histMetric := prometheus.NewHistogram(prometheus.HistogramOpts{
					Name:    metrics.FluxMeterMetricName,
					Buckets: fluxMeter.buckets,
					// TODO flux-meter-refactor: remove metricID, instead use policyName, fluxMeterName, policyHash
					ConstLabels: prometheus.Labels{
						metrics.PolicyNameLabel:    fluxMeter.GetPolicyName(),
						metrics.FluxMeterNameLabel: fluxMeter.GetFluxMeterName(),
						metrics.PolicyHashLabel:    fluxMeter.GetPolicyHash(),
						metrics.DecisionTypeLabel:  decisionType.String(),
					},
				})
				fluxMeter.histMetrics[decisionType] = histMetric
				// Register metric with Prometheus
				err := prometheusRegistry.Register(histMetric)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to register metric with Prometheus registry for FluxMeter %s", fluxMeter.fluxMeterName)
					return err
				}
			}

			// Register metric with PCA
			err := engineAPI.RegisterFluxMeter(fluxMeter)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to register FluxMeter %s with PolicyConfigAPI", fluxMeter.fluxMeterName)
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			// Unregister metric with PCA
			err := engineAPI.UnregisterFluxMeter(fluxMeter)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to unregister FluxMeter %s with PolicyConfigAPI", fluxMeter.fluxMeterName)
			}
			for _, histMetric := range fluxMeter.histMetrics {
				// Unregister metrics with Prometheus
				unregistered := prometheusRegistry.Unregister(histMetric)
				if !unregistered {
					log.Error().Err(err).Msgf("Failed to unregister metric %s with Prometheus registry", fluxMeter.fluxMeterName)
				}
			}

			return err
		},
	})
}

// GetSelector returns the selector.
func (fluxMeter *FluxMeter) GetSelector() *selectorv1.Selector {
	return fluxMeter.selector
}

// GetFluxMeterProto returns the flux meter proto.
func (fluxMeter *FluxMeter) GetFluxMeterProto() *policylangv1.FluxMeter {
	return fluxMeter.fluxMeterProto
}

// GetFluxMeterName returns the metric name.
func (fluxMeter *FluxMeter) GetFluxMeterName() string {
	return fluxMeter.fluxMeterName
}

// GetFluxMeterID returns the flux meter ID.
func (fluxMeter *FluxMeter) GetFluxMeterID() iface.FluxMeterID {
	return iface.FluxMeterID{
		PolicyName:    fluxMeter.GetPolicyName(),
		FluxMeterName: fluxMeter.GetFluxMeterName(),
		PolicyHash:    fluxMeter.GetPolicyHash(),
	}
}

// GetHistogram returns the histogram.
func (fluxMeter *FluxMeter) GetHistogram(decisionType flowcontrolv1.DecisionType) prometheus.Histogram {
	return fluxMeter.histMetrics[decisionType]
}

// GetBuckets returns the buckets.
func (fluxMeter *FluxMeter) GetBuckets() []float64 {
	return fluxMeter.buckets
}
