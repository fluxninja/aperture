local selectors_defaults = [{
  service: '__REQUIRED_FIELD__',
  control_point: '__REQUIRED_FIELD__',
}];

local service_protection_core_defaults = {
  overload_confirmations: [],

  adaptive_load_scheduler: {
    load_scheduler: {
      selectors: selectors_defaults,
    },
    gradient: {
      slope: -1,
      min_gradient: 0.1,
      max_gradient: 1.0,
    },
    alerter: {
      alert_name: 'Load Throttling Event',
    },
    max_load_multiplier: 2.0,
    load_multiplier_linear_increment: 0.025,
  },

  dry_run: false,
};

{
  /**
  * @param (policy.policy_name: string) Name of the policy.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  * @param (policy.resources: aperture.spec.v1.Resources) Additional resources.
  * @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
  * @param (policy.service_protection_core.overload_confirmations: []overload_confirmation) List of overload confirmation criteria. Load scheduler can throttle flows when all of the specified overload confirmation criteria are met.
  * @schema (overload_confirmation.query_string: string) The Prometheus query to be run. Must return a scalar or a vector with a single element.
  * @schema (overload_confirmation.threshold: float64) The threshold for the overload confirmation criteria.
  * @schema (overload_confirmation.operator: string) The operator for the overload confirmation criteria. oneof: `gt | lt | gte | lte | eq | neq`
  * @param (policy.service_protection_core.adaptive_load_scheduler: aperture.spec.v1.AdaptiveLoadSchedulerParameters) Parameters for Adaptive Load Scheduler.
  * @param (policy.service_protection_core.dry_run: bool) Default configuration for setting dry run mode on Load Scheduler. In dry run mode, the Load Scheduler acts as a passthrough and does not throttle flows. This config can be updated at runtime without restarting the policy.
  */
  policy: {
    policy_name: '__REQUIRED_FIELD__',
    components: [],
    resources: {
      flow_control: {
        classifiers: [],
      },
    },
    evaluation_interval: '10s',
    service_protection_core: service_protection_core_defaults,
  },

  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) Time from of dashboard.
  * @param (dashboard.time_to: string) Time to of dashboard.
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  */
  dashboard: {
    refresh_interval: '15s',
    time_from: 'now-15m',
    time_to: 'now',
    extra_filters: {},
    title: 'Aperture Service Protection',
    /**
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
    variant_name: 'Service Protection',
  },
  selectors: selectors_defaults,
}
