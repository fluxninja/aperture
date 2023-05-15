local autoScalingDefaults = import '../../auto-scaling/base/config-defaults.libsonnet';
local averageLatencyServiceProtection = import '../../service-protection/average-latency/config.libsonnet';

/**
* @param (policy.policy_name: string required) Name of the policy.
* @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
* @param (policy.resources: aperture.spec.v1.Resources) Additional resources.
* @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
* @param (policy.service_protection_core.overload_confirmations: []policies/service-protection/average-latency:schema:overload_confirmation) List of overload confirmation criteria. Load scheduler can throttle flows when all of the specified overload confirmation criteria are met.
* @param (policy.service_protection_core.adaptive_load_scheduler: aperture.spec.v1.AdaptiveLoadSchedulerParameters required) Parameters for Adaptive Load Scheduler.
* @param (policy.service_protection_core.dry_run: bool) Default configuration for setting dry run mode on Load Scheduler. In dry run mode, the Load Scheduler acts as a passthrough and does not throttle flows. This config can be updated at runtime without restarting the policy.
* @param (policy.latency_baseliner.flux_meter: aperture.spec.v1.FluxMeter required) Flux Meter defines the scope of latency measurements.
* @param (policy.latency_baseliner.ema: aperture.spec.v1.EMAParameters) EMA parameters.
* @param (policy.latency_baseliner.latency_tolerance_multiplier: float64) Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.
* @param (policy.latency_baseliner.latency_ema_limit_multiplier: float64) Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.
* @param (policy.auto_scaling.kubernetes_replicas: aperture.spec.v1.AutoScalerScalingBackendKubernetesReplicas required) Kubernetes replicas scaling backend.
* @param (policy.auto_scaling.kubernetes_replicas.kubernetes_object_selector: aperture.spec.v1.KubernetesObjectSelector required) Kubernetes object selector.
* @param (policy.auto_scaling.kubernetes_replicas.min_replicas: string required) Minimum number of replicas.
* @param (policy.auto_scaling.kubernetes_replicas.max_replicas: string required) Maximum number of replicas.
* @param (policy.auto_scaling.dry_run: bool) Dry run mode ensures that no scaling is invoked by this auto scaler.
* @param (policy.auto_scaling.dry_run_config_key: string) Configuration key for overriding dry run setting through dynamic configuration.
* @param (policy.auto_scaling.scaling_parameters: aperture.spec.v1.AutoScalerScalingParameters required) Parameters that define the scaling behavior.
* @param (policy.auto_scaling.disable_periodic_scale_in: bool) Disable periodic scale in.
* @param (policy.auto_scaling.periodic_decrease: aperture.spec.v1.PeriodicDecreaseParameters) Parameters for periodic scale in.
* @param (policy.auto_scaling.periodic_decrease.period: string) Period for periodic scale in.
* @param (policy.auto_scaling.periodic_decrease.scale_in_percentage: int) Percentage of replicas to scale in.
*/

/**
* @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
* @param (dashboard.time_from: string) From time of dashboard.
* @param (dashboard.time_to: string) To time of dashboard.
* @param (dashboard.datasource.name: string) Datasource name.
* @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
*/

averageLatencyServiceProtection {
  policy+: {
    auto_scaling: {
      dry_run: autoScalingDefaults.policy.dry_run,
      dry_run_config_key: autoScalingDefaults.policy.dry_run_config_key,
      kubernetes_replicas+: {
        kubernetes_object_selector: '__REQUIRED_FIELD__',
        min_replicas: '__REQUIRED_FIELD__',
        max_replicas: '__REQUIRED_FIELD__',
      },
      scaling_parameters: autoScalingDefaults.policy.scaling_parameters,
      disable_periodic_scale_in: false,
      periodic_decrease: {
        period: '60s',
        scale_in_percentage: 10,
      },
    },
  },
}
