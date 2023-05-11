package classifier

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/agentinfo"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// Module is a default set of components to enable flow classification
//
// Note: this module provides just a Classifier datastructure, with no API endpoint.
// Example API endpoint to the classifier is pkg/envoy.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: setupEtcdClassifierWatcher,
				Name:   "classifier",
			},
			ProvideClassificationEngine,
		),
	)
}

func setupEtcdClassifierWatcher(etcdClient *etcdclient.Client, lc fx.Lifecycle, ai *agentinfo.AgentInfo) (notifiers.Watcher, error) {
	agentGroup := ai.GetAgentGroup()
	etcdPath := path.Join(paths.ClassifiersPath,
		paths.AgentGroupPrefix(agentGroup))
	etcdWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return etcdWatcher.Start()
		},
		OnStop: func(_ context.Context) error {
			return etcdWatcher.Stop()
		},
	})

	return etcdWatcher, nil
}

// ClassificationEngineIn holds parameters for ProvideClassificationEngine.
type ClassificationEngineIn struct {
	fx.In
	Watcher      notifiers.Watcher `name:"classifier"`
	Lifecycle    fx.Lifecycle
	Registry     status.Registry
	PromRegistry *prometheus.Registry
	AgentInfo    *agentinfo.AgentInfo
}

// ProvideClassificationEngine provides a classifier that loads the rules from config file.
func ProvideClassificationEngine(in ClassificationEngineIn) (iface.ClassificationEngine, *ClassificationEngine) {
	reg := in.Registry.Child("resource", "classifiers")

	classificationEngine := NewClassificationEngine(in.AgentInfo, reg)

	fxDriver, err := notifiers.NewFxDriver(reg, in.PromRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{classificationEngine.provideClassifierFxOptions})
	if err != nil {
		reg.GetLogger().Fatal().Err(err).Msg("Failed to create fx driver")
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := in.PromRegistry.Register(classificationEngine.counterVec)
			if err != nil {
				return err
			}
			err = in.Watcher.AddPrefixNotifier(fxDriver)
			if err != nil {
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			err := in.Watcher.RemovePrefixNotifier(fxDriver)
			if err != nil {
				return err
			}
			if !in.PromRegistry.Unregister(classificationEngine.counterVec) {
				return fmt.Errorf("failed to unregister %s metric", metrics.ClassifierCounterMetricName)
			}

			return nil
		},
	})
	// Return the same object once as an interface and once as a normal classifier engine - for authz
	return classificationEngine, classificationEngine
}

// Per classifier fx app.
func (c *ClassificationEngine) provideClassifierFxOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry status.Registry,
) (fx.Option, error) {
	return fx.Options(
		fx.Invoke(c.invokeMiniApp),
	), nil
}

func (c *ClassificationEngine) invokeMiniApp(
	lc fx.Lifecycle,
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
) error {
	logger := c.registry.GetLogger()
	wrapperMessage := &policysyncv1.ClassifierWrapper{}
	errM := unmarshaller.Unmarshal(wrapperMessage)
	if errM != nil || wrapperMessage.Classifier == nil {
		logger.Warn().Err(errM).Msg("Failed to unmarshal classifier config wrapper")
		return errM
	}
	var activeRuleset ActiveRuleset
	classifier := Classifier{
		classifierProto: wrapperMessage.GetClassifier(),
		classifierID: iface.ClassifierID{
			PolicyName:      wrapperMessage.ClassifierAttributes.PolicyName,
			PolicyHash:      wrapperMessage.ClassifierAttributes.PolicyHash,
			ClassifierIndex: wrapperMessage.ClassifierAttributes.ClassifierIndex,
		},
	}

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = wrapperMessage.ClassifierAttributes.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = wrapperMessage.ClassifierAttributes.GetPolicyHash()
	metricLabels[metrics.ClassifierIndexLabel] = strconv.FormatInt(wrapperMessage.ClassifierAttributes.GetClassifierIndex(), 10)

	lc.Append(
		fx.Hook{
			OnStart: func(startCtx context.Context) error {
				counter, err := c.counterVec.GetMetricWith(metricLabels)
				if err != nil {
					return fmt.Errorf("%w: failed to get classifier counter from vector", err)
				}
				classifier.counter = counter

				err = c.RegisterClassifier(&classifier)
				if err != nil {
					return err
				}

				activeRuleset, err = c.AddRules(startCtx, string(key), wrapperMessage)
				if err != nil {
					return err
				}
				return nil
			},
			OnStop: func(_ context.Context) error {
				var errMulti error
				activeRuleset.Drop()

				err := c.UnregisterClassifier(&classifier)
				if err != nil {
					errMulti = multierr.Append(errMulti, err)
				}

				deleted := c.counterVec.Delete(metricLabels)
				if !deleted {
					errMulti = multierr.Append(errMulti, errors.New("failed to delete classifier_counter from its metric vector"))
				}
				return errMulti
			},
		},
	)

	return nil
}
