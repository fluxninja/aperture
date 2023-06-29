local selectors_defaults = [{
  service: '__REQUIRED_FIELD__',
  control_point: '__REQUIRED_FIELD__',
}];

/**
* @param (policy.drivers.average_latency_drivers: []average_latency_driver) List of drivers that compare average latency against forward, backward and reset thresholds.
* @schema (average_latency_driver.selectors: []aperture.spec.v1.Selector) Identify the service and flows whose latency needs to be measured.
* @schema (average_latency_driver.criteria.forward.threshold: float64) The threshold for the forward criteria.
* @schema (average_latency_driver.criteria.backward.threshold: float64) The threshold for the backward criteria.
* @schema (average_latency_driver.criteria.reset.threshold: float64) The threshold for the reset criteria.
*/
local average_latency_driver_defaults = {
  selectors: selectors_defaults,
  criteria: {
    forward: {
      threshold: '__REQUIRED_FIELD__',
    },
  },
};

/**
* @param (policy.drivers.percentile_latency_drivers: []percentile_latency_driver) List of drivers that compare percentile latency against forward, backward and reset thresholds.
* @schema (percentile_latency_driver.flux_meter: aperture.spec.v1.FluxMeter) FluxMeter specifies the flows whose latency needs to be measured and parameters for the histogram metrics.
* @schema (percentile_latency_driver.percentile: float64) The percentile to be used for latency measurement.
* @schema (percentile_latency_driver.criteria.forward.threshold: float64) The threshold for the forward criteria.
* @schema (percentile_latency_driver.criteria.backward.threshold: float64) The threshold for the backward criteria.
* @schema (percentile_latency_driver.criteria.reset.threshold: float64) The threshold for the reset criteria.
*/
local percentile_latency_driver_defaults = {
  flux_meter: {
    selector: selectors_defaults,
    static_buckets: {
      buckets: [5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2500.0, 5000.0, 10000.0],
    },
  },
  percentile: 95,
  criteria: {
    forward: {
      threshold: '__REQUIRED_FIELD__',
    },
  },
};

/**
* @param (policy.drivers.promql_drivers: []promql_driver) List of promql drivers that compare results of a Prometheus query against forward, backward and reset thresholds.
* @schema (promql_driver.query_string: string) The Prometheus query to be run. Must return a scalar or a vector with a single element.
* @schema (promql_driver.criteria.forward.threshold: float64) The threshold for the forward criteria.
* @schema (promql_driver.criteria.forward.operator: string) The operator for the forward criteria. oneof: `gt | lt | gte | lte | eq | neq`
* @schema (promql_driver.criteria.backward.threshold: float64) The threshold for the backward criteria.
* @schema (promql_driver.criteria.backward.operator: string) The operator for the backward criteria. oneof: `gt | lt | gte | lte | eq | neq`
* @schema (promql_driver.criteria.reset.threshold: float64) The threshold for the reset criteria.
* @schema (promql_driver.criteria.reset.operator: string) The operator for the reset criteria. oneof: `gt | lt | gte | lte | eq | neq`
*/
local promql_driver_defaults = {
  query_string: '__REQUIRED_FIELD__',
  criteria: {
    forward: {
      threshold: '__REQUIRED_FIELD__',
      operator: '__REQUIRED_FIELD__',
    },
  },
};


local ramp_policy_base_defaults = {
  /**
  * @param (policy.policy_name: string) Name of the policy.
  * @param (policy.load_ramp: aperture.spec.v1.LoadRampParameters) Identify the service and flows of the feature that needs to be rolled out. And specify load ramp steps.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  * @param (policy.resources: aperture.spec.v1.Resources) List of additional resources.
  * @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
  * @param (policy.start: bool) Whether to start the ramp. This setting may be overridden at runtime via dynamic configuration.
  */
  policy_name: '__REQUIRED_FIELD__',
  load_ramp: {
    sampler: {
      selectors: selectors_defaults,
      label_key: '',
    },
    steps: [
      {
        target_accept_percentage: '__REQUIRED_FIELD__',
        duration: '__REQUIRED_FIELD__',
      },
    ],
  },

  drivers: {
    promql_drivers: [],
    average_latency_drivers: [],
    percentile_latency_drivers: [],
  },

  start: false,

  components: [],

  resources: {
    flow_control: {
      classifiers: [],
    },
  },

  evaluation_interval: '10s',
};

{
  policy: ramp_policy_base_defaults,

  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) From time of dashboard.
  * @param (dashboard.time_to: string) To time of dashboard.
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  */
  dashboard: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    extra_filters: {},
    title: 'Aperture Load Ramp',
    /**
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
  },

  // defaults for the schemas
  promql_driver: promql_driver_defaults,
  average_latency_driver: average_latency_driver_defaults,
  percentile_latency_driver: percentile_latency_driver_defaults,
}
