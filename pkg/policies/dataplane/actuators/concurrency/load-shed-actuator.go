package concurrency

import (
	"context"
	"fmt"
	"path"

	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators/concurrency/scheduler"
	"github.com/fluxninja/aperture/pkg/status"
)

type loadShedActuatorFactory struct {
	loadShedDecisionWatcher            notifiers.Watcher
	tokenBucketLSFGaugeVec             *prometheus.GaugeVec
	tokenBucketFillRateGaugeVec        *prometheus.GaugeVec
	tokenBucketBucketCapacityGaugeVec  *prometheus.GaugeVec
	tokenBucketAvailableTokensGaugeVec *prometheus.GaugeVec
	agentGroupName                     string
}

// newLoadShedActuatorFactory sets up the load shed module in the main fx app.
func newLoadShedActuatorFactory(
	lc fx.Lifecycle,
	etcdClient *etcdclient.Client,
	agentGroup string,
	prometheusRegistry *prometheus.Registry,
) (*loadShedActuatorFactory, error) {
	// Scope the sync to the agent group.
	etcdPath := path.Join(paths.LoadShedDecisionsPath, paths.AgentGroupPrefix(agentGroup))
	loadShedDecisionWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}
	f := &loadShedActuatorFactory{
		loadShedDecisionWatcher: loadShedDecisionWatcher,
		agentGroupName:          agentGroup,
	}
	// Initialize and register the WFQ and Token Bucket Metric Vectors
	f.tokenBucketLSFGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketMetricName,
			Help: "A gauge that tracks the load shed factor",
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
			err := prometheusRegistry.Register(f.tokenBucketLSFGaugeVec)
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

			err = loadShedDecisionWatcher.Start()
			if err != nil {
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			var errMulti error
			err := loadShedDecisionWatcher.Stop()
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			if !prometheusRegistry.Unregister(f.tokenBucketLSFGaugeVec) {
				err := fmt.Errorf("failed to unregister token_bucket_lsf metric")
				errMulti = multierr.Append(errMulti, err)
			}
			if !prometheusRegistry.Unregister(f.tokenBucketFillRateGaugeVec) {
				err := fmt.Errorf("failed to unregister token_bucket_fill_rate metric")
				errMulti = multierr.Append(errMulti, err)
			}
			if !prometheusRegistry.Unregister(f.tokenBucketBucketCapacityGaugeVec) {
				err := fmt.Errorf("failed to unregister token_bucket_capacity metric")
				errMulti = multierr.Append(errMulti, err)
			}
			if !prometheusRegistry.Unregister(f.tokenBucketAvailableTokensGaugeVec) {
				err := fmt.Errorf("failed to unregister token_bucket_available_tokens metric ")
				errMulti = multierr.Append(errMulti, err)
			}
			return errMulti
		},
	})
	return f, nil
}

// newLoadShedActuator creates a new load shed actuator based on proto spec.
func (lsaFactory *loadShedActuatorFactory) newLoadShedActuator(conLimiter *concurrencyLimiter,
	registry status.Registry,
	clock clockwork.Clock,
	lifecycle fx.Lifecycle,
	metricLabels prometheus.Labels,
) (*loadShedActuator, error) {
	reg := registry.Child("load_shed_actuator")

	lsa := &loadShedActuator{
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
		notifiers.Key(paths.DataplaneComponentKey(lsaFactory.agentGroupName, lsa.conLimiter.GetPolicyName(), lsa.conLimiter.GetComponentIndex())),
		unmarshaller,
		lsa.decisionUpdateCallback,
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			retErr := func(err error) error {
				lsa.statusRegistry.SetStatus(status.NewStatus(nil, err))
				return err
			}

			tokenBucketLSFGauge, err := lsaFactory.tokenBucketLSFGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(errors.Wrap(err, "Failed to get token bucket LSF gauge"))
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

			tokenBucketMetrics := &scheduler.TokenBucketLoadShedMetrics{
				LSFGauge: tokenBucketLSFGauge,
				TokenBucketMetrics: &scheduler.TokenBucketMetrics{
					FillRateGauge:        tokenBucketFillRateGauge,
					BucketCapacityGauge:  tokenBucketBucketCapacityGauge,
					AvailableTokensGauge: tokenBucketAvailableTokensGauge,
				},
			}

			// Initialize the token bucket
			lsa.tokenBucketLoadShed = scheduler.NewTokenBucketLoadShed(clock.Now(), tokenBucketMetrics)

			err = lsaFactory.loadShedDecisionWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				return retErr(err)
			}
			return retErr(nil)
		},
		OnStop: func(context.Context) error {
			var errMulti error
			err := lsaFactory.loadShedDecisionWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			deleted := lsaFactory.tokenBucketLSFGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete token_bucket_lsf gauge from its metric vector"))
			}
			deleted = lsaFactory.tokenBucketFillRateGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete token_bucket_fill_rate gauge from its metric vector"))
			}
			deleted = lsaFactory.tokenBucketBucketCapacityGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete token_bucket_capacity gauge from its metric vector"))
			}
			deleted = lsaFactory.tokenBucketAvailableTokensGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete token_bucket_available_tokens gauge from its metric vector"))
			}

			lsa.statusRegistry.SetStatus(status.NewStatus(nil, errMulti))
			return errMulti
		},
	})
	return lsa, nil
}

// loadShedActuator saves load shed decisions received from controller.
type loadShedActuator struct {
	conLimiter          *concurrencyLimiter
	clock               clockwork.Clock
	tokenBucketLoadShed *scheduler.TokenBucketLoadShed
	statusRegistry      status.Registry
}

func (lsa *loadShedActuator) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	if event.Type == notifiers.Remove {
		log.Debug().Msg("Decision was removed")
		return
	}

	var wrapperMessage wrappersv1.LoadShedDecsisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	loadShedDecision := wrapperMessage.LoadShedDecision
	if err != nil || loadShedDecision == nil {
		statusMsg := "Failed to unmarshal config wrapper"
		log.Warn().Err(err).Msg(statusMsg)
		lsa.statusRegistry.SetStatus(status.NewStatus(nil, err))
		return
	}
	// check if this decision is for the same policy id as what we have
	if wrapperMessage.PolicyHash != lsa.conLimiter.GetPolicyHash() {
		err = errors.New("policy id mismatch")
		statusMsg := fmt.Sprintf("Expected policy hash: %s, Got: %s", lsa.conLimiter.GetPolicyHash(), wrapperMessage.PolicyHash)
		log.Warn().Err(err).Msg(statusMsg)
		lsa.statusRegistry.SetStatus(status.NewStatus(nil, err))
		return
	}

	log.Trace().Float64("loadShedFactor", loadShedDecision.LoadShedFactor).Msg("Setting load shed factor")
	lsa.tokenBucketLoadShed.SetLoadShedFactor(lsa.clock.Now(), loadShedDecision.LoadShedFactor)
}
