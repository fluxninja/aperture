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
};

commonConfig {
  /**
  * @param (policy.promql_scale_out_controllers: []promql_scale_out_controller) List of scale out controllers.
  * @param (policy.promql_scale_in_controllers: []promql_scale_in_controller) List of scale in controllers.
  * @param (policy.scaling_parameters: aperture.spec.v1.AutoScalerScalingParameters) Parameters that define the scaling behavior.
  * @param (policy.scaling_backend: aperture.spec.v1.AutoScalerScalingBackend) Scaling backend for the policy.
  * @param (policy.dry_run: bool) Dry run mode ensures that no scaling is invoked by this auto scaler.
  */
  policy+: auto_scaling_defaults,

  // schema defaults are below
  promql_scale_out_controller: promql_scale_controller_defaults,
  promql_scale_in_controller: promql_scale_controller_defaults,
}
