package flowregulator

import (
	"context"
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"path"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/distcache"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
)

const flowRegulatorStatusRoot = "flow_regulators"

var (
	fxNameTag       = config.NameTag(flowRegulatorStatusRoot)
	metricLabelKeys = []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.DecisionTypeLabel,
		metrics.RegulatorDroppedLabel,
	}
)

// Module returns the fx options for flow regulator.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideWatchers,
				fx.ResultTags(fxNameTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupFlowRegulatorFactory,
				fx.ParamTags(fxNameTag),
			),
		),
	)
}

func provideWatchers(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	agentGroupName := ai.GetAgentGroup()

	etcdPath := path.Join(paths.FlowRegulatorConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type flowRegulatorFactory struct {
	engineAPI            iface.Engine
	registry             status.Registry
	distCache            *distcache.DistCache
	decisionsWatcher     notifiers.Watcher
	dynamicConfigWatcher notifiers.Watcher
	counterVector        *prometheus.CounterVec
	agentGroupName       string
}

func setupFlowRegulatorFactory(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	distCache *distcache.DistCache,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) error {
	agentGroupName := ai.GetAgentGroup()
	etcdPath := path.Join(paths.FlowRegulatorDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	dynamicConfigWatcher, err := etcdwatcher.NewWatcher(etcdClient,
		paths.FlowRegulatorDynamicConfigPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", flowRegulatorStatusRoot)
	counterVector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.FlowRegulatorCounterMetricName,
			Help: "Total number of decisions made by flow regulators.",
		},
		metricLabelKeys,
	)

	factory := &flowRegulatorFactory{
		engineAPI:            e,
		registry:             statusRegistry,
		distCache:            distCache,
		decisionsWatcher:     decisionsWatcher,
		dynamicConfigWatcher: dynamicConfigWatcher,
		counterVector:        counterVector,
		agentGroupName:       agentGroupName,
	}

	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{factory.newFlowRegulatorOptions})
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := prometheusRegistry.Register(factory.counterVector)
			if err != nil {
				return err
			}
			err = decisionsWatcher.Start()
			if err != nil {
				return err
			}
			err = dynamicConfigWatcher.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			var err, merr error
			err = decisionsWatcher.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = dynamicConfigWatcher.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(factory.counterVector) {
				err2 := fmt.Errorf("failed to unregister flow_regulator_counter metric")
				merr = multierr.Append(merr, err2)
			}
			reg.Detach()
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

func (frf *flowRegulatorFactory) newFlowRegulatorOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := frf.registry.GetLogger()
	wrapperMessage := &policysyncv1.FlowRegulatorWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.FlowRegulator == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal flow regulator config")
		return fx.Options(), err
	}

	flowRegulatorProto := wrapperMessage.FlowRegulator
	fr := &flowRegulator{
		Component:         wrapperMessage.GetCommonAttributes(),
		proto:             flowRegulatorProto,
		labelKey:          flowRegulatorProto.GetParameters().GetLabelKey(),
		factory:           frf,
		registry:          reg,
		enableLabelValues: make(map[string]bool),
	}
	fr.name = iface.ComponentKey(fr)

	return fx.Options(
		fx.Invoke(
			fr.setup,
		),
	), nil
}

type flowRegulator struct {
	enableValuesMutex sync.RWMutex
	iface.Component
	registry          status.Registry
	factory           *flowRegulatorFactory
	proto             *policylangv1.FlowRegulator
	enableLabelValues map[string]bool
	name              string
	labelKey          string
	acceptPercentage  float64
}

// Make sure flowRegulator implements iface.Limiter.
var _ iface.Limiter = (*flowRegulator)(nil)

func (fr *flowRegulator) setup(lifecycle fx.Lifecycle) error {
	logger := fr.registry.GetLogger()
	etcdKey := paths.AgentComponentKey(fr.factory.agentGroupName,
		fr.GetPolicyName(),
		fr.GetComponentId())
	// decision notifier
	decisionUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaller,
		fr.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}
	// dynamic config notifier
	dynamicConfigUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	dynamicConfigNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		dynamicConfigUnmarshaller,
		fr.dynamicConfigUpdateCallback,
	)
	if err != nil {
		return err
	}

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = fr.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = fr.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = fr.GetComponentId()
	counterVec := fr.factory.counterVector

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			fr.updateDynamicConfig(fr.proto.GetDefaultConfig())

			// add decisions notifier
			err = fr.factory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}
			// add dynamic config notifier
			err = fr.factory.dynamicConfigWatcher.AddKeyNotifier(dynamicConfigNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add dynamic config notifier")
				return err
			}

			// add to data engine
			err = fr.factory.engineAPI.RegisterFlowRegulator(fr)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to register flow regulator")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error
			deleted := counterVec.DeletePartialMatch(metricLabels)
			if deleted == 0 {
				logger.Warn().Msg("Could not delete flow regulator counter from its metric vector. No traffic to generate metrics?")
			}
			// remove from data engine
			err = fr.factory.engineAPI.UnregisterFlowRegulator(fr)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to unregister rate limiter")
				merr = multierr.Append(merr, err)
			}
			// remove dynamic config notifier
			err = fr.factory.dynamicConfigWatcher.RemoveKeyNotifier(dynamicConfigNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove dynamic config notifier")
				merr = multierr.Append(merr, err)
			}
			// remove decisions notifier
			err = fr.factory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierr.Append(merr, err)
			}
			fr.registry.SetStatus(status.NewStatus(nil, merr))

			return merr
		},
	})
	return nil
}

func (fr *flowRegulator) updateDynamicConfig(dynamicConfig *policylangv1.FlowRegulator_DynamicConfig) {
	logger := fr.registry.GetLogger()

	if dynamicConfig == nil {
		return
	}
	labelValues := make(map[string]bool)
	// loop through overrides
	for _, labelValue := range dynamicConfig.EnableLabelValues {
		labelValues[labelValue] = true
	}

	logger.Debug().Interface("enable values", labelValues).Str("name", fr.name).Msgf("Updating dynamic config for flow regulator")

	fr.setEnableValues(labelValues)
}

func (fr *flowRegulator) setEnableValues(labelValues map[string]bool) {
	fr.enableValuesMutex.Lock()
	defer fr.enableValuesMutex.Unlock()
	fr.enableLabelValues = labelValues
}

// GetFlowSelector returns the selector for the rate limiter.
func (fr *flowRegulator) GetFlowSelector() *policylangv1.FlowSelector {
	return fr.proto.Parameters.GetFlowSelector()
}

// Decide runs the limiter.
func (fr *flowRegulator) Decide(ctx context.Context,
	labels map[string]string,
) *flowcontrolv1.LimiterDecision {
	var (
		labelValue  string
		hasLabelKey bool
	)
	labelKey := fr.proto.GetParameters().GetLabelKey()
	if labelKey != "" {
		labelValue, hasLabelKey = labels[fr.proto.GetParameters().GetLabelKey()]
	}

	// Initialize LimiterDecision
	limiterDecision := &flowcontrolv1.LimiterDecision{
		PolicyName:  fr.GetPolicyName(),
		PolicyHash:  fr.GetPolicyHash(),
		ComponentId: fr.GetComponentId(),
		Dropped:     false,
		Reason:      flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED,
		Details: &flowcontrolv1.LimiterDecision_FlowRegulatorInfo_{
			FlowRegulatorInfo: &flowcontrolv1.LimiterDecision_FlowRegulatorInfo{
				Label: labelKey + ":" + labelValue,
			},
		},
	}

	// If label_key is a non-empty string and is found within labels
	if labelKey != "" && hasLabelKey {
		hashValue := fnv.New32a()
		hashValue.Write([]byte(labelValue))
		hash := hashValue.Sum32()

		// Allow only acceptPercentage proportion of requests
		if float64(hash)/float64(math.MaxUint32) <= fr.acceptPercentage/100.0 {
			limiterDecision.Dropped = false
		} else {
			limiterDecision.Dropped = true
		}
	} else {
		// Else, label_key is empty or not found in labels
		// Randomly accept only acceptPercentage proportion of requests
		// #nosec G404
		// G404: Use of weak random number generator (math/rand instead of crypto/rand) (gosec)
		// This is not a security issue as we do not need cryptographic randomness for load management.
		val := rand.Float64()

		if val <= fr.acceptPercentage/100.0 {
			limiterDecision.Dropped = false
		} else {
			limiterDecision.Dropped = true
		}
	}

	return limiterDecision
}

// Revert implements the Revert method of the flowcontrolv1.FlowRegulator interface.
func (fr *flowRegulator) Revert(_ map[string]string, _ *flowcontrolv1.LimiterDecision) {
	// No-op
}

func (fr *flowRegulator) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := fr.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		fr.acceptPercentage = 100
		return
	}

	var wrapperMessage policysyncv1.FlowRegulatorDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.FlowRegulatorDecision == nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Common attributes not found")
		return
	}
	if commonAttributes.PolicyHash != fr.GetPolicyHash() {
		return
	}
	flowRegulatorDecision := wrapperMessage.FlowRegulatorDecision
	fr.acceptPercentage = flowRegulatorDecision.AcceptPercentage
}

func (fr *flowRegulator) dynamicConfigUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := fr.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Dynamic config removed")
		fr.enableLabelValues = make(map[string]bool)
		return
	}

	var wrapperMessage policysyncv1.FlowRegulatorDynamicConfigWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.FlowRegulatorDynamicConfig == nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Common attributes not found")
		return
	}
	if commonAttributes.PolicyHash != fr.GetPolicyHash() {
		return
	}
	dynamicConfig := wrapperMessage.FlowRegulatorDynamicConfig
	fr.updateDynamicConfig(dynamicConfig)
}

// GetLimiterID returns the limiter ID.
func (fr *flowRegulator) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  fr.GetPolicyName(),
		PolicyHash:  fr.GetPolicyHash(),
		ComponentID: fr.GetComponentId(),
	}
}

// GetRequestCounter returns counter for tracking number of times flowRegulator was triggered.
func (fr *flowRegulator) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := fr.factory.counterVector.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get counter")
		return nil
	}

	return counter
}
