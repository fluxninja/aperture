package regulator

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

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/agentinfo"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/distcache"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

const regulatorStatusRoot = "load_regulators"

var (
	fxNameTag       = config.NameTag(regulatorStatusRoot)
	metricLabelKeys = []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.DecisionTypeLabel,
		metrics.RegulatorDroppedLabel,
	}
)

// Module returns the fx options for load regulator.
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
				setupRegulatorFactory,
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

	etcdPath := path.Join(paths.RegulatorConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type regulatorFactory struct {
	engineAPI            iface.Engine
	registry             status.Registry
	distCache            *distcache.DistCache
	decisionsWatcher     notifiers.Watcher
	dynamicConfigWatcher notifiers.Watcher
	counterVector        *prometheus.CounterVec
	agentGroupName       string
}

func setupRegulatorFactory(
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
	etcdPath := path.Join(paths.RegulatorDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	dynamicConfigWatcher, err := etcdwatcher.NewWatcher(etcdClient,
		paths.RegulatorDynamicConfigPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", regulatorStatusRoot)
	counterVector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.RegulatorCounterTotalMetricName,
			Help: "Total number of decisions made by load regulators.",
		},
		metricLabelKeys,
	)

	factory := &regulatorFactory{
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
		[]notifiers.FxOptionsFunc{factory.newRegulatorOptions})
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
				err2 := fmt.Errorf("failed to unregister regulator_counter_total metric")
				merr = multierr.Append(merr, err2)
			}
			reg.Detach()
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

func (frf *regulatorFactory) newRegulatorOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := frf.registry.GetLogger()
	wrapperMessage := &policysyncv1.RegulatorWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.Regulator == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal load regulator config")
		return fx.Options(), err
	}

	regulatorProto := wrapperMessage.Regulator
	fr := &regulator{
		Component:         wrapperMessage.GetCommonAttributes(),
		proto:             regulatorProto,
		labelKey:          regulatorProto.GetParameters().GetLabelKey(),
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

type regulator struct {
	enableValuesMutex sync.RWMutex
	iface.Component
	registry          status.Registry
	factory           *regulatorFactory
	proto             *policylangv1.Regulator
	enableLabelValues map[string]bool
	name              string
	labelKey          string
	acceptPercentage  float64
}

// Make sure regulator implements iface.Limiter.
var _ iface.Limiter = (*regulator)(nil)

func (fr *regulator) setup(lifecycle fx.Lifecycle) error {
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
			err = fr.factory.engineAPI.RegisterRegulator(fr)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to register load regulator")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error
			deleted := counterVec.DeletePartialMatch(metricLabels)
			if deleted == 0 {
				logger.Warn().Msg("Could not delete load regulator counter from its metric vector. No traffic to generate metrics?")
			}
			// remove from data engine
			err = fr.factory.engineAPI.UnregisterRegulator(fr)
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

func (fr *regulator) updateDynamicConfig(dynamicConfig *policylangv1.Regulator_DynamicConfig) {
	logger := fr.registry.GetLogger()

	if dynamicConfig == nil {
		return
	}
	labelValues := make(map[string]bool)
	// loop through overrides
	for _, labelValue := range dynamicConfig.EnableLabelValues {
		labelValues[labelValue] = true
	}

	logger.Debug().Interface("enable values", labelValues).Str("name", fr.name).Msgf("Updating dynamic config for load regulator")

	fr.setEnableValues(labelValues)
}

func (fr *regulator) setEnableValues(labelValues map[string]bool) {
	fr.enableValuesMutex.Lock()
	defer fr.enableValuesMutex.Unlock()
	fr.enableLabelValues = labelValues
}

// GetSelectors returns the selectors for the load regulator.
func (fr *regulator) GetSelectors() []*policylangv1.Selector {
	return fr.proto.Parameters.GetSelectors()
}

// Decide runs the limiter.
func (fr *regulator) Decide(ctx context.Context,
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
		Details: &flowcontrolv1.LimiterDecision_RegulatorInfo_{
			RegulatorInfo: &flowcontrolv1.LimiterDecision_RegulatorInfo{
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

// Revert implements the Revert method of the flowcontrolv1.Regulator interface.
func (fr *regulator) Revert(_ map[string]string, _ *flowcontrolv1.LimiterDecision) {
	// No-op
}

func (fr *regulator) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := fr.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		fr.acceptPercentage = 100
		return
	}

	var wrapperMessage policysyncv1.RegulatorDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.RegulatorDecision == nil {
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
	regulatorDecision := wrapperMessage.RegulatorDecision
	fr.acceptPercentage = regulatorDecision.AcceptPercentage
}

func (fr *regulator) dynamicConfigUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := fr.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Dynamic config removed")
		fr.enableLabelValues = make(map[string]bool)
		return
	}

	var wrapperMessage policysyncv1.RegulatorDynamicConfigWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.RegulatorDynamicConfig == nil {
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
	dynamicConfig := wrapperMessage.RegulatorDynamicConfig
	fr.updateDynamicConfig(dynamicConfig)
}

// GetLimiterID returns the limiter ID.
func (fr *regulator) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  fr.GetPolicyName(),
		PolicyHash:  fr.GetPolicyHash(),
		ComponentID: fr.GetComponentId(),
	}
}

// GetRequestCounter returns counter for tracking number of times regulator was triggered.
func (fr *regulator) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := fr.factory.counterVector.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get counter")
		return nil
	}

	return counter
}
