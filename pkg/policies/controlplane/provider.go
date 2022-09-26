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
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	kuberneteswatcher "github.com/fluxninja/aperture/pkg/policies/watcher/kubernetes"
)

// swagger:operation POST /policies common-configuration PoliciesConfig
// ---
// x-fn-config-env: true
// parameters:
// - name: promql_jobs_scheduler
//   in: body
//   schema:
//     "$ref": "#/definitions/JobGroupConfig"

// Module - Controller can be initialized by passing options from Module() to fx app.
func Module() fx.Option {
	policiesFxTag := "Policies"
	return fx.Options(
		fx.Provide(providePolicyValidator),
		// Syncing policies config to etcd
		kuberneteswatcher.Constructor{Name: policiesFxTag}.Annotate(), // Create a new watcher
		fx.Invoke(
			fx.Annotate(
				setupPoliciesNotifier,
				fx.ParamTags(config.NameTag(policiesFxTag)),
			),
		),
		// Policy factory
		policyFactoryModule(),
	)
}

// Sync policies config directory with etcd.
func setupPoliciesNotifier(w notifiers.Watcher, etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
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

	notifier := etcdnotifier.NewPrefixToEtcdNotifier(
		common.PoliciesConfigPath,
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
