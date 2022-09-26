package classifier

import (
	"context"
	"path"

	"go.uber.org/fx"

	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/status"
)

// Module is a default set of components to enable flow classification
//
// Note: this module provides just a Classifier datastructure, with no API endpoint.
// Example API endpoint to the classifier is pkg/envoy.
var Module fx.Option = fx.Options(
	fx.Provide(
		fx.Annotated{
			Target: setupEtcdClassifierWatcher,
			Name:   "classifier",
		},
		ProvideClassificationEngine,
	),
)

func setupEtcdClassifierWatcher(etcdClient *etcdclient.Client, lc fx.Lifecycle, ai *agentinfo.AgentInfo) (notifiers.Watcher, error) {
	agentGroup := ai.GetAgentGroup()
	etcdPath := path.Join(common.ClassifiersConfigPath, common.AgentGroupPrefix(agentGroup))
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
	Watcher   notifiers.Watcher `name:"classifier"`
	Lifecycle fx.Lifecycle
	Registry  status.Registry
}

// ProvideClassificationEngine provides a classifier that loads the rules from config file.
func ProvideClassificationEngine(in ClassificationEngineIn) *ClassificationEngine {
	reg := in.Registry.Child("classifiers")

	classificationEngine := NewClassificationEngine(reg)

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: []notifiers.FxOptionsFunc{classificationEngine.provideClassifierFxOptions},
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.NewProtobufUnmarshaller,
		},
		StatusRegistry: reg,
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return in.Watcher.AddPrefixNotifier(fxDriver)
		},
		OnStop: func(context.Context) error {
			return in.Watcher.RemovePrefixNotifier(fxDriver)
		},
	})
	return classificationEngine
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
	wrapperMessage := &wrappersv1.ClassifierWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.Classifier == nil {
		logger.Warn().Err(err).Msg("Failed to unmarshal classifier config wrapper")
		return err
	}
	var activeRuleset ActiveRuleset
	lc.Append(
		fx.Hook{
			OnStart: func(startCtx context.Context) error {
				var err error
				activeRuleset, err = c.AddRules(startCtx, string(key), wrapperMessage)
				if err != nil {
					return err
				}

				return nil
			},
			OnStop: func(_ context.Context) error {
				activeRuleset.Drop()
				return nil
			},
		},
	)

	return nil
}
