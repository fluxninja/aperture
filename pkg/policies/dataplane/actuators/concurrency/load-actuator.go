package concurrency

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators/concurrency/scheduler"
	"github.com/fluxninja/aperture/pkg/status"
)

type loadActuatorFactory struct {
	loadDecisionWatcher                notifiers.Watcher
	tokenBucketLMGaugeVec              *prometheus.GaugeVec
	tokenBucketFillRateGaugeVec        *prometheus.GaugeVec
	tokenBucketBucketCapacityGaugeVec  *prometheus.GaugeVec
	tokenBucketAvailableTokensGaugeVec *prometheus.GaugeVec
	agentGroupName                     string
}

// newLoadActuatorFactory sets up the load actuator module in the main fx app.
func newLoadActuatorFactory(
	lc fx.Lifecycle,
	etcdClient *etcdclient.Client,
	agentGroup string,
	prometheusRegistry *prometheus.Registry,
) (*loadActuatorFactory, error) {
	// Scope the sync to the agent group.
	etcdPath := path.Join(common.LoadDecisionsPath, common.AgentGroupPrefix(agentGroup))
	loadDecisionWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}
	f := &loadActuatorFactory{
		loadDecisionWatcher: loadDecisionWatcher,
		agentGroupName:      agentGroup,
	}
	// Initialize and register the WFQ and Token Bucket Metric Vectors
	f.tokenBucketLMGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketLMMetricName,
			Help: "A gauge that tracks the load multiplier",
		},
		metricLabelKeys,
	)
	f.tokenBucketFillRateGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketFillRateMetricName,
			Help: "A gauge that tracks the fill rate of token bucket in tokens/sec",
		},
		metricLabelKeys,
	)
	f.tokenBucketBucketCapacityGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketCapacityMetricName,
			Help: "A gauge that tracks the capacity of token bucket",
		},
		metricLabelKeys,
	)
	f.tokenBucketAvailableTokensGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketAvailableMetricName,
			Help: "A gauge that tracks the number of tokens available in token bucket",
		},
		metricLabelKeys,
	)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := prometheusRegistry.Register(f.tokenBucketLMGaugeVec)
			if err != nil {
				return err
			}
			err = prometheusRegistry.Register(f.tokenBucketFillRateGaugeVec)
			if err != nil {
				return err
			}
			err = prometheusRegistry.Register(f.tokenBucketBucketCapacityGaugeVec)
			if err != nil {
				return err
			}
			err = prometheusRegistry.Register(f.tokenBucketAvailableTokensGaugeVec)
			if err != nil {
				return err
			}

			err = loadDecisionWatcher.Start()
			if err != nil {
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			var errMulti error
			err := loadDecisionWatcher.Stop()
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			if !prometheusRegistry.Unregister(f.tokenBucketLMGaugeVec) {
				err := fmt.Errorf("failed to unregister " + metrics.TokenBucketLMMetricName)
				errMulti = multierr.Append(errMulti, err)
			}
			if !prometheusRegistry.Unregister(f.tokenBucketFillRateGaugeVec) {
				err := fmt.Errorf("failed to unregister " + metrics.TokenBucketFillRateMetricName)
				errMulti = multierr.Append(errMulti, err)
			}
			if !prometheusRegistry.Unregister(f.tokenBucketBucketCapacityGaugeVec) {
				err := fmt.Errorf("failed to unregister " + metrics.TokenBucketCapacityMetricName)
				errMulti = multierr.Append(errMulti, err)
			}
			if !prometheusRegistry.Unregister(f.tokenBucketAvailableTokensGaugeVec) {
				err := fmt.Errorf("failed to unregister " + metrics.TokenBucketAvailableMetricName)
				errMulti = multierr.Append(errMulti, err)
			}
			return errMulti
		},
	})
	return f, nil
}

// newLoadActuator creates a new load actuator based on proto spec.
func (lsaFactory *loadActuatorFactory) newLoadActuator(conLimiter *concurrencyLimiter,
	registry status.Registry,
	clock clockwork.Clock,
	lifecycle fx.Lifecycle,
	metricLabels prometheus.Labels,
) (*loadActuator, error) {
	reg := registry.Child("load_actuator")

	lsa := &loadActuator{
		conLimiter:     conLimiter,
		clock:          clock,
		statusRegistry: reg,
	}

	unmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return nil, err
	}

	// decision notifier
	decisionNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(common.DataplaneComponentKey(lsaFactory.agentGroupName, lsa.conLimiter.GetPolicyName(), lsa.conLimiter.GetComponentIndex())),
		unmarshaller,
		lsa.decisionUpdateCallback,
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			retErr := func(err error) error {
				lsa.statusRegistry.SetStatus(status.NewStatus(nil, err))
				return err
			}

			tokenBucketLMGauge, err := lsaFactory.tokenBucketLMGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(errors.Wrap(err, "Failed to get token bucket LM gauge"))
			}

			tokenBucketFillRateGauge, err := lsaFactory.tokenBucketFillRateGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(errors.Wrap(err, "Failed to get token bucket fill rate gauge"))
			}

			tokenBucketBucketCapacityGauge, err := lsaFactory.tokenBucketBucketCapacityGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(errors.Wrap(err, "Failed to get token bucket bucket capacity gauge"))
			}

			tokenBucketAvailableTokensGauge, err := lsaFactory.tokenBucketAvailableTokensGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(errors.Wrap(err, "Failed to get token bucket available tokens gauge"))
			}

			tokenBucketMetrics := &scheduler.TokenBucketLoadMultiplierMetrics{
				LMGauge: tokenBucketLMGauge,
				TokenBucketMetrics: &scheduler.TokenBucketMetrics{
					FillRateGauge:        tokenBucketFillRateGauge,
					BucketCapacityGauge:  tokenBucketBucketCapacityGauge,
					AvailableTokensGauge: tokenBucketAvailableTokensGauge,
				},
			}

			// Initialize the token bucket (non continuous tracking mode)
			lsa.tokenBucketLoadMultiplier = scheduler.NewTokenBucketLoadMultiplier(clock.Now(), 10, time.Second, tokenBucketMetrics)

			err = lsaFactory.loadDecisionWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				return retErr(err)
			}
			return retErr(nil)
		},
		OnStop: func(context.Context) error {
			var errMulti error
			err := lsaFactory.loadDecisionWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			deleted := lsaFactory.tokenBucketLMGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete "+metrics.TokenBucketLMMetricName+" from its metric vector"))
			}
			deleted = lsaFactory.tokenBucketFillRateGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete "+metrics.TokenBucketFillRateMetricName+" gauge from its metric vector"))
			}
			deleted = lsaFactory.tokenBucketBucketCapacityGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete "+metrics.TokenBucketCapacityMetricName+" gauge from its metric vector"))
			}
			deleted = lsaFactory.tokenBucketAvailableTokensGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete "+metrics.TokenBucketAvailableMetricName+" gauge from its metric vector"))
			}

			lsa.statusRegistry.SetStatus(status.NewStatus(nil, errMulti))
			return errMulti
		},
	})
	return lsa, nil
}

// loadActuator saves load decisions received from controller.
type loadActuator struct {
	conLimiter                *concurrencyLimiter
	clock                     clockwork.Clock
	tokenBucketLoadMultiplier *scheduler.TokenBucketLoadMultiplier
	statusRegistry            status.Registry
}

func (lsa *loadActuator) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := lsa.statusRegistry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision was removed")
		return
	}

	var wrapperMessage wrappersv1.LoadDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	loadDecision := wrapperMessage.LoadDecision
	if err != nil || loadDecision == nil {
		statusMsg := "Failed to unmarshal config wrapper"
		logger.Warn().Err(err).Msg(statusMsg)
		lsa.statusRegistry.SetStatus(status.NewStatus(nil, err))
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		statusMsg := "Failed to get common attributes from config wrapperShedFactor"
		logger.Error().Err(err).Msg(statusMsg)
		lsa.statusRegistry.SetStatus(status.NewStatus(nil, err))
		return
	}
	// check if this decision is for the same policy id as what we have
	if commonAttributes.PolicyHash != lsa.conLimiter.GetPolicyHash() {
		err = errors.New("policy id mismatch")
		statusMsg := fmt.Sprintf("Expected policy hash: %s, Got: %s", lsa.conLimiter.GetPolicyHash(), commonAttributes.PolicyHash)
		logger.Warn().Err(err).Msg(statusMsg)
		lsa.statusRegistry.SetStatus(status.NewStatus(nil, err))
		return
	}

	logger.Trace().Float64("loadMultiplier", loadDecision.LoadMultiplier).Msg("Setting load multiplier")
	lsa.tokenBucketLoadMultiplier.SetLoadMultiplier(lsa.clock.Now(), loadDecision.LoadMultiplier)
}
