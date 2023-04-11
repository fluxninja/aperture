package controlplane

import (
	"context"
	"encoding/json"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdnotifier "github.com/fluxninja/aperture/pkg/etcd/notifier"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/crwatcher"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"sigs.k8s.io/yaml"
)

const (
	policiesTrackerFxTag             = "PoliciesTracker"
	policiesDynameConfigTrackerFxTag = "PoliciesDynamicConfigTracker"
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
		fx.Provide(providePolicyValidator,
			fx.Annotate(
				provideTrackers,
				fx.ResultTags(
					config.NameTag(policiesTrackerFxTag),
					config.NameTag(policiesDynameConfigTrackerFxTag),
				),
			),
		),
		// Syncing policies config to etcd
		crwatcher.Constructor{
			TrackersName:              policiesTrackerFxTag,
			DynamicConfigTrackersName: policiesDynameConfigTrackerFxTag,
		}.Annotate(), // Create a new watcher
		fx.Invoke(
			fx.Annotate(
				setupPoliciesNotifier,
				fx.ParamTags(
					config.NameTag(policiesTrackerFxTag),
					config.NameTag(policiesDynameConfigTrackerFxTag),
				),
			),
		),
		// Policy factory
		policyFactoryModule(),
	)
}

func provideTrackers(lifecycle fx.Lifecycle) (notifiers.Trackers, notifiers.Trackers, error) {
	trackers := notifiers.NewDefaultTrackers()
	dynamicConfigTrackers := notifiers.NewDefaultTrackers()

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return multierr.Combine(
				trackers.Start(),
				dynamicConfigTrackers.Start(),
			)
		},
		OnStop: func(ctx context.Context) error {
			return multierr.Combine(
				trackers.Stop(),
				dynamicConfigTrackers.Stop(),
			)
		},
	})

	return trackers, dynamicConfigTrackers, nil
}

// Sync policies config directory with etcd.
func setupPoliciesNotifier(crWatcher, dynamicConfigWatcher notifiers.Trackers, etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	if crWatcher == nil || dynamicConfigWatcher == nil {
		log.Debug().Msg("Kubernetes watcher is disabled")
		return
	}
	wrapPolicy := func(key notifiers.Key, bytes []byte, etype notifiers.EventType) (notifiers.Key, []byte, error) {
		var dat []byte
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

	crNotifier := etcdnotifier.NewPrefixToEtcdNotifier(
		paths.PoliciesConfigPath,
		etcdClient,
		true)
	// content transform callback to wrap policy in config properties wrapper
	crNotifier.SetTransformFunc(wrapPolicy)

	dynamicConfigNotifier := etcdnotifier.NewPrefixToEtcdNotifier(
		paths.PoliciesDynamicConfigPath,
		etcdClient,
		true)

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := crNotifier.Start()
			if err != nil {
				return err
			}
			err = crWatcher.AddPrefixNotifier(crNotifier)
			if err != nil {
				return err
			}

			err = dynamicConfigNotifier.Start()
			if err != nil {
				return err
			}
			err = dynamicConfigWatcher.AddPrefixNotifier(dynamicConfigNotifier)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			var merr, err error
			err = crWatcher.RemovePrefixNotifier(crNotifier)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = crNotifier.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}

			err = dynamicConfigWatcher.RemovePrefixNotifier(dynamicConfigNotifier)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = dynamicConfigNotifier.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})
}
