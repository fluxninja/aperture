local serviceProtectionDefaults = import '../base/config-defaults.libsonnet';

/**
* @param (common.policy_name: string required) Name of the policy.
* @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
* @param (dashboard.time_from: string) From time of dashboard.
* @param (dashboard.time_to: string) To time of dashboard.
* @param (dashboard.datasource.name: string) Datasource name.
* @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
* @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
* @param (policy.resources: aperture.spec.v1.Resources) List of additional resources.
* @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
* @param (policy.service_protection_core: policies/service-protection/base:schema:service_protection_core required) Core parameters for Service Protection policy.
* @param (policy.service_protection_core.overload_confirmations: []overload_confirmation) List of overload confirmation criteria. Load scheduler can shed flows when all of the specified criteria are met.
* @param (policy.service_protection_core.adaptive_load_scheduler.flow_selector: aperture.spec.v1.FlowSelector required) Concurrency Limiter flow selector.
* @param (policy.service_protection_core.adaptive_load_scheduler.scheduler: aperture.spec.v1.Schedulerschemaeters) Scheduler schemaeters.
* @param (policy.service_protection_core.adaptive_load_scheduler.gradient: aperture.spec.v1.GradientControllerschemaeters) Gradient Controller schemaeters.
* @param (policy.service_protection_core.adaptive_load_scheduler.alerter: aperture.spec.v1.Alerterschemaeters) Whether tokens for workloads are computed dynamically or set statically by the user.
* @param (policy.service_protection_core.adaptive_load_scheduler.max_load_multiplier: float64) Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.
* @param (policy.service_protection_core.adaptive_load_scheduler.load_multiplier_linear_increment: float64) Linear increment to load multiplier in each execution tick (0.5s) when the system is not in overloaded state.
* @param (policy.service_protection_core.adaptive_load_scheduler.default_config: aperture.spec.v1.LoadActuatorDynamicConfig) Default configuration for concurrency controller that can be updated at the runtime without shutting down the
*/
serviceProtectionDefaults {
  /**
  */
  policy+: {
    /**
    * @param (policy.flux_meter: aperture.spec.v1.FluxMeter required) Flux Meter.
    */
    flux_meter: {
      flow_selector: serviceProtectionDefaults.flow_selector,
    },
    /**
    * @param (policy.latency_baseliner.ema: aperture.spec.v1.EMAParameters) EMA parameters.
    * @param (policy.latency_baseliner.latency_tolerance_multiplier: float64) Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.
    * @param (policy.latency_baseliner.latency_ema_limit_multiplier: float64) Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.
    */
    latency_baseliner: {
      ema: {
        ema_window: '1500s',
        warmup_window: '60s',
        correction_factor_on_max_envelope_violation: 0.95,
      },
      latency_tolerance_multiplier: 1.1,
      latency_ema_limit_multiplier: 2.0,
    },
  },
}
