package controlplane

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdnotifier "github.com/fluxninja/aperture/v2/pkg/etcd/notifier"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/crwatcher"
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
			ProvidePolicyValidator,
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
	etcdClient *etcdclient.Client,
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
			wrapper, err := policyWrapperFromJSONBytes(string(key), bytes, source)
			if err != nil {
				log.Warn().Err(err).Msgf("Failed to unmarshal policy wrapper from json bytes %+v", bytes)
				return key, nil, err
			}

			var marshalWrapErr error
			dat, marshalWrapErr = proto.Marshal(wrapper)
			if marshalWrapErr != nil {
				log.Warn().Err(marshalWrapErr).Msgf("Failed to marshal config wrapper for proto message %+v", &wrapper)
				return key, nil, marshalWrapErr
			}
		}
		return key, dat, nil
	}

	policyEtcdToEtcdNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesConfigPath, etcdClient)
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
	policyCRToEtcdNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesConfigPath, etcdClient)
	policyCRToEtcdNotifier.SetTransformFunc(
		func(key notifiers.Key, bytes []byte, etype notifiers.EventType) (notifiers.Key, []byte, error) {
			return wrapPolicy(key, bytes, etype, policysyncv1.PolicyWrapper_K8S)
		},
	)

	policyDynamicConfigEtcdNotifier := etcdnotifier.NewPrefixToEtcdNotifier(paths.PoliciesDynamicConfigPath, etcdClient)

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
