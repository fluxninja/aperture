local autoScalingDefaults = import '../base/config-defaults.libsonnet';

/**
* @schema (promql_scale_controller.query_string: string) The Prometheus query to be run. Must return a scalar or a vector with a single element.
* @schema (promql_scale_controller.setpoint: float64) Setpoint for the controller.
* @schema (promql_scale_controller.gradient: aperture.spec.v1.IncreasingGradientParameters) Gradient parameters for the controller.
* @schema (promql_scale_controller.alerter: aperture.spec.v1.AlerterParameters) Alerter parameters for the controller.
*/

autoScalingDefaults {
  /**
  * @param (policy.policy_name: string) Name of the policy.
  * @param (policy.promql_scale_out_controllers: []promql_scale_controller) List of scale out controllers.
  * @param (policy.promql_scale_in_controllers: []promql_scale_controller) List of scale in controllers.
  * @param (policy.scaling_parameters: aperture.spec.v1.AutoScalerScalingParameters) Parameters that define the scaling behavior.
  * @param (policy.scaling_backend: aperture.spec.v1.AutoScalerScalingBackend) Scaling backend for the policy.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  * @param (policy.resources: aperture.spec.v1.Resources) List of additional resources.
  * @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
  * @param (policy.dry_run: bool) Dry run mode ensures that no scaling is invoked by this auto scaler.
  */
  policy+: {
    scaling_backend: {
      kubernetes_replicas: '__REQUIRED_FIELD__',
    },
  },

  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) Time from of dashboard.
  * @param (dashboard.time_to: string) Time to of dashboard.
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  */
  dashboard+: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    extra_filters: {},
    title: 'Aperture Auto-scale',
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
