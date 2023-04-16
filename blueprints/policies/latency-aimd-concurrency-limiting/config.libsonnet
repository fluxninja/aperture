local flow_selector_defaults = {
  service_selector: {
    service: '__REQUIRED_FIELD__',
  },
  flow_matcher: {
    control_point: '__REQUIRED_FIELD__',
  },
};

{
  /**
  * @param (common.policy_name: string required) Name of the policy.
  */
  common: {
    policy_name: '__REQUIRED_FIELD__',
  },
  policy: {
    /**
    * @param (policy.flux_meter: aperture.spec.v1.FluxMeter required) Flux Meter.
    */
    flux_meter: {
      flow_selector: flow_selector_defaults,
    },
    /**
    * @param (policy.classifiers: []aperture.spec.v1.Classifier) List of classification rules.
    */
    classifiers: [],
    /**
    * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
    */
    components: [],
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
    /**
    * @param (policy.concurrency_controller.flow_selector: aperture.spec.v1.FlowSelector required) Concurrency Limiter flow selector.
    * @param (policy.concurrency_controller.scheduler: aperture.spec.v1.SchedulerParameters) Scheduler parameters.
    * @param (policy.concurrency_controller.scheduler.auto_tokens: bool) Automatically estimate cost (tokens) for workload requests.
    * @param (policy.concurrency_controller.gradient: aperture.spec.v1.GradientControllerParameters) Gradient Controller parameters.
    * @param (policy.concurrency_controller.alerter: aperture.spec.v1.AlerterParameters) Whether tokens for workloads are computed dynamically or set statically by the user.
    * @param (policy.concurrency_controller.max_load_multiplier: float64) Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.
    * @param (policy.concurrency_controller.load_multiplier_linear_increment: float64) Linear increment to load multiplier in each execution tick (0.5s) when the system is not in overloaded state.
    * @param (policy.concurrency_controller.default_config: aperture.spec.v1.LoadActuatorDynamicConfig) Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.
    */
    concurrency_controller: {
      flow_selector: flow_selector_defaults,
      scheduler: {
        auto_tokens: true,
      },
      gradient: {
        slope: -1,
        min_gradient: 0.1,
        max_gradient: 1.0,
      },
      alerter: {
        alert_name: 'Load Shed Event',
      },
      max_load_multiplier: 2.0,
      load_multiplier_linear_increment: 0.0025,
      default_config: {
        dry_run: false,
      },
    },
  },
  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) From time of dashboard.
  * @param (dashboard.time_to: string) To time of dashboard.
  */
  dashboard: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    /**
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
  },
}
