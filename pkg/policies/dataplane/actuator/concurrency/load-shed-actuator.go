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

	configv1 "aperture.tech/aperture/api/gen/proto/go/aperture/common/config/v1"
	policydecisionsv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	"aperture.tech/aperture/pkg/config"
	etcdclient "aperture.tech/aperture/pkg/etcd/client"
	etcdwatcher "aperture.tech/aperture/pkg/etcd/watcher"
	"aperture.tech/aperture/pkg/flowcontrol/scheduler"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/notifiers"
	"aperture.tech/aperture/pkg/paths"
	"aperture.tech/aperture/pkg/policies/dataplane/component"
	"aperture.tech/aperture/pkg/status"
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
	agentGroupName string,
	prometheusRegistry *prometheus.Registry,
) (*loadShedActuatorFactory, error) {
	// Scope the sync to the agent group.
	etcdPath := path.Join(paths.LoadShedDecisionsPath, paths.AgentGroupPrefix(agentGroupName))
	loadShedDecisionWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}
	f := &loadShedActuatorFactory{
		loadShedDecisionWatcher: loadShedDecisionWatcher,
		agentGroupName:          agentGroupName,
	}
	// Initialize and register the WFQ and Token Bucket Metric Vectors
	f.tokenBucketLSFGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_bucket_lsf",
			Help: "A gauge that tracks the load shed factor",
		},
		metricLabelKeys,
	)
	f.tokenBucketFillRateGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_bucket_bucket_fill_rate",
			Help: "A gauge that tracks the fill rate of token bucket",
		},
		metricLabelKeys,
	)
	f.tokenBucketBucketCapacityGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_bucket_bucket_capacity",
			Help: "A gauge that tracks the capacity of token bucket",
		},
		metricLabelKeys,
	)
	f.tokenBucketAvailableTokensGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_bucket_available_tokens",
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
				err := fmt.Errorf("failed to unregister token_bucket_bucket_fill_rate metric")
				errMulti = multierr.Append(errMulti, err)
			}
			if !prometheusRegistry.Unregister(f.tokenBucketBucketCapacityGaugeVec) {
				err := fmt.Errorf("failed to unregister token_bucket_bucket_capacity metric")
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
func (lsaFactory *loadShedActuatorFactory) newLoadShedActuator(registryPath string, componentAPI component.ComponentAPI, registry *status.Registry, clock clockwork.Clock, lifecycle fx.Lifecycle, metricLabels prometheus.Labels) (*loadShedActuator, error) {
	lsa := &loadShedActuator{
		componentAPI:   componentAPI,
		clock:          clock,
		statusRegistry: registry,
		registryPath:   fmt.Sprintf("%s.load_shed", registryPath),
	}

	unmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return nil, err
	}
	// decision notifier
	decisionNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(paths.IdentifierForComponent(lsaFactory.agentGroupName, lsa.componentAPI.GetPolicyName(), lsa.componentAPI.GetComponentIndex())),
		unmarshaller,
		lsa.decisionUpdateCallback,
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			retErr := func(err error) error {
				s := status.NewStatus(nil, err)
				errStatus := lsa.statusRegistry.Push(lsa.registryPath, s)
				if errStatus != nil {
					errStatus = errors.Wrap(errStatus, "failed to push status")
					return multierr.Append(err, errStatus)
				}
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
				errMulti = multierr.Append(errMulti, errors.New("failed to delete token_bucket_bucket_fill_rate gauge from its metric vector"))
			}
			deleted = lsaFactory.tokenBucketBucketCapacityGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete token_bucket_bucket_capacity gauge from its metric vector"))
			}
			deleted = lsaFactory.tokenBucketAvailableTokensGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete token_bucket_available_tokens gauge from its metric vector"))
			}

			s := status.NewStatus(nil, errMulti)
			err = lsa.statusRegistry.Push(lsa.registryPath, s)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}
			return errMulti
		},
	})
	return lsa, nil
}

// loadShedActuator saves load shed decisions received from controller.
type loadShedActuator struct {
	componentAPI        component.ComponentAPI
	clock               clockwork.Clock
	tokenBucketLoadShed *scheduler.TokenBucketLoadShed
	statusRegistry      *status.Registry
	registryPath        string
}

func (lsa *loadShedActuator) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	if event.Type == notifiers.Remove {
		log.Debug().Msg("Decision was removed")
		return
	}

	var wrapperMessage configv1.ConfigPropertiesWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Config == nil {
		statusMsg := "Failed to unmarshal config wrapper"
		log.Warn().Err(err).Msg(statusMsg)
		s := status.NewStatus(nil, err)
		rPErr := lsa.statusRegistry.Push(lsa.registryPath, s)
		if rPErr != nil {
			log.Error().Err(rPErr).Msg("Failed to push status")
		}
		return
	}
	// check if this decision is for the same policy id as what we have
	if wrapperMessage.PolicyHash != lsa.componentAPI.GetPolicyHash() {
		err = errors.New("policy id mismatch")
		statusMsg := fmt.Sprintf("Expected policy hash: %s, Got: %s", lsa.componentAPI.GetPolicyHash(), wrapperMessage.PolicyHash)
		log.Warn().Err(err).Msg(statusMsg)
		s := status.NewStatus(nil, err)
		rPErr := lsa.statusRegistry.Push(lsa.registryPath, s)
		if rPErr != nil {
			log.Error().Err(rPErr).Msg("Failed to push status")
		}
		return
	}
	var loadShedDecision policydecisionsv1.LoadShedDecision
	err = wrapperMessage.Config.UnmarshalTo(&loadShedDecision)
	if err != nil {
		statusMsg := "Failed to unmarshal policy decision"
		log.Warn().Err(err).Msg(statusMsg)
		s := status.NewStatus(nil, err)
		rPErr := lsa.statusRegistry.Push(lsa.registryPath, s)
		if rPErr != nil {
			log.Error().Err(rPErr).Msg("Failed to push status")
		}
		return
	}

	log.Info().Float64("loadShedFactor", loadShedDecision.LoadShedFactor).
		Msg("Setting load shed factor")
	lsa.tokenBucketLoadShed.SetLoadShedFactor(lsa.clock.Now(), loadShedDecision.LoadShedFactor)
}
