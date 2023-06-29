local commonConfig = import '../../common/config-defaults.libsonnet';

/**
* @schema (kubeletstats_infra_meter_label_filter.key: string) Key represents the key or name of the field or labels that a filter can apply on.
* @schema (kubeletstats_infra_meter_label_filter.value: string) Value represents the value associated with the key that a filter operation specified by the `Op` field applies on.
* @schema (kubeletstats_infra_meter_label_filter.op: string) Op represents the filter operation to apply on the given Key: Value pair. The supported operations are: equals, not-equals, exists, does-not-exist.
* @schema (kubeletstats_infra_meter_filter.node: string) Node represents a k8s node or host. If specified, any pods not running on the specified node will be ignored by the tagger.
* @schema (kubeletstats_infra_meter_filter.node_from_env_var: string) odeFromEnv can be used to extract the node name from an environment variable. For example: `NODE_NAME`.
* @schema (kubeletstats_infra_meter_filter.namespace: string) Namespace filters all pods by the provided namespace. All other pods are ignored.
* @schema (kubeletstats_infra_meter_filter.fields: []kubeletstats_infra_meter_label_filter) Fields allows to filter pods by generic k8s fields. Supported operations are: equals, not-equals.
* @schema (kubeletstats_infra_meter_filter.labels: []kubeletstats_infra_meter_label_filter) Labels allows to filter pods by generic k8s pod labels.
* @schema (kubeletstats_infra_meter.enabled: bool) Adds infra_meter for scraping Kubelet metrics.
* @schema (kubeletstats_infra_meter.agent_group: string) Agent group to be used for the infra_meter.
* @schema (kubeletstats_infra_meter.filter: kubeletstats_infra_meter_filter) Filter to be applied to the infra_meter.
*/
local kubeletstats_infra_meter = {
  enabled: false,
  agent_group: 'default',
  filter: {},
};

local service_protection_core_defaults = {
  overload_confirmations: [],

  adaptive_load_scheduler: {
    load_scheduler: {
      selectors: commonConfig.selectors_defaults,
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

commonConfig {
  /**
  * @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
  * @param (policy.service_protection_core.overload_confirmations: []overload_confirmation) List of overload confirmation criteria. Load scheduler can throttle flows when all of the specified overload confirmation criteria are met.
  * @schema (overload_confirmation.query_string: string) The Prometheus query to be run. Must return a scalar or a vector with a single element.
  * @schema (overload_confirmation.threshold: float64) The threshold for the overload confirmation criteria.
  * @schema (overload_confirmation.operator: string) The operator for the overload confirmation criteria. oneof: `gt | lt | gte | lte | eq | neq`
  * @param (policy.service_protection_core.adaptive_load_scheduler: aperture.spec.v1.AdaptiveLoadSchedulerParameters) Parameters for Adaptive Load Scheduler.
  * @param (policy.service_protection_core.dry_run: bool) Default configuration for setting dry run mode on Load Scheduler. In dry run mode, the Load Scheduler acts as a passthrough and does not throttle flows. This config can be updated at runtime without restarting the policy.
  * @param (policy.kubeletstats_infra_meter: kubeletstats_infra_meter) Infra meter for scraping Kubelet metrics.
  */
  policy+: {
    evaluation_interval: '10s',
    service_protection_core: service_protection_core_defaults,
    kubeletstats_infra_meter: kubeletstats_infra_meter,
  },

  dashboard+: {
    title: 'Aperture Service Protection',
    variant_name: 'Service Protection',
  },

  kubeletstats_infra_meter: kubeletstats_infra_meter,
}
