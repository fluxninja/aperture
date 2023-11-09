local commonConfig = import '../../common/config-defaults.libsonnet';

/**
* @schema (criteria.threshold: float64) The threshold for the criteria.
*/
local criteria_defaults = {
  threshold: '__REQUIRED_FIELD__',
};

/**
* @schema (promql_criteria.threshold: float64) The threshold for the criteria.
* @schema (promql_criteria.operator: string) The operator for the criteria. oneof: `gt | lt | gte | lte | eq | neq`.
*/
local promql_criteria_defaults = criteria_defaults {
  operator: '__REQUIRED_FIELD__',
};

/**
* @schema (driver_criteria.forward: criteria) The forward criteria.
* @schema (driver_criteria.backward: criteria) The backward criteria.
* @schema (driver_criteria.reset: criteria) The reset criteria.
* @schema (promql_driver_criteria.forward: promql_criteria) The forward criteria.
* @schema (promql_driver_criteria.backward: promql_criteria) The backward criteria.
* @schema (promql_driver_criteria.reset: promql_criteria) The reset criteria.
*/
local driver_criteria_defaults = {
  forward: {},
  backward: {},
  reset: {},
};

/**
* @schema (kubelet_metrics_criteria.pod_cpu: driver_criteria) The criteria of the pod cpu usage driver.
* @schema (kubelet_metrics_criteria.pod_memory: driver_criteria) The criteria of the pod memory usage driver.
*/
local kubelet_metrics_criteria_defaults = {
  pod_cpu: {},
  pod_memory: {},
};

/**
* @param (policy.drivers.average_latency_drivers: []average_latency_driver) List of drivers that compare average latency against forward, backward and reset thresholds.
* @schema (average_latency_driver.selectors: []aperture.spec.v1.Selector) Identify the service and flows whose latency needs to be measured.
* @schema (average_latency_driver.criteria: driver_criteria) The criteria of the driver.
*/
local average_latency_driver_defaults = {
  selectors: commonConfig.selectors_defaults,
  criteria: '__REQUIRED_FIELD__',
};

/**
* @param (policy.drivers.percentile_latency_drivers: []percentile_latency_driver) List of drivers that compare percentile latency against forward, backward and reset thresholds.
* @schema (percentile_latency_driver.flux_meter: aperture.spec.v1.FluxMeter) FluxMeter specifies the flows whose latency needs to be measured and parameters for the histogram metrics.
* @schema (percentile_latency_driver.percentile: float64) The percentile to be used for latency measurement.
* @schema (percentile_latency_driver.criteria: driver_criteria) The criteria of the driver.
*/
local percentile_latency_driver_defaults = {
  flux_meter: {
    selector: commonConfig.selectors_defaults,
    static_buckets: {
      buckets: [5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2500.0, 5000.0, 10000.0],
    },
  },
  percentile: 95,
  criteria: '__REQUIRED_FIELD__',
};

/**
* @param (policy.drivers.promql_drivers: []promql_driver) List of promql drivers that compare results of a Prometheus query against forward, backward and reset thresholds.
* @schema (promql_driver.query_string: string) The Prometheus query to be run. Must return a scalar or a vector with a single element.
* @schema (promql_driver.criteria: promql_driver_criteria) The criteria of the driver.
*/
local promql_driver_defaults = {
  query_string: '__REQUIRED_FIELD__',
  criteria: '__REQUIRED_FIELD__',
};

/**
* @param (policy.kubelet_metrics: kubelet_metrics) Kubelet metrics configuration.
* @schema (kubelet_metrics.infra_context: aperture.spec.v1.KubernetesObjectSelector) Kubernetes selector for scraping metrics.
* @schema (kubelet_metrics.criteria: kubelet_metrics_criteria) Criteria.
*/
local kubelet_metrics_defaults = {
  infra_context: '__REQUIRED_FIELD__',
  criteria: '__REQUIRED_FIELD__',
};

local ramp_policy_base_defaults = {
  /**
  * @param (policy.load_ramp: aperture.spec.v1.LoadRampParameters) Identify the service and flows of the feature that needs to be rolled out. And specify load ramp steps.
  * @param (policy.start: bool) Whether to start the ramp. This setting may be overridden at runtime via dynamic configuration.
  */
  load_ramp: {
    sampler: {
      selectors: commonConfig.selectors_defaults,
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

  kubelet_metrics: {},
};

commonConfig {
  policy+: ramp_policy_base_defaults,

  // defaults for the schemas
  promql_driver: promql_driver_defaults,
  average_latency_driver: average_latency_driver_defaults,
  percentile_latency_driver: percentile_latency_driver_defaults,
  kubelet_metrics: kubelet_metrics_defaults,
  criteria: criteria_defaults,
  promql_criteria: promql_criteria_defaults,
  driver_criteria: driver_criteria_defaults,
  promql_driver_criteria: driver_criteria_defaults,
  kubelet_metrics_criteria: kubelet_metrics_criteria_defaults,
}
