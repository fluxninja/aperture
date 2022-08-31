package controller

import (
	"context"
	"path"

	"github.com/ghodss/yaml"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdnotifier "github.com/fluxninja/aperture/pkg/etcd/notifier"
	filesystemwatcher "github.com/fluxninja/aperture/pkg/filesystem/watcher"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
)

var (
	// swagger:operation POST /controller common-configuration Controller
	// ---
	// x-fn-config-env: true
	// parameters:
	// - name: policies_path
	//   in: query
	//   type: string
	//   description: Directory containing policies rules
	//   x-go-default: "/etc/aperture/aperture-controller/policies"

	policiesDefaultPath = path.Join(config.DefaultAssetsDirectory, "policies")
	policiesPathKey     = "controller.policies_path"
	policiesFxTag       = "Policies"
)

// Module - Controller can be initialized by passing options from Module() to fx app.
func Module() fx.Option {
	return fx.Options(
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
func setupPoliciesNotifier(w notifiers.Watcher, etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	wrapPolicy := func(key notifiers.Key, bytes []byte, etype notifiers.EventType) (notifiers.Key, []byte, error) {
		var dat []byte
		switch etype {
		case notifiers.Write:
			unmarshaller, _ := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			policyMessage := &policylangv1.Policy{}
			unmarshalErr := unmarshaller.Unmarshal(policyMessage)
			if unmarshalErr != nil {
				log.Warn().Err(unmarshalErr).Msg("Failed to unmarshal policy")
				return key, nil, unmarshalErr
			}

			wrapper, wrapErr := controlplane.HashAndPolicyWrap(policyMessage, string(key))
			if wrapErr != nil {
				log.Warn().Err(wrapErr).Msg("Failed to wrap message in config properties")
				return key, nil, wrapErr
			}
			var marshalWrapErr error
			dat, marshalWrapErr = yaml.Marshal(wrapper)
			if marshalWrapErr != nil {
				log.Warn().Err(marshalWrapErr).Msgf("Failed to marshal config wrapper for proto message %+v", &wrapper)
				return key, nil, marshalWrapErr
			}
		}
		return key, dat, nil
	}

	notifier := etcdnotifier.NewPrefixToEtcdNotifier(
		paths.PoliciesConfigPath,
		etcdClient,
		true)
	// content transform callback to wrap policy in config properties wrapper
	notifier.SetTransformFunc(wrapPolicy)

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
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
		OnStop: func(_ context.Context) error {
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
