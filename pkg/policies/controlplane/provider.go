package controlplane

import (
	"context"
	"encoding/json"

	"go.uber.org/fx"
	"go.uber.org/multierr"
	"sigs.k8s.io/yaml"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdnotifier "github.com/fluxninja/aperture/v2/pkg/etcd/notifier"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/crwatcher"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
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
					config.NameTag(policiesAPITrackerFxTag),
					config.NameTag(policiesAPIDynamicConfigTrackerFxTag),
				),
			),
		),
		// Policy factory
		policyFactoryModule(),
	)
}

// provideTrackers provides new trackers for policies and policies dynamic config.
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
	policyAPIWatcher notifiers.Watcher,
	policyAPIDynamicConfigWatcher notifiers.Watcher,
	scopedKV *etcdclient.SessionScopedKV,
	lifecycle fx.Lifecycle,
) {
	wrapPolicy := func(
		key notifiers.Key,
		bytes []byte,
		etype notifiers.EventType,
		source policysyncv1.PolicyWrapper_Source,
	) (notifiers.Key, []byte, error) {
		var dat []byte
		if bytes == nil {
			return key, dat, nil
		}
		switch etype {
		case notifiers.Write:
			policyMessage := &iface.PolicyMessage{
				Policy: &languagev1.Policy{},
			}
			unmarshalErr := json.Unmarshal(bytes, policyMessage)
			if unmarshalErr != nil {
				log.Warn().Err(unmarshalErr).Msg("Failed to unmarshal policy")
				return key, nil, unmarshalErr
			}
			wrapper, wrapErr := hashAndPolicyWrap(policyMessage.Policy, string(key))
			if wrapErr != nil {
				log.Warn().Err(wrapErr).Msg("Failed to wrap message in config properties")
				return key, nil, wrapErr
			}
			wrapper.Source = source
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

	policyEtcdToEtcdNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesConfigPath, &scopedKV.KVWrapper)
	// content transform callback to wrap policy in config properties wrapper
	policyEtcdToEtcdNotifier.SetTransformFunc(
		func(key notifiers.Key, bytes []byte, etype notifiers.EventType) (notifiers.Key, []byte, error) {
			return wrapPolicy(key, bytes, etype, policysyncv1.PolicyWrapper_ETCD)
		},
	)

	// FIXME: Don't use multiple etcdnotifiers writing to the same prefix.
	// Note: This solution is not perfect, but still works:
	// * The whole etcd subtree is purged on Start(), but that shouldn't be a
	//   problem, as no other logic runs between these starts.
	// * Multiple notifiers writing the same key to etcd are a problem, but
	//   it's not handled correctly anyway (FIXME).
	policyCRToEtcdNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesConfigPath, &scopedKV.KVWrapper)
	policyCRToEtcdNotifier.SetTransformFunc(
		func(key notifiers.Key, bytes []byte, etype notifiers.EventType) (notifiers.Key, []byte, error) {
			return wrapPolicy(key, bytes, etype, policysyncv1.PolicyWrapper_K8S)
		},
	)

	policyDynamicConfigEtcdNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesDynamicConfigPath, &scopedKV.KVWrapper)

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return multierr.Combine(
				policyEtcdToEtcdNotifier.Start(),
				policyCRToEtcdNotifier.Start(),
				policyTrackers.AddPrefixNotifier(policyCRToEtcdNotifier),
				policyAPIWatcher.AddPrefixNotifier(policyEtcdToEtcdNotifier),
				policyDynamicConfigEtcdNotifier.Start(),
				policyDynamicConfigTrackers.AddPrefixNotifier(policyDynamicConfigEtcdNotifier),
				policyAPIDynamicConfigWatcher.AddPrefixNotifier(policyDynamicConfigEtcdNotifier),
			)
		},
		OnStop: func(_ context.Context) error {
			return multierr.Combine(
				policyTrackers.RemovePrefixNotifier(policyCRToEtcdNotifier),
				policyAPIWatcher.RemovePrefixNotifier(policyEtcdToEtcdNotifier),
				policyEtcdToEtcdNotifier.Stop(),
				policyCRToEtcdNotifier.Stop(),
				policyDynamicConfigTrackers.RemovePrefixNotifier(policyDynamicConfigEtcdNotifier),
				policyAPIDynamicConfigWatcher.RemovePrefixNotifier(policyDynamicConfigEtcdNotifier),
				policyDynamicConfigEtcdNotifier.Stop(),
			)
		},
	})
}
