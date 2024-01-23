package concurrencylimiter

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/types/known/durationpb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	concurrencylimiter "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/concurrency-limiter"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/labelstatus"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

const concurrencyLimiterStatusRoot = "concurrency_limiters"

var (
	fxTag           = config.NameTag(concurrencyLimiterStatusRoot)
	metricLabelKeys = []string{metrics.PolicyNameLabel, metrics.PolicyHashLabel, metrics.ComponentIDLabel, metrics.DecisionTypeLabel, metrics.LimiterDroppedLabel}
)

func concurrencyLimiterModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideConcurrencyLimiterWatchers,
				fx.ResultTags(fxTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupConcurrencyLimiterFactory,
				fx.ParamTags(fxTag),
			),
		),
	)
}

func provideConcurrencyLimiterWatchers(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	agentGroupName := ai.GetAgentGroup()

	etcdPath := path.Join(paths.ConcurrencyLimiterConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type concurrencyLimiterFactory struct {
	engineAPI           iface.Engine
	registry            status.Registry
	distCache           *distcache.DistCache
	decisionsWatcher    notifiers.Watcher
	counterVector       *prometheus.CounterVec
	agentGroupName      string
	labelStatusJobGroup *jobs.JobGroup
	labelStatusFactory  *labelstatus.LabelStatusFactory
}

// main fx app.
func setupConcurrencyLimiterFactory(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	distCache *distcache.DistCache,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
	labelStatusFactory *labelstatus.LabelStatusFactory,
) error {
	agentGroupName := ai.GetAgentGroup()
	etcdPath := path.Join(paths.ConcurrencyLimiterDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", concurrencyLimiterStatusRoot)
	logger := reg.GetLogger()

	counterVector := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.ConcurrencyLimiterCounterTotalMetricName,
		Help: "A counter measuring the number of times Concurrency Limiter was triggered",
	}, metricLabelKeys)

	labelStatusJobGroup, err := jobs.NewJobGroup(reg.Child("label_status", concurrencyLimiterStatusRoot), jobs.JobGroupConfig{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create labels status job group")
		return err
	}

	concurrencyLimiterFactory := &concurrencyLimiterFactory{
		engineAPI:           e,
		distCache:           distCache,
		decisionsWatcher:    decisionsWatcher,
		agentGroupName:      agentGroupName,
		registry:            reg,
		counterVector:       counterVector,
		labelStatusJobGroup: labelStatusJobGroup,
		labelStatusFactory:  labelStatusFactory,
	}

	fxDriver, err := notifiers.NewFxDriver(
		reg,
		prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{concurrencyLimiterFactory.newConcurrencyLimiterOptions},
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := prometheusRegistry.Register(concurrencyLimiterFactory.counterVector)
			if err != nil {
				return err
			}
			err = labelStatusJobGroup.Start()
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
			err = labelStatusJobGroup.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(concurrencyLimiterFactory.counterVector) {
				err2 := fmt.Errorf("failed to unregister metric")
				merr = multierr.Append(merr, err2)
			}
			reg.Detach()
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// per component fx app.
func (clFactory *concurrencyLimiterFactory) newConcurrencyLimiterOptions(key notifiers.Key, unmarshaller config.Unmarshaller, reg status.Registry) (fx.Option, error) {
	logger := clFactory.registry.GetLogger()

	wrapperMessage := &policysyncv1.ConcurrencyLimiterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.ConcurrencyLimiter == nil {
		reg.SetStatus(status.NewStatus(nil, err), nil)
		logger.Warn().Err(err).Msg("Failed to unmarshal concurrency limiter config")
		return fx.Options(), err
	}

	clProto := wrapperMessage.ConcurrencyLimiter
	cl := &concurrencyLimiter{
		Component:           wrapperMessage.GetCommonAttributes(),
		clProto:             clProto,
		clFactory:           clFactory,
		registry:            reg,
		labelStatusJobGroup: clFactory.labelStatusJobGroup,
	}
	cl.name = iface.ComponentKey(cl)
	cl.tokensLabelKeyStatus = clFactory.labelStatusFactory.New("tokens_label_key", cl.GetPolicyName(), cl.GetComponentId())
	cl.limitByLabelKeyStatus = clFactory.labelStatusFactory.New("limit_by_label_key", cl.GetPolicyName(), cl.GetComponentId())

	return fx.Options(
		fx.Invoke(
			cl.setup,
		),
	), nil
}

// concurrencyLimiter implements concurrency limiter on the data plane side.
type concurrencyLimiter struct {
	iface.Component
	registry              status.Registry
	clFactory             *concurrencyLimiterFactory
	limiter               concurrencylimiter.ConcurrencyLimiter
	clProto               *policylangv1.ConcurrencyLimiter
	name                  string
	labelStatusJobGroup   *jobs.JobGroup
	tokensLabelKeyStatus  *labelstatus.LabelStatus
	limitByLabelKeyStatus *labelstatus.LabelStatus
}

// Make sure concurrencyLimiter implements iface.Limiter.
var _ iface.Limiter = (*concurrencyLimiter)(nil)

func (cl *concurrencyLimiter) setup(lifecycle fx.Lifecycle) error {
	logger := cl.registry.GetLogger()
	etcdKey := paths.AgentComponentKey(
		cl.clFactory.agentGroupName,
		cl.GetPolicyName(),
		cl.GetComponentId(),
	)
	// decision notifier
	decisionUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaller,
		cl.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = cl.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = cl.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = cl.GetComponentId()

	// setup labels status
	cl.tokensLabelKeyStatus.Setup(cl.clFactory.labelStatusJobGroup, lifecycle)
	cl.limitByLabelKeyStatus.Setup(cl.clFactory.labelStatusJobGroup, lifecycle)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error

			inner, err := concurrencylimiter.NewGlobalTokenCounter(
				cl.clFactory.distCache,
				cl.name,
				cl.clProto.Parameters.GetMaxIdleTime().AsDuration(),
				cl.clProto.Parameters.GetMaxInflightDuration().AsDuration(),
			)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create limiter")
				return err
			}
			cl.limiter = inner

			// add decisions notifier
			err = cl.clFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}

			// add to data engine
			err = cl.clFactory.engineAPI.RegisterConcurrencyLimiter(cl)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to register concurrency limiter")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error
			// remove from data engine
			err = cl.clFactory.engineAPI.UnregisterConcurrencyLimiter(cl)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to unregister concurrency limiter")
				merr = multierr.Append(merr, err)
			}
			// remove decisions notifier
			err = cl.clFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierr.Append(merr, err)
			}
			cl.limiter.Close()
			deleted := cl.clFactory.counterVector.DeletePartialMatch(metricLabels)
			if deleted == 0 {
				logger.Warn().Msg("Could not delete concurrency limiter counter from its metric vector. No traffic to generate metrics?")
			}
			cl.registry.SetStatus(status.NewStatus(nil, merr), nil)

			return merr
		},
	})
	return nil
}

// GetSelectors returns the selectors for the concurrency limiter.
func (cl *concurrencyLimiter) GetSelectors() []*policylangv1.Selector {
	return cl.clProto.GetSelectors()
}

// Decide runs the limiter.
func (cl *concurrencyLimiter) Decide(ctx context.Context, labels labels.Labels) *flowcontrolv1.LimiterDecision {
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED

	tokens := float64(1)
	// get tokens from labels
	rParams := cl.clProto.GetRequestParameters()
	var deniedResponseStatusCode flowcontrolv1.StatusCode
	if rParams != nil {
		deniedResponseStatusCode = rParams.GetDeniedResponseStatusCode()
		tokensLabelKey := rParams.GetTokensLabelKey()
		if tokensLabelKey != "" {
			val, ok := labels.Get(tokensLabelKey)
			if ok {
				if parsedTokens, err := strconv.ParseFloat(val, 64); err == nil {
					tokens = parsedTokens
				}
			} else {
				cl.tokensLabelKeyStatus.SetMissing()
			}
		}
	}

	label, ok, waitTime, remaining, current, reqID := cl.takeIfAvailable(ctx, labels, tokens)

	tokensConsumed := float64(0)
	if ok {
		tokensConsumed = tokens
	}

	if label == "" {
		reason = flowcontrolv1.LimiterDecision_LIMITER_REASON_KEY_NOT_FOUND
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:               cl.GetPolicyName(),
		PolicyHash:               cl.GetPolicyHash(),
		ComponentId:              cl.GetComponentId(),
		Dropped:                  !ok,
		DeniedResponseStatusCode: deniedResponseStatusCode,
		Reason:                   reason,
		WaitTime:                 durationpb.New(waitTime),
		Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo_{
			ConcurrencyLimiterInfo: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo{
				Label: label,
				TokensInfo: &flowcontrolv1.LimiterDecision_TokensInfo{
					Remaining: remaining,
					Current:   current,
					Consumed:  tokensConsumed,
				},
				RequestId: reqID,
			},
		},
	}
}

// Revert returns the tokens to the limiter.
func (cl *concurrencyLimiter) Revert(ctx context.Context, labels labels.Labels, decision *flowcontrolv1.LimiterDecision) {
	if concurrencyLimiterDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo_); ok {
		tokens := concurrencyLimiterDecision.ConcurrencyLimiterInfo.TokensInfo.Consumed
		if tokens > 0 {
			_, err := cl.limiter.Return(ctx, concurrencyLimiterDecision.ConcurrencyLimiterInfo.Label, tokens,
				concurrencyLimiterDecision.ConcurrencyLimiterInfo.RequestId)
			if err != nil {
				log.Autosample().Error().Err(err).Msg("Failed to return tokens")
			}
		}
	}
}

// Return returns the tokens to the limiter.
func (cl *concurrencyLimiter) Return(ctx context.Context, label string, tokens float64, requestID string) (bool, error) {
	return cl.limiter.Return(ctx, label, tokens, requestID)
}

// takeIfAvailable takes n tokens from the limiter.
func (cl *concurrencyLimiter) takeIfAvailable(
	ctx context.Context,
	labels labels.Labels,
	n float64,
) (label string, ok bool, waitTime time.Duration, remaining float64, current float64, reqID string) {
	if cl.limiter.GetPassThrough() {
		return label, true, 0, 0, 0, ""
	}

	label, found := cl.getLimitLabelFromLabels(labels)
	if !found {
		return label, true, 0, 0, 0, ""
	}

	ok, waitTime, remaining, current, reqID = cl.limiter.TakeIfAvailable(ctx, label, n)
	return
}

func (cl *concurrencyLimiter) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := cl.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		cl.limiter.SetPassThrough(true)
		return
	}

	var wrapperMessage policysyncv1.ConcurrencyLimiterDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.ConcurrencyLimiterDecision == nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Common attributes not found")
		return
	}
	if commonAttributes.PolicyHash != cl.GetPolicyHash() {
		return
	}
	limitDecision := wrapperMessage.ConcurrencyLimiterDecision
	cl.limiter.SetCapacity(limitDecision.MaxConcurrency)
	cl.limiter.SetPassThrough(limitDecision.PassThrough)
}

// GetLimiterID returns the limiter ID.
func (cl *concurrencyLimiter) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  cl.GetPolicyName(),
		PolicyHash:  cl.GetPolicyHash(),
		ComponentID: cl.GetComponentId(),
	}
}

// GetRequestCounter returns counter for tracking number of times concurrencyLimiter was triggered.
func (cl *concurrencyLimiter) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := cl.clFactory.counterVector.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get counter")
		return nil
	}
	return counter
}

// GetRampMode is always false for concurrencyLimiters.
func (cl *concurrencyLimiter) GetRampMode() bool {
	return false
}

// getLabelFromLabels returns the label value from labels and bool indicating if the label was found.
func (cl *concurrencyLimiter) getLimitLabelFromLabels(labels labels.Labels) (label string, ok bool) {
	labelKey := cl.clProto.Parameters.GetLimitByLabelKey()
	if labelKey == "" {
		label = "default"
		return label, true
	}

	labelValue, found := labels.Get(labelKey)
	if !found {
		cl.limitByLabelKeyStatus.SetMissing()
		return "", false
	}

	label = labelKey + ":" + labelValue
	return label, true
}
