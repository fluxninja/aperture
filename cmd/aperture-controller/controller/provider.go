package controller

import (
	"context"
	"path"

	"github.com/ghodss/yaml"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	classificationv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/classification/v1"
	policylangv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/FluxNinja/aperture/pkg/config"
	etcdclient "github.com/FluxNinja/aperture/pkg/etcd/client"
	etcdnotifier "github.com/FluxNinja/aperture/pkg/etcd/notifier"
	filesystemwatcher "github.com/FluxNinja/aperture/pkg/filesystem/watcher"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/notifiers"
	"github.com/FluxNinja/aperture/pkg/paths"
	"github.com/FluxNinja/aperture/pkg/policies/controlplane"
	"github.com/FluxNinja/aperture/pkg/status"
	"github.com/FluxNinja/aperture/pkg/utils"
)

var (
	// swagger:operation POST /controller common-configuration Controller
	// ---
	// x-fn-config-env: true
	// parameters:
	// - name: classifiers_path
	//   in: query
	//   type: string
	//   description: Directory containing classification rules
	//   x-go-default: "/etc/aperture/aperture-controller/classifiers"
	// - name: policies_path
	//   in: query
	//   type: string
	//   description: Directory containing policies rules
	//   x-go-default: "/etc/aperture/aperture-controller/policies"

	classifiersDefaultPath = path.Join(config.DefaultAssetsDirectory, "classifiers")
	classifiersPathKey     = "controller.classifiers_path"
	classifiersFxTag       = "Classifiers"
	rulesetKey             = "classification.ruleset"

	policiesDefaultPath = path.Join(config.DefaultAssetsDirectory, "policies")
	policiesPathKey     = "controller.policies_path"
	policiesFxTag       = "Policies"
	policyKey           = "policy"
)

// Module - Controller can be initialized by passing options from Module() to fx app.
func Module() fx.Option {
	return fx.Options(
		// Syncing classifiers config to etcd
		fx.Provide(provideClassifiersPathFlag),
		filesystemwatcher.Constructor{Name: classifiersFxTag, PathKey: classifiersPathKey, Path: classifiersDefaultPath}.Annotate(), // Create a new filesystemwatcher
		fx.Invoke(
			fx.Annotate(
				setupClassifiersNotifier,
				fx.ParamTags(config.NameTag(classifiersFxTag)),
			),
		),
		// Syncing policies config to etcd
		fx.Provide(providePoliciesPathFlag),
		filesystemwatcher.Constructor{Name: policiesFxTag, PathKey: policiesPathKey, Path: policiesDefaultPath}.Annotate(), // Create a new watcher
		fx.Invoke(
			fx.Annotate(
				setupPoliciesNotifier,
				fx.ParamTags(config.NameTag(policiesFxTag)),
			),
		),
		// Policy factory
		controlplane.PolicyFactoryModule(),
	)
}

// provideClassifiersPathFlag registers a command line flag builder function.
func provideClassifiersPathFlag() config.FlagSetBuilderOut {
	return config.FlagSetBuilderOut{Builder: setClassifiersPathFlag}
}

// setClassifiersPathFlag registers command line flags.
func setClassifiersPathFlag(fs *pflag.FlagSet) error {
	fs.String(classifiersPathKey, classifiersDefaultPath, "path to classifiers directory")
	return nil
}

// Sync classifiers config directory with etcd.
func setupClassifiersNotifier(w notifiers.Watcher, etcdClient *etcdclient.Client, lifecycle fx.Lifecycle, statusRegistry *status.Registry) {
	wrapClassifier := func(key notifiers.Key, bytes []byte, etype notifiers.EventType) ([]byte, error) {
		unmarshaller, _ := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
		classifierMsg := &classificationv1.Classifier{}
		unmarshalErr := unmarshaller.Unmarshal(classifierMsg)
		if unmarshalErr != nil {
			log.Warn().Err(unmarshalErr).Msg("Failed to unmarshal classification")
			return nil, unmarshalErr
		}

		statusPath := rulesetKey + "." + key.String()

		if etype == notifiers.Write {
			s := status.NewStatus(classifierMsg, nil)
			pushErr := statusRegistry.Push(statusPath, s)
			if pushErr != nil {
				log.Error().Err(pushErr).Msg("could not push classification rules to status registry")
			}

			wrapper, wrapErr := utils.HashAndWrapWithConfProps(classifierMsg, string(key))
			if wrapErr != nil {
				log.Warn().Err(wrapErr).Msg("Failed to wrap message in config properties")
				return nil, wrapErr
			}
			dat, marshalWrapErr := yaml.Marshal(wrapper)
			if marshalWrapErr != nil {
				log.Warn().Err(marshalWrapErr).Msgf("Failed to marshal config wrapper for proto message %+v", &wrapper)
				return nil, marshalWrapErr
			}
			return dat, nil
		}

		if etype == notifiers.Remove {
			statusRegistry.Delete(statusPath)
		}

		return nil, nil
	}

	notifier := etcdnotifier.NewPrefixToEtcdNotifier(
		paths.Classifiers,
		etcdClient,
		true,
	)

	// content transform callback to wrap policy in config properties wrapper
	notifier.SetTransformFunc(wrapClassifier)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := notifier.Start()
			if err != nil {
				return err
			}
			err = w.AddPrefixNotifier(notifier)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			err = w.RemovePrefixNotifier(notifier)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = notifier.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})
}

// providePoliciesPathFlag registers a command line flag builder function.
func providePoliciesPathFlag() config.FlagSetBuilderOut {
	return config.FlagSetBuilderOut{Builder: setPoliciesPathFlag}
}

// setPoliciesPathFlag registers command line flags.
func setPoliciesPathFlag(fs *pflag.FlagSet) error {
	fs.String(policiesPathKey, policiesDefaultPath, "path to Policies directory")
	return nil
}

// Sync policies config directory with etcd.
func setupPoliciesNotifier(w notifiers.Watcher, etcdClient *etcdclient.Client, lifecycle fx.Lifecycle, statusRegistry *status.Registry) {
	wrapPolicy := func(key notifiers.Key, bytes []byte, etype notifiers.EventType) ([]byte, error) {
		statusPath := policyKey + "." + key.String()

		var dat []byte

		switch etype {
		case notifiers.Write:
			unmarshaller, _ := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			policyMessage := &policylangv1.Policy{}
			unmarshalErr := unmarshaller.Unmarshal(policyMessage)
			if unmarshalErr != nil {
				log.Warn().Err(unmarshalErr).Msg("Failed to unmarshal policy")
				return nil, unmarshalErr
			}

			wrapper, wrapErr := utils.HashAndWrapWithConfProps(policyMessage, string(key))
			if wrapErr != nil {
				log.Warn().Err(wrapErr).Msg("Failed to wrap message in config properties")
				return nil, wrapErr
			}
			var marshalWrapErr error
			dat, marshalWrapErr = yaml.Marshal(wrapper)
			if marshalWrapErr != nil {
				log.Warn().Err(marshalWrapErr).Msgf("Failed to marshal config wrapper for proto message %+v", &wrapper)
				return nil, marshalWrapErr
			}

			s := status.NewStatus(wrapper, nil)
			pushErr := statusRegistry.Push(statusPath, s)
			if pushErr != nil {
				log.Error().Err(pushErr).Msg("could not push classification rules to status registry")
			}

		case notifiers.Remove:
			statusRegistry.Delete(statusPath)
		}
		return dat, nil
	}

	notifier := etcdnotifier.NewPrefixToEtcdNotifier(
		paths.Policies,
		etcdClient,
		true)
	// content transform callback to wrap policy in config properties wrapper
	notifier.SetTransformFunc(wrapPolicy)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := notifier.Start()
			if err != nil {
				return err
			}
			err = w.AddPrefixNotifier(notifier)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			err = w.RemovePrefixNotifier(notifier)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = notifier.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})
}
