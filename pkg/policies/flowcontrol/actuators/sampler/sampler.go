package sampler

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
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

const samplerStatusRoot = "samplers"

var (
	fxNameTag       = config.NameTag(samplerStatusRoot)
	metricLabelKeys = []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.DecisionTypeLabel,
		metrics.SamplerDroppedLabel,
	}
)

// Module returns the fx options for sampler.
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
				setupSamplerFactory,
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

	etcdPath := path.Join(paths.SamplerConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type samplerFactory struct {
	engineAPI        iface.Engine
	registry         status.Registry
	distCache        *distcache.DistCache
	decisionsWatcher notifiers.Watcher
	counterVector    *prometheus.CounterVec
	agentGroupName   string
}

func setupSamplerFactory(
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
	etcdPath := path.Join(paths.SamplerDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", samplerStatusRoot)
	counterVector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.SamplerCounterTotalMetricName,
			Help: "Total number of decisions made by samplers.",
		},
		metricLabelKeys,
	)

	factory := &samplerFactory{
		engineAPI:        e,
		registry:         statusRegistry,
		distCache:        distCache,
		decisionsWatcher: decisionsWatcher,
		counterVector:    counterVector,
		agentGroupName:   agentGroupName,
	}

	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{factory.newSamplerOptions})
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
			return nil
		},
		OnStop: func(context.Context) error {
			var err, merr error
			err = decisionsWatcher.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(factory.counterVector) {
				err2 := fmt.Errorf("failed to unregister sampler_counter_total metric")
				merr = multierr.Append(merr, err2)
			}
			reg.Detach()
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

func (frf *samplerFactory) newSamplerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := frf.registry.GetLogger()
	wrapperMessage := &policysyncv1.SamplerWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.Sampler == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal sampler config")
		return fx.Options(), err
	}

	samplerProto := wrapperMessage.Sampler
	fr := &sampler{
		Component:              wrapperMessage.GetCommonAttributes(),
		proto:                  samplerProto,
		labelKey:               samplerProto.GetParameters().GetLabelKey(),
		factory:                frf,
		registry:               reg,
		passthroughLabelValues: make(map[string]bool),
	}
	fr.name = iface.ComponentKey(fr)

	return fx.Options(
		fx.Invoke(
			fr.setup,
		),
	), nil
}

type sampler struct {
	iface.Component
	registry                    status.Registry
	factory                     *samplerFactory
	proto                       *policylangv1.Sampler
	passthroughLabelValues      map[string]bool
	name                        string
	labelKey                    string
	acceptPercentage            float64
	passthroughLabelValuesMutex sync.RWMutex
}

// Make sure sampler implements iface.Limiter.
var _ iface.Limiter = (*sampler)(nil)

func (fr *sampler) setup(lifecycle fx.Lifecycle) error {
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

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = fr.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = fr.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = fr.GetComponentId()
	counterVec := fr.factory.counterVector

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			// add decisions notifier
			err = fr.factory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}

			// add to data engine
			err = fr.factory.engineAPI.RegisterSampler(fr)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to register sampler")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error
			deleted := counterVec.DeletePartialMatch(metricLabels)
			if deleted == 0 {
				logger.Warn().Msg("Could not delete sampler counter from its metric vector. No traffic to generate metrics?")
			}
			// remove from data engine
			err = fr.factory.engineAPI.UnregisterSampler(fr)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to unregister rate limiter")
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

func (fr *sampler) setPassthroughLabelValues(labelList []string) {
	labelSet := make(map[string]bool)
	for _, labelValue := range labelList {
		labelSet[labelValue] = true
	}
	fr.passthroughLabelValuesMutex.Lock()
	defer fr.passthroughLabelValuesMutex.Unlock()
	fr.passthroughLabelValues = labelSet
}

// GetSelectors returns the selectors for the sampler.
func (fr *sampler) GetSelectors() []*policylangv1.Selector {
	return fr.proto.Parameters.GetSelectors()
}

// Decide runs the limiter.
func (fr *sampler) Decide(ctx context.Context,
	labels labels.Labels,
) *flowcontrolv1.LimiterDecision {
	var (
		labelValue  string
		hasLabelKey bool
	)
	labelKey := fr.proto.GetParameters().GetLabelKey()
	if labelKey != "" {
		labelValue, hasLabelKey = labels.Get(fr.proto.GetParameters().GetLabelKey())
	}

	// Initialize LimiterDecision
	limiterDecision := &flowcontrolv1.LimiterDecision{
		PolicyName:               fr.GetPolicyName(),
		PolicyHash:               fr.GetPolicyHash(),
		ComponentId:              fr.GetComponentId(),
		Dropped:                  false,
		DeniedResponseStatusCode: fr.proto.GetParameters().GetDeniedResponseStatusCode(),
		Reason:                   flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED,
		Details: &flowcontrolv1.LimiterDecision_SamplerInfo_{
			SamplerInfo: &flowcontrolv1.LimiterDecision_SamplerInfo{
				Label: labelKey + ":" + labelValue,
			},
		},
	}

	// Check if labelValue is in passthroughLabelValues
	fr.passthroughLabelValuesMutex.RLock()
	_, ok := fr.passthroughLabelValues[labelValue]
	fr.passthroughLabelValuesMutex.RUnlock()
	if ok {
		limiterDecision.Dropped = false
		return limiterDecision
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

// Revert implements the Revert method of the flowcontrolv1.Sampler interface.
func (fr *sampler) Revert(_ context.Context, _ labels.Labels, _ *flowcontrolv1.LimiterDecision) {
	// No-op
}

func (fr *sampler) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := fr.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		if fr.proto.Parameters.RampMode {
			fr.acceptPercentage = 0
		} else {
			fr.acceptPercentage = 100
		}
		return
	}

	var wrapperMessage policysyncv1.SamplerDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.SamplerDecision == nil {
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
	samplerDecision := wrapperMessage.SamplerDecision
	fr.setPassthroughLabelValues(samplerDecision.PassThroughLabelValues)
	fr.acceptPercentage = samplerDecision.AcceptPercentage
}

// GetLimiterID returns the limiter ID.
func (fr *sampler) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  fr.GetPolicyName(),
		PolicyHash:  fr.GetPolicyHash(),
		ComponentID: fr.GetComponentId(),
	}
}

// GetRequestCounter returns counter for tracking number of times sampler was triggered.
func (fr *sampler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := fr.factory.counterVector.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get counter")
		return nil
	}

	return counter
}

// GetRampMode returns the ramp mode flag of the sampler.
func (fr *sampler) GetRampMode() bool {
	return fr.proto.Parameters.RampMode
}
