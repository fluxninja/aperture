package controlplane

import (
	"github.com/fluxninja/aperture/pkg/policies/controlplane/crwatcher"
	"go.uber.org/fx"
)

const (
	policiesFxTag              = "Policies"
	policiesDynamicConfigFxTag = "PoliciesDynamicConfig"
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
	return fx.Options(
		fx.Provide(providePolicyValidator),
		// Syncing policies config to etcd
		crwatcher.Constructor{Name: policiesFxTag, DynamicConfigName: policiesDynamicConfigFxTag}.Annotate(), // Create a new watcher
		// Policy factory
		policyFactoryModule(),
	)
}
