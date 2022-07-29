package classification

import (
	"context"
	"path"

	"go.uber.org/fx"

	classificationv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/classification/v1"
	configv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/common/config/v1"
	"github.com/FluxNinja/aperture/pkg/agentinfo"
	"github.com/FluxNinja/aperture/pkg/config"
	etcdclient "github.com/FluxNinja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/FluxNinja/aperture/pkg/etcd/watcher"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/notifiers"
	"github.com/FluxNinja/aperture/pkg/paths"
	"github.com/FluxNinja/aperture/pkg/status"
	"github.com/FluxNinja/aperture/pkg/webhooks/validation"
)

// Module is a default set of components to enable flow classification
//
// Note: this module provides just a Classifier datastructure, with no API endpoint.
// Example API endpoint to the classifier is pkg/apis/authz.
var Module fx.Option = fx.Options(
	fx.Provide(
		fx.Annotated{
			Target: setupEtcdClassifierWatcher,
			Name:   fxTag,
		},
		fx.Annotated{
			Target: ProvideEmptyClassifier,
			Name:   "empty",
		},
		ProvideClassifier,
	),
)

var agentGroup string

const (
	configKey              = "classification"
	classificationJobGroup = "classification"
	fxTag                  = "classifier"
)

func setupEtcdClassifierWatcher(etcdClient *etcdclient.Client, lc fx.Lifecycle, ai *agentinfo.AgentInfo) (notifiers.Watcher, error) {
	agentGroup = ai.GetAgentGroup()
	etcdPath := path.Join(paths.Classifiers)
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
			GetUnmarshallerFunc: config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
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

// ProvideCMFileValidator provides classification config map file validator
//
// Note: This validator must be registered to be accessible.
func ProvideCMFileValidator() *CMFileValidator {
	return &CMFileValidator{}
}

// RegisterCMFileValidator registers classification configmap validator as configmap file validator.
func RegisterCMFileValidator(validator *CMFileValidator, configMapValidator *validation.CMValidator) {
	// The path is not configurable – if one doesn't want default path, one
	// could just write their own Register function
	configMapValidator.RegisterCMFileValidator(validator)
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
	var wrapperMessage configv1.ConfigPropertiesWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Config == nil {
		if err != nil {
			return err
		}
	}

	var rs classificationv1.Classifier
	if err := wrapperMessage.Config.UnmarshalTo(&rs); err != nil {
		return err
	}

	if rs.Selector.AgentGroup != agentGroup {
		log.Trace().Msg("Could not create classifier - agent group mismatch")
		return nil
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
