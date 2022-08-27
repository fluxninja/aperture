package classifier

import (
	"context"
	"path"

	"go.uber.org/fx"

	classificationv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/paths"
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
		fx.Annotated{
			Target: ProvideEmptyClassifier,
			Name:   "empty",
		},
		ProvideClassifier,
	),
)

func setupEtcdClassifierWatcher(etcdClient *etcdclient.Client, lc fx.Lifecycle, ai *agentinfo.AgentInfo) (notifiers.Watcher, error) {
	agentGroup := ai.GetAgentGroup()
	etcdPath := path.Join(paths.ClassifiersConfigPath, paths.AgentGroupPrefix(agentGroup))
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

// ProvideEmptyClassifier provides a classifier that is empty
//
// The classifier could be populated by calling UpdateRules.
func ProvideEmptyClassifier() *Classifier { return New() }

// ProvideClassifierIn holds parameters for ProvideClassifier.
type ProvideClassifierIn struct {
	fx.In
	Classifier *Classifier       `name:"empty"`
	Watcher    notifiers.Watcher `name:"classifier"`
	Lifecycle  fx.Lifecycle
	Registry   *status.Registry
}

// ProvideClassifier provides a classifier that loads the rules from config file.
func ProvideClassifier(in ProvideClassifierIn) *Classifier {
	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: []notifiers.FxOptionsFunc{in.Classifier.provideClassifierFxOptions},
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.NewProtobufUnmarshaller,
		},
		StatusRegistry: in.Registry,
		StatusPath:     "classifier-driver",
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return in.Watcher.AddPrefixNotifier(fxDriver)
		},
		OnStop: func(context.Context) error {
			return in.Watcher.RemovePrefixNotifier(fxDriver)
		},
	})
	return in.Classifier
}

// Per classifier fx app.
func (c *Classifier) provideClassifierFxOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry *status.Registry,
) (fx.Option, error) {
	return fx.Options(
		fx.Supply(c),
		fx.Invoke(invokeMiniApp),
	), nil
}

func invokeMiniApp(
	lc fx.Lifecycle,
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	classifier *Classifier,
) error {
	var rs classificationv1.Classifier
	if err := unmarshaller.Unmarshal(&rs); err != nil {
		return err
	}

	var activeRuleset ActiveRuleset

	lc.Append(
		fx.Hook{
			OnStart: func(startCtx context.Context) error {
				var err error
				activeRuleset, err = classifier.AddRules(startCtx, string(key), &rs)
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
