package concurrency

import (
	"context"
	"fmt"
	"math"
	"path"
	"strconv"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators/concurrency/scheduler"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/selectors"
	"github.com/fluxninja/aperture/pkg/status"
)

var (
	// FxNameTag is Concurrency Limiter Watcher's Fx Tag.
	fxNameTag = config.NameTag("concurrency_limiter_watcher")

	// Array of Label Keys for WFQ and Token Bucket Metrics.
	metricLabelKeys = []string{metrics.PolicyNameLabel, metrics.PolicyHashLabel, metrics.ComponentIndexLabel}
)

// concurrencyLimiterModule returns the fx options for dataplane side pieces of concurrency limiter in the main fx app.
func concurrencyLimiterModule() fx.Option {
	return fx.Options(
		// Tag the watcher so that other modules can find it.
		fx.Provide(
			fx.Annotate(
				provideWatcher,
				fx.ResultTags(fxNameTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupConcurrencyLimiterFactory,
				fx.ParamTags(fxNameTag),
			),
		),
	)
}

// provideWatcher provides pointer to concurrency limiter watcher.
func provideWatcher(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	// Get Agent Group from host info gatherer
	agentGroupName := ai.GetAgentGroup()
	// Scope the sync to the agent group.
	etcdPath := path.Join(common.ConcurrencyLimiterConfigPath, common.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}

type concurrencyLimiterFactory struct {
	engineAPI  iface.Engine
	metricsAPI iface.ResponseMetricsAPI
	registry   status.Registry

	autoTokensFactory       *autoTokensFactory
	loadShedActuatorFactory *loadShedActuatorFactory

	// WFQ Metrics.
	wfqFlowsGaugeVec    *prometheus.GaugeVec
	wfqRequestsGaugeVec *prometheus.GaugeVec

	// TODO: following will be moved to scheduler.
	incomingConcurrencyCounterVec *prometheus.CounterVec
	acceptedConcurrencyCounterVec *prometheus.CounterVec
}

// setupConcurrencyLimiterFactory sets up the concurrency limiter module in the main fx app.
func setupConcurrencyLimiterFactory(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
	m iface.ResponseMetricsAPI,
) error {
	agentGroup := ai.GetAgentGroup()

	// Create factories
	loadShedActuatorFactory, err := newLoadShedActuatorFactory(lifecycle, etcdClient, agentGroup, prometheusRegistry)
	if err != nil {
		return err
	}

	autoTokensFactory, err := newAutoTokensFactory(lifecycle, etcdClient, agentGroup)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("concurrency_limiter")

	conLimiterFactory := &concurrencyLimiterFactory{
		engineAPI:               e,
		metricsAPI:              m,
		autoTokensFactory:       autoTokensFactory,
		loadShedActuatorFactory: loadShedActuatorFactory,
		registry:                reg,
	}

	conLimiterFactory.wfqFlowsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.WFQFlowsMetricName,
			Help: "A gauge that tracks the number of flows in the WFQScheduler",
		},
		metricLabelKeys,
	)
	conLimiterFactory.wfqRequestsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.WFQRequestsMetricName,
			Help: "A gauge that tracks the number of queued requests in the WFQScheduler",
		},
		metricLabelKeys,
	)
	conLimiterFactory.incomingConcurrencyCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.IncomingConcurrencyMetricName,
			Help: "A counter measuring incoming concurrency into Scheduler",
		},
		metricLabelKeys,
	)
	conLimiterFactory.acceptedConcurrencyCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.AcceptedConcurrencyMetricName,
			Help: "A counter measuring the concurrency admitted by Scheduler",
		},
		metricLabelKeys,
	)

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: []notifiers.FxOptionsFunc{conLimiterFactory.newConcurrencyLimiterOptions},
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.NewProtobufUnmarshaller,
		},
		StatusRegistry:     reg,
		PrometheusRegistry: prometheusRegistry,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var merr error

			err := prometheusRegistry.Register(conLimiterFactory.wfqFlowsGaugeVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(conLimiterFactory.wfqRequestsGaugeVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(conLimiterFactory.incomingConcurrencyCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(conLimiterFactory.acceptedConcurrencyCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}

			return merr
		},
		OnStop: func(_ context.Context) error {
			var merr error

			if !prometheusRegistry.Unregister(conLimiterFactory.wfqFlowsGaugeVec) {
				err := fmt.Errorf("failed to unregister wfq_flows metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(conLimiterFactory.wfqRequestsGaugeVec) {
				err := fmt.Errorf("failed to unregister wfq_requests metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(conLimiterFactory.incomingConcurrencyCounterVec) {
				err := fmt.Errorf("failed to unregister incoming_concurrency metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(conLimiterFactory.acceptedConcurrencyCounterVec) {
				err := fmt.Errorf("failed to unregister accepted_concurrency metric")
				merr = multierr.Append(merr, err)
			}

			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	matchedWorkloads map[int]*policylangv1.Scheduler_WorkloadParameters
}

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[int, multiMatchResult]

// newConcurrencyLimiterOptions returns fx options for the concurrency limiter fx app.
func (conLimiterFactory *concurrencyLimiterFactory) newConcurrencyLimiterOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := conLimiterFactory.registry.GetLogger()
	wrapperMessage := &wrappersv1.ConcurrencyLimiterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	concurrencyLimiterMessage := wrapperMessage.ConcurrencyLimiter
	if err != nil || concurrencyLimiterMessage == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal concurrency limiter config wrapper")
		return fx.Options(), err
	}

	// Loop through the workloads
	schedulerProto := concurrencyLimiterMessage.Scheduler
	if schedulerProto == nil {
		err = fmt.Errorf("no scheduler specified")
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal scheduler")
		return fx.Options(), err
	}
	mm := multimatcher.New[int, multiMatchResult]()
	for workloadIndex, workloadProto := range schedulerProto.Workloads {
		labelMatcher, err := selectors.MMExprFromLabelMatcher(workloadProto.GetLabelMatcher())
		if err != nil {
			return fx.Options(), err
		}
		wm := &workloadMatcher{
			workloadIndex: workloadIndex,
			workloadProto: workloadProto,
		}
		err = mm.AddEntry(workloadIndex, labelMatcher, wm.matchCallback)
		if err != nil {
			return fx.Options(), err
		}
	}

	conLimiter := &concurrencyLimiter{
		Component:                      wrapperMessage,
		concurrencyLimiterProto:        concurrencyLimiterMessage,
		registry:                       reg,
		concurrencyLimiterFactory:      conLimiterFactory,
		workloadMultiMatcher:           mm,
		defaultWorkloadParametersProto: schedulerProto.DefaultWorkloadParameters,
		schedulerProto:                 schedulerProto,
	}

	return fx.Options(
		fx.Invoke(
			conLimiter.setup,
		),
	), nil
}

type workloadMatcher struct {
	workloadProto *policylangv1.Scheduler_Workload
	workloadIndex int
}

func (wm *workloadMatcher) matchCallback(mmr multiMatchResult) multiMatchResult {
	// mmr.matchedWorkloads is nil on first match.
	if mmr.matchedWorkloads == nil {
		mmr.matchedWorkloads = make(map[int]*policylangv1.Scheduler_WorkloadParameters)
	}
	mmr.matchedWorkloads[wm.workloadIndex] = wm.workloadProto.GetWorkloadParameters()
	return mmr
}

// concurrencyLimiter implements concurrency limiter on the dataplane side.
type concurrencyLimiter struct {
	iface.Component
	scheduler                      scheduler.Scheduler
	registry                       status.Registry
	incomingConcurrencyCounter     prometheus.Counter
	acceptedConcurrencyCounter     prometheus.Counter
	concurrencyLimiterProto        *policylangv1.ConcurrencyLimiter
	concurrencyLimiterFactory      *concurrencyLimiterFactory
	autoTokens                     *autoTokens
	workloadMultiMatcher           *multiMatcher
	defaultWorkloadParametersProto *policylangv1.Scheduler_WorkloadParameters
	schedulerProto                 *policylangv1.Scheduler
}

// Make sure ConcurrencyLimiter implements the iface.ConcurrencyLimiter.
var _ iface.Limiter = &concurrencyLimiter{}

func (conLimiter *concurrencyLimiter) setup(lifecycle fx.Lifecycle) error {
	// Factories
	conLimiterFactory := conLimiter.concurrencyLimiterFactory
	loadShedActuatorFactory := conLimiterFactory.loadShedActuatorFactory
	autoTokensFactory := conLimiterFactory.autoTokensFactory
	// Form metric labels
	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = conLimiter.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = conLimiter.GetPolicyHash()
	metricLabels[metrics.ComponentIndexLabel] = strconv.FormatInt(conLimiter.GetComponentIndex(), 10)
	// Create sub components.
	clock := clockwork.NewRealClock()
	loadShedActuator, err := loadShedActuatorFactory.newLoadShedActuator(conLimiter, conLimiter.registry, clock, lifecycle, metricLabels)
	if err != nil {
		return err
	}
	if conLimiter.schedulerProto.AutoTokens {
		autoTokens, err := autoTokensFactory.newAutoTokens(
			conLimiter.GetPolicyName(), conLimiter.GetPolicyHash(),
			lifecycle, conLimiter.GetComponentIndex(), conLimiter.registry)
		if err != nil {
			return err
		}
		conLimiter.autoTokens = autoTokens
	}

	engineAPI := conLimiterFactory.engineAPI
	metricsAPI := conLimiterFactory.metricsAPI
	wfqFlowsGaugeVec := conLimiterFactory.wfqFlowsGaugeVec
	wfqRequestsGaugeVec := conLimiterFactory.wfqRequestsGaugeVec
	incomingConcurrencyCounterVec := conLimiterFactory.incomingConcurrencyCounterVec
	acceptedConcurrencyCounterVec := conLimiterFactory.acceptedConcurrencyCounterVec

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			retErr := func(err error) error {
				conLimiter.registry.SetStatus(status.NewStatus(nil, err))
				return err
			}

			wfqFlowsGauge, err := wfqFlowsGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return errors.Wrap(err, "failed to get wfq flows gauge")
			}

			wfqRequestsGauge, err := wfqRequestsGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return errors.Wrap(err, "failed to get wfq requests gauge")
			}

			wfqMetrics := &scheduler.WFQMetrics{
				FlowsGauge:        wfqFlowsGauge,
				HeapRequestsGauge: wfqRequestsGauge,
			}

			// setup scheduler
			conLimiter.scheduler = scheduler.NewWFQScheduler(conLimiter.schedulerProto.MaxTimeout.AsDuration(), loadShedActuator.tokenBucketLoadShed, clock, wfqMetrics)

			incomingConcurrencyCounter, err := incomingConcurrencyCounterVec.GetMetricWith(metricLabels)
			if err != nil {
				return err
			}
			conLimiter.incomingConcurrencyCounter = incomingConcurrencyCounter
			acceptedConcurrencyCounter, err := acceptedConcurrencyCounterVec.GetMetricWith(metricLabels)
			if err != nil {
				return err
			}
			conLimiter.acceptedConcurrencyCounter = acceptedConcurrencyCounter

			err = engineAPI.RegisterConcurrencyLimiter(conLimiter)
			if err != nil {
				return retErr(err)
			}

			return retErr(nil)
		},
		OnStop: func(context.Context) error {
			var errMulti error

			err := engineAPI.UnregisterConcurrencyLimiter(conLimiter)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			metricsAPI.DeleteTokenLatencyHistogram(metricLabels)

			// Remove metrics from metric vectors
			deleted := wfqFlowsGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete wfq_flows gauge from its metric vector"))
			}
			deleted = wfqRequestsGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete wfq_requests gauge from its metric vector"))
			}
			deleted = incomingConcurrencyCounterVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete incoming_concurrency counter from its metric vector"))
			}
			deleted = acceptedConcurrencyCounterVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete accepted_concurrency counter from its metric vector"))
			}

			conLimiter.registry.SetStatus(status.NewStatus(nil, errMulti))
			return errMulti
		},
	})

	return nil
}

// GetSelector returns selector.
func (conLimiter *concurrencyLimiter) GetSelector() *selectorv1.Selector {
	return conLimiter.schedulerProto.GetSelector()
}

// RunLimiter .
func (conLimiter *concurrencyLimiter) RunLimiter(labels map[string]string) *flowcontrolv1.LimiterDecision {
	var matchedWorkloadProto *policylangv1.Scheduler_WorkloadParameters
	var matchedWorkloadIndex string
	// match labels against conLimiter.workloadMultiMatcher
	mmr := conLimiter.workloadMultiMatcher.Match(multimatcher.Labels(labels))
	// if at least one match, return workload with lowest index
	if len(mmr.matchedWorkloads) > 0 {
		// select the smallest workloadIndex
		smallestWorkloadIndex := math.MaxInt32
		for workloadIndex := range mmr.matchedWorkloads {
			if workloadIndex < smallestWorkloadIndex {
				smallestWorkloadIndex = workloadIndex
			}
		}
		matchedWorkloadProto = mmr.matchedWorkloads[smallestWorkloadIndex]
		matchedWorkloadIndex = strconv.Itoa(smallestWorkloadIndex)
	} else {
		// no match, return default workload
		matchedWorkloadProto = conLimiter.defaultWorkloadParametersProto
		matchedWorkloadIndex = metrics.DefaultWorkloadIndex
	}

	fairnessLabel := "workload:" + matchedWorkloadIndex

	if val, ok := labels[matchedWorkloadProto.FairnessKey]; ok {
		fairnessLabel = fairnessLabel + "," + val
	}
	// Lookup tokens for the workload
	var tokens uint64
	if conLimiter.schedulerProto.AutoTokens {
		tokensAuto, ok := conLimiter.autoTokens.GetTokensForWorkload(matchedWorkloadIndex)
		if !ok {
			// default to 1 if auto tokens not found
			tokens = 1
		} else {
			tokens = tokensAuto
		}
	} else {
		tokens = matchedWorkloadProto.Tokens
	}

	// timeout is tokens(which is in milliseconds) * conLimiter.schedulerProto.TimeoutFactor(float64)
	timeout := time.Duration(float64(tokens)*conLimiter.schedulerProto.TimeoutFactor) * time.Millisecond

	if timeout > conLimiter.schedulerProto.MaxTimeout.AsDuration() {
		timeout = conLimiter.schedulerProto.MaxTimeout.AsDuration()
	}

	reqContext := scheduler.RequestContext{
		FairnessLabel: fairnessLabel,
		Priority:      uint8(matchedWorkloadProto.Priority),
		Timeout:       timeout,
		Tokens:        tokens,
	}

	accepted := conLimiter.scheduler.Schedule(reqContext)

	// update concurrency metrics and decisionType
	conLimiter.incomingConcurrencyCounter.Add(float64(reqContext.Tokens))

	if accepted {
		conLimiter.acceptedConcurrencyCounter.Add(float64(reqContext.Tokens))
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:     conLimiter.GetPolicyName(),
		PolicyHash:     conLimiter.GetPolicyHash(),
		ComponentIndex: conLimiter.GetComponentIndex(),
		Dropped:        !accepted,
		Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo_{
			ConcurrencyLimiterInfo: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo{
				WorkloadIndex: matchedWorkloadIndex,
			},
		},
	}
}

// GetLimiterID returns the limiter ID.
func (conLimiter *concurrencyLimiter) GetLimiterID() iface.LimiterID {
	// TODO: move this to limiter base.
	return iface.LimiterID{
		PolicyName:     conLimiter.GetPolicyName(),
		ComponentIndex: conLimiter.GetComponentIndex(),
		PolicyHash:     conLimiter.GetPolicyHash(),
	}
}
