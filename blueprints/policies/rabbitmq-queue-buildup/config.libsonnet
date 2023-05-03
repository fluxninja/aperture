{
  /**
  * @param (common.policy_name: string required) Name of the policy.
  * @param (common.queue_name: string required) Name of the queue to watch for buildup.
  */
  common: {
    policy_name: '__REQUIRED_FIELD__',
    queue_name: '__REQUIRED_FIELD__',
  },
  policy: {
    /**
    * @param (policy.classifiers: []aperture.spec.v1.Classifier) List of classification rules.
    */
    classifiers: [],
    /**
    * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
    */
    components: [],
    /**
    * @param (policy.concurrency_controller.selectors: []aperture.spec.v1.Selector required) Concurrency Limiter flow selectors.
    * @param (policy.concurrency_controller.scheduler: aperture.spec.v1.SchedulerParameters) Scheduler parameters.
    * @param (policy.concurrency_controller.scheduler.auto_tokens: bool) Automatically estimate cost (tokens) for workload requests.
    * @param (policy.concurrency_controller.gradient: aperture.spec.v1.GradientControllerParameters) Gradient Controller parameters.
    * @param (policy.concurrency_controller.alerter: aperture.spec.v1.AlerterParameters) Whether tokens for workloads are computed dynamically or set statically by the user.
    * @param (policy.concurrency_controller.max_load_multiplier: float64) Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.
    * @param (policy.concurrency_controller.queue_buildup_setpoint: float64) Queue buildup setpoint in number of messages.
    * @param (policy.concurrency_controller.load_multiplier_linear_increment: float64) Linear increment to load multiplier in each execution tick (0.5s) when the system is not in overloaded state.
    * @param (policy.concurrency_controller.default_config: aperture.spec.v1.LoadActuatorDynamicConfig) Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.
    */
    concurrency_controller: {
      selectors: [{
        service: '__REQUIRED_FIELD__',
        control_point: '__REQUIRED_FIELD__',
      }],
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
      queue_buildup_setpoint: '__REQUIRED_FIELD__',
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
