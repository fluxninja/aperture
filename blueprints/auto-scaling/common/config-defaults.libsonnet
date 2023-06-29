local commonConfig = import '../../common/config-defaults.libsonnet';

local scaling_parameters_defaults = {
  scale_in_alerter: {
    alert_name: 'Auto-scaler is scaling in',
  },
  scale_out_alerter: {
    alert_name: 'Auto-scaler is scaling out',
  },
};

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
  enabled: true,
  agent_group: 'default',
  filter: {},
};

/**
* @schema (promql_scale_out_controller.query_string: string) The Prometheus query to be run. Must return a scalar or a vector with a single element.
* @schema (promql_scale_out_controller.setpoint: float64) Setpoint for the controller.
* @schema (promql_scale_out_controller.gradient: aperture.spec.v1.IncreasingGradientParameters) Gradient parameters for the controller.
* @schema (promql_scale_out_controller.alerter: aperture.spec.v1.AlerterParameters) Alerter parameters for the controller.
* @schema (promql_scale_in_controller.query_string: string) The Prometheus query to be run. Must return a scalar or a vector with a single element.
* @schema (promql_scale_in_controller.setpoint: float64) Setpoint for the controller.
* @schema (promql_scale_in_controller.gradient: aperture.spec.v1.DecreasingGradientParameters) Gradient parameters for the controller.
* @schema (promql_scale_in_controller.alerter: aperture.spec.v1.AlerterParameters) Alerter parameters for the controller.
*/
local promql_scale_controller_defaults = {
  query_string: '__REQUIRED_FIELD__',
  setpoint: '__REQUIRED_FIELD__',
  gradient: '__REQUIRED_FIELD__',
  alerter: '__REQUIRED_FIELD__',
};

local auto_scaling_defaults = {
  promql_scale_out_controllers: [],

  promql_scale_in_controllers: [],

  scaling_parameters: scaling_parameters_defaults,

  scaling_backend: '__REQUIRED_FIELD__',

  dry_run: false,

  kubeletstats_infra_meter: kubeletstats_infra_meter,
};

commonConfig {
  /**
  * @param (policy.promql_scale_out_controllers: []promql_scale_out_controller) List of scale out controllers.
  * @param (policy.promql_scale_in_controllers: []promql_scale_in_controller) List of scale in controllers.
  * @param (policy.scaling_parameters: aperture.spec.v1.AutoScalerScalingParameters) Parameters that define the scaling behavior.
  * @param (policy.scaling_backend: aperture.spec.v1.AutoScalerScalingBackend) Scaling backend for the policy.
  * @param (policy.dry_run: bool) Dry run mode ensures that no scaling is invoked by this auto scaler.
  * @param (policy.kubeletstats_infra_meter: kubeletstats_infra_meter) Infra meter for scraping Kubelet metrics.
  */
  policy+: auto_scaling_defaults,

  dashboard+: {
    title: 'Aperture Auto-scale',
  },

  // schema defaults are below
  promql_scale_out_controller: promql_scale_controller_defaults,
  promql_scale_in_controller: promql_scale_controller_defaults,
  kubeletstats_infra_meter: kubeletstats_infra_meter,
}
