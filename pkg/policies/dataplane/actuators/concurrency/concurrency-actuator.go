package concurrency

import (
	"context"
	"errors"
	"fmt"
	"path"
	"sync"

	policydecisionsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators/concurrency/scheduler"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"
)

type concurrencyActuatorFactory struct {
	concurrencyDecisionWatcher        notifiers.Watcher
	defaultConcurrencyDecisionWatcher notifiers.Watcher
	multiplierDecisionWatcher         notifiers.Watcher
	agentGroupName                    string
}

// newConcurrencyActuatorFactory sets up the load shed module in the main fx app.
func newConcurrencyActuatorFactory(
	lc fx.Lifecycle,
	etcdClient *etcdclient.Client,
	agentGroup string,
	_ *prometheus.Registry,
) (*concurrencyActuatorFactory, error) {
	// Scope the sync to the agent group.
	concurrencyDemandPath := path.Join(common.ConcurrencyDemandDecisionsPath, common.AgentPrefix(info.UUID)+"-"+common.AgentGroupPrefix(agentGroup))
	concurrencyDemandDecisionWatcher, err := etcdwatcher.NewWatcher(etcdClient, concurrencyDemandPath)
	if err != nil {
		return nil, err
	}
	concurrencyDemandDefaultPath := path.Join(common.ConcurrencyDemandDefaultDecisionsPath, common.AgentGroupPrefix(agentGroup))
	concurrencyDemandDefaultDecisionWatcher, err := etcdwatcher.NewWatcher(etcdClient, concurrencyDemandDefaultPath)
	if err != nil {
		return nil, err
	}
	concurrencyMultiplierPath := path.Join(common.ConcurrencyMultiplierDecisionsPath, common.AgentGroupPrefix(agentGroup))
	concurrencyMultiplierDecisionWatcher, err := etcdwatcher.NewWatcher(etcdClient, concurrencyMultiplierPath)
	if err != nil {
		return nil, err
	}
	f := &concurrencyActuatorFactory{
		concurrencyDecisionWatcher:        concurrencyDemandDecisionWatcher,
		defaultConcurrencyDecisionWatcher: concurrencyDemandDefaultDecisionWatcher,
		multiplierDecisionWatcher:         concurrencyMultiplierDecisionWatcher,
		agentGroupName:                    agentGroup,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err = concurrencyDemandDecisionWatcher.Start()
			if err != nil {
				return err
			}
			err = concurrencyDemandDefaultDecisionWatcher.Start()
			if err != nil {
				return err
			}
			err = concurrencyMultiplierDecisionWatcher.Start()
			if err != nil {
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			var errMulti error
			err := concurrencyMultiplierDecisionWatcher.Stop()
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}
			err = concurrencyDemandDefaultDecisionWatcher.Stop()
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}
			err = concurrencyDemandDecisionWatcher.Stop()
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			return errMulti
		},
	})
	return f, nil
}

// newConcurrencyActuator creates a new concurrency actuator based on proto spec.
func (caFactory *concurrencyActuatorFactory) newConcurrencyActuator(conLimiter *concurrencyLimiter,
	registry status.Registry,
	clock clockwork.Clock,
	lifecycle fx.Lifecycle,
	metricLabels prometheus.Labels,
) (*concurrencyActuator, error) {
	reg := registry.Child("concurrency_actuator")

	ca := &concurrencyActuator{
		conLimiter:     conLimiter,
		clock:          clock,
		statusRegistry: reg,
	}

	unmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return nil, err
	}

	// decision notifiers
	concurrencyNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(common.DataplaneComponentAgentKey(info.UUID, caFactory.agentGroupName, ca.conLimiter.GetPolicyName(), ca.conLimiter.GetComponentIndex())),
		unmarshaller,
		ca.concurrencyUpdateCallback,
	)
	multiplierNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(common.DataplaneComponentKey(caFactory.agentGroupName, ca.conLimiter.GetPolicyName(), ca.conLimiter.GetComponentIndex())),
		unmarshaller,
		ca.multiplierUpdateCallback,
	)
	defaultConcurrencyNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(common.DataplaneComponentKey(caFactory.agentGroupName, ca.conLimiter.GetPolicyName(), ca.conLimiter.GetComponentIndex())),
		unmarshaller,
		ca.defaultConcurrencyUpdateCallback,
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			retErr := func(err error) error {
				ca.statusRegistry.SetStatus(status.NewStatus(nil, err))
				return err
			}

			// Initialize the token bucket

			// TODO: token bucket metrics lifecycle needs to be moved to token bucket
			ca.basicTokenBucket = scheduler.NewBasicTokenBucket(clock.Now(), 0.0, nil)
			// Set pass through initially
			ca.basicTokenBucket.SetPassThrough(true)

			err = caFactory.concurrencyDecisionWatcher.AddKeyNotifier(concurrencyNotifier)
			if err != nil {
				return retErr(err)
			}
			err = caFactory.multiplierDecisionWatcher.AddKeyNotifier(multiplierNotifier)
			if err != nil {
				return retErr(err)
			}
			err = caFactory.defaultConcurrencyDecisionWatcher.AddKeyNotifier(defaultConcurrencyNotifier)
			if err != nil {
				return retErr(err)
			}
			return retErr(nil)
		},
		OnStop: func(context.Context) error {
			var errMulti error
			err := caFactory.defaultConcurrencyDecisionWatcher.RemoveKeyNotifier(defaultConcurrencyNotifier)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}
			err = caFactory.multiplierDecisionWatcher.RemoveKeyNotifier(multiplierNotifier)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}
			err = caFactory.concurrencyDecisionWatcher.RemoveKeyNotifier(concurrencyNotifier)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			ca.statusRegistry.SetStatus(status.NewStatus(nil, errMulti))
			return errMulti
		},
	})
	return ca, nil
}

type concurrencyActuator struct {
	// A reader writer lock for concurrent access of decisions
	lock                          sync.Mutex
	conLimiter                    *concurrencyLimiter
	clock                         clockwork.Clock
	statusRegistry                status.Registry
	concurrencyDecision           *policydecisionsv1.ConcurrencyDemandDecision
	defaultConcurrencyDecision    *policydecisionsv1.ConcurrencyDemandDecision
	concurrencyMultiplierDecision *policydecisionsv1.ConcurrencyMultiplierDecision
	basicTokenBucket              *scheduler.BasicTokenBucket
}

func (ca *concurrencyActuator) concurrencyUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := ca.statusRegistry.GetLogger()
	ca.lock.Lock()
	defer ca.lock.Unlock()
	logger.Info().Msg("concurrencyUpdateCallback")

	concurrencyDecision, err := ca.processConcurrencyEvent(event, unmarshaller)
	if err == nil {
		ca.concurrencyDecision = concurrencyDecision
	}

	ca.updateConcurrency()
}

func (ca *concurrencyActuator) defaultConcurrencyUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := ca.statusRegistry.GetLogger()
	ca.lock.Lock()
	defer ca.lock.Unlock()
	logger.Info().Msg("defaultConcurrencyUpdateCallback")

	defaultConcurrencyDecision, err := ca.processConcurrencyEvent(event, unmarshaller)
	if err == nil {
		ca.defaultConcurrencyDecision = defaultConcurrencyDecision
	}

	ca.updateConcurrency()
}

func (ca *concurrencyActuator) multiplierUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := ca.statusRegistry.GetLogger()
	ca.lock.Lock()
	defer ca.lock.Unlock()
	logger.Info().Msg("multiplierUpdateCallback")

	concurrencyMultiplierDecision, err := ca.processMultiplierEvent(event, unmarshaller)
	if err == nil {
		ca.concurrencyMultiplierDecision = concurrencyMultiplierDecision
	}

	ca.updateConcurrency()
}

func (ca *concurrencyActuator) updateConcurrency() {
	logger := ca.statusRegistry.GetLogger()
	if ca.concurrencyMultiplierDecision != nil {
		if ca.concurrencyDecision == nil {
			logger.Info().Msg("No concurrency decision found, using default concurrency decision")
			if ca.defaultConcurrencyDecision != nil {
				concurrency := ca.defaultConcurrencyDecision.GetDemand() * ca.concurrencyMultiplierDecision.Multiplier
				ca.basicTokenBucket.SetFillRate(ca.clock.Now(), concurrency)
				ca.basicTokenBucket.SetPassThrough(false)
			} else {
				logger.Info().Msg("No default concurrency decision found, set pass through")
				ca.basicTokenBucket.SetPassThrough(true)
			}
		} else {
			logger.Info().Msg("Using concurrency decision")
			// Ensure that tick concurrency decision and multiplier's ticks match
			if ca.concurrencyDecision.GetTick() == ca.concurrencyMultiplierDecision.GetTick() {
				concurrency := ca.concurrencyDecision.GetDemand() * ca.concurrencyMultiplierDecision.Multiplier
				ca.basicTokenBucket.SetFillRate(ca.clock.Now(), concurrency)
				ca.basicTokenBucket.SetPassThrough(false)
			} else {
				logger.Debug().Msg("Concurrency decision and multiplier's ticks do not match. Wait until they match to update concurrency")
			}
		}
	} else {
		ca.basicTokenBucket.SetPassThrough(true)
	}
}

func (ca *concurrencyActuator) processConcurrencyEvent(event notifiers.Event, unmarshaller config.Unmarshaller) (*policydecisionsv1.ConcurrencyDemandDecision, error) {
	logger := ca.statusRegistry.GetLogger()
	if event.Type == notifiers.Remove {
		return nil, nil
	} else {
		var wrapperMessage wrappersv1.ConcurrencyDemandDecisionWrapper
		err := unmarshaller.Unmarshal(&wrapperMessage)
		concurrencyDecision := wrapperMessage.ConcurrencyDemandDecision
		if err != nil || concurrencyDecision == nil {
			statusMsg := "Failed to unmarshal config wrapper"
			logger.Warn().Err(err).Msg(statusMsg)
			ca.statusRegistry.SetStatus(status.NewStatus(nil, err))
			return nil, err
		}
		commonAttributes := wrapperMessage.GetCommonAttributes()
		if commonAttributes == nil {
			statusMsg := "Failed to get common attributes from config wrapper"
			logger.Error().Err(err).Msg(statusMsg)
			ca.statusRegistry.SetStatus(status.NewStatus(nil, err))
			return nil, err
		}
		// check if this decision is for the same policy id as what we have
		if commonAttributes.PolicyHash != ca.conLimiter.GetPolicyHash() {
			err = errors.New("policy id mismatch")
			statusMsg := fmt.Sprintf("Expected policy hash: %s, Got: %s", ca.conLimiter.GetPolicyHash(), commonAttributes.PolicyHash)
			logger.Warn().Err(err).Msg(statusMsg)
			ca.statusRegistry.SetStatus(status.NewStatus(nil, err))
			return nil, err
		}
		return concurrencyDecision, nil
	}
}

func (ca *concurrencyActuator) processMultiplierEvent(event notifiers.Event, unmarshaller config.Unmarshaller) (*policydecisionsv1.ConcurrencyMultiplierDecision, error) {
	logger := ca.statusRegistry.GetLogger()
	if event.Type == notifiers.Remove {
		return nil, nil
	} else {
		var wrapperMessage wrappersv1.ConcurrencyMultiplierDecisionWrapper
		err := unmarshaller.Unmarshal(&wrapperMessage)
		multiplierDecision := wrapperMessage.ConcurrencyMultiplierDecision
		if err != nil || multiplierDecision == nil {
			statusMsg := "Failed to unmarshal config wrapper"
			logger.Warn().Err(err).Msg(statusMsg)
			ca.statusRegistry.SetStatus(status.NewStatus(nil, err))
			return nil, err
		}
		commonAttributes := wrapperMessage.GetCommonAttributes()
		if commonAttributes == nil {
			statusMsg := "Failed to get common attributes from config wrapper"
			logger.Error().Err(err).Msg(statusMsg)
			ca.statusRegistry.SetStatus(status.NewStatus(nil, err))
			return nil, err
		}
		// check if this decision is for the same policy id as what we have
		if commonAttributes.PolicyHash != ca.conLimiter.GetPolicyHash() {
			err = errors.New("policy id mismatch")
			statusMsg := fmt.Sprintf("Expected policy hash: %s, Got: %s", ca.conLimiter.GetPolicyHash(), commonAttributes.PolicyHash)
			logger.Warn().Err(err).Msg(statusMsg)
			ca.statusRegistry.SetStatus(status.NewStatus(nil, err))
			return nil, err
		}
		return multiplierDecision, nil
	}
}
