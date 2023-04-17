package controlplane

import (
	"context"
	"encoding/json"

	"go.uber.org/fx"
	"go.uber.org/multierr"
	"sigs.k8s.io/yaml"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdnotifier "github.com/fluxninja/aperture/pkg/etcd/notifier"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/crwatcher"
	"github.com/fluxninja/aperture/pkg/policies/paths"
)

const (
	policiesTrackerFxTag                 = "PoliciesTracker"
	policiesDynamicConfigTrackerFxTag    = "PoliciesDynamicConfigTracker"
	policiesAPITrackerFxTag              = "PoliciesAPITracker"
	policiesAPIDynamicConfigTrackerFxTag = "PoliciesAPIDynamicConfigTracker"
)

// swagger:operation POST /policies common-configuration PoliciesConfig
// ---
// x-fn-config-env: true
// parameters:
// - name: promql_jobs_scheduler
//   in: body
//   schema:
//     "$ref": "#/definitions/JobGroupConfig"
// - name: cr_watcher
//   in: body
//   schema:
//     "$ref": "#/definitions/CRWatcherConfig"

// Module - Controller can be initialized by passing options from Module() to fx app.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			providePolicyValidator,
			fx.Annotate(
				provideTrackers,
				fx.ResultTags(
					config.NameTag(policiesTrackerFxTag),
					config.NameTag(policiesDynamicConfigTrackerFxTag),
				),
			),
		),
		// Create a new watcher for CRs
		crwatcher.Constructor{
			Name:              policiesTrackerFxTag,
			DynamicConfigName: policiesDynamicConfigTrackerFxTag,
		}.Annotate(),
		// Create a new watcher for policies API
		etcdwatcher.Constructor{
			Name:     policiesAPITrackerFxTag,
			EtcdPath: paths.PoliciesAPIConfigPath,
		}.Annotate(),
		// Create a new watcher for policies API dynamic config
		etcdwatcher.Constructor{
			Name:     policiesAPIDynamicConfigTrackerFxTag,
			EtcdPath: paths.PoliciesAPIDynamicConfigPath,
		}.Annotate(),
		fx.Invoke(
			fx.Annotate(
				setupPoliciesNotifier,
				fx.ParamTags(
					config.NameTag(policiesTrackerFxTag),
					config.NameTag(policiesDynamicConfigTrackerFxTag),
				),
			),
			fx.Annotate(
				setupPoliciesAPINotifier,
				fx.ParamTags(
					config.NameTag(policiesAPITrackerFxTag),
					config.NameTag(policiesAPIDynamicConfigTrackerFxTag),
				),
			),
		),
		// Policy factory
		policyFactoryModule(),
	)
}

func setupPoliciesAPINotifier(
	policiesWatcher notifiers.Watcher,
	dynamicConfigWatcher notifiers.Watcher,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
) {
	policiesNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesConfigPath, etcdClient, true)
	dynamicConfigNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesDynamicConfigPath, etcdClient, true)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return multierr.Combine(
				policiesWatcher.AddPrefixNotifier(policiesNotifier),
				dynamicConfigWatcher.AddPrefixNotifier(dynamicConfigNotifier),
			)
		},
		OnStop: func(ctx context.Context) error {
			return multierr.Combine(
				policiesWatcher.RemovePrefixNotifier(policiesNotifier),
				dynamicConfigWatcher.RemovePrefixNotifier(dynamicConfigNotifier),
			)
		},
	})
}

func provideTrackers(lifecycle fx.Lifecycle) (notifiers.Trackers, notifiers.Trackers, error) {
	policyTrackers := notifiers.NewDefaultTrackers()
	policyDynamicConfigTrackers := notifiers.NewDefaultTrackers()

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return multierr.Combine(
				policyTrackers.Start(),
				policyDynamicConfigTrackers.Start(),
			)
		},
		OnStop: func(ctx context.Context) error {
			return multierr.Combine(
				policyTrackers.Stop(),
				policyDynamicConfigTrackers.Stop(),
			)
		},
	})

	return policyTrackers, policyDynamicConfigTrackers, nil
}

// Sync policies config directory with etcd.
func setupPoliciesNotifier(
	policyTrackers notifiers.Trackers,
	policyDynamicConfigTrackers notifiers.Trackers,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
) {
	if policyTrackers == nil || policyDynamicConfigTrackers == nil {
		log.Debug().Msg("Kubernetes watcher is disabled")
		return
	}
	wrapPolicy := func(key notifiers.Key, bytes []byte, etype notifiers.EventType) (notifiers.Key, []byte, error) {
		var dat []byte
		if bytes == nil {
			return key, dat, nil
		}
		switch etype {
		case notifiers.Write:
			policyMessage := &policylangv1.Policy{}
			unmarshalErr := config.UnmarshalYAML(bytes, policyMessage)
			if unmarshalErr != nil {
				log.Warn().Err(unmarshalErr).Msg("Failed to unmarshal policy")
				return key, nil, unmarshalErr
			}
			wrapper, wrapErr := hashAndPolicyWrap(policyMessage, string(key))
			if wrapErr != nil {
				log.Warn().Err(wrapErr).Msg("Failed to wrap message in config properties")
				return key, nil, wrapErr
			}
			var marshalWrapErr error
			jsonDat, marshalWrapErr := json.Marshal(wrapper)
			if marshalWrapErr != nil {
				log.Warn().Err(marshalWrapErr).Msgf("Failed to marshal config wrapper for proto message %+v", &wrapper)
				return key, nil, marshalWrapErr
			}
			// convert to yaml
			dat, marshalWrapErr = yaml.JSONToYAML(jsonDat)
			if marshalWrapErr != nil {
				log.Warn().Err(marshalWrapErr).Msgf("Failed to marshal config wrapper for proto message %+v", &wrapper)
				return key, nil, marshalWrapErr
			}
		}
		return key, dat, nil
	}

	policyEtcdNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesConfigPath, etcdClient, true)
	// content transform callback to wrap policy in config properties wrapper
	policyEtcdNotifier.SetTransformFunc(wrapPolicy)

	policyDynamicConfigEtcdNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesDynamicConfigPath, etcdClient, true)

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := policyEtcdNotifier.Start()
			if err != nil {
				return err
			}
			err = policyTrackers.AddPrefixNotifier(policyEtcdNotifier)
			if err != nil {
				return err
			}

			err = policyDynamicConfigEtcdNotifier.Start()
			if err != nil {
				return err
			}
			err = policyDynamicConfigTrackers.AddPrefixNotifier(policyDynamicConfigEtcdNotifier)
			if err != nil {
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			var merr, err error
			err = policyTrackers.RemovePrefixNotifier(policyEtcdNotifier)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = policyEtcdNotifier.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}

			err = policyDynamicConfigTrackers.RemovePrefixNotifier(policyDynamicConfigEtcdNotifier)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = policyDynamicConfigEtcdNotifier.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}

			return merr
		},
	})
}
