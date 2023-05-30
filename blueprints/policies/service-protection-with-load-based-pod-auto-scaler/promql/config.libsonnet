local baseServiceProtectionDefaults = import '../../auto-scaling/base/config-defaults.libsonnet';
local promqlServiceProtection = import '../../service-protection/promql/config.libsonnet';

/**
* @param (policy.policy_name: string required) Name of the policy.
* @param (policy.promql_query: string required) PromQL query.
* @param (policy.setpoint: float64 required) Setpoint.
* @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
* @param (policy.resources: aperture.spec.v1.Resources) Additional resources.
* @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
* @param (policy.service_protection_core.overload_confirmations: []policies/service-protection/promql:schema:overload_confirmation) List of overload confirmation criteria. Load scheduler can throttle flows when all of the specified overload confirmation criteria are met.
* @param (policy.service_protection_core.adaptive_load_scheduler: aperture.spec.v1.AdaptiveLoadSchedulerParameters required) Parameters for Adaptive Load Scheduler.
* @param (policy.service_protection_core.dry_run: bool) Default configuration for setting dry run mode on Load Scheduler. In dry run mode, the Load Scheduler acts as a passthrough and does not throttle flows. This config can be updated at runtime without restarting the policy.
* @param (policy.auto_scaling: policies/service-protection/average-latency:schema:auto_scaling) Auto scaling escalation.
*/

/**
* @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
* @param (dashboard.time_from: string) From time of dashboard.
* @param (dashboard.time_to: string) To time of dashboard.
* @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
* @param (dashboard.datasource.name: string) Datasource name.
* @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
*/

promqlServiceProtection {
  policy+: {
    auto_scaling: baseServiceProtectionDefaults.auto_scaling_pods,
  },
}
