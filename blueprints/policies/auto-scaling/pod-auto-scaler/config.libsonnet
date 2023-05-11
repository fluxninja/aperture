local autoScalingDefaults = import '../base/config-defaults.libsonnet';

/**
* @param (policy.policy_name: string required) Name of the policy.
* @schema (promql_scale_out_controller.query_string: string required) The Prometheus query to be run. Must return a scalar or a vector with a single element.
* @schema (promql_scale_out_controller.threshold: float64 required) Threshold for the controller.
* @schema (promql_scale_out_controller.gradient: aperture.spec.v1.IncreasingGradientParameters required) Gradient parameters for the controller.
* @param (policy.promql_scale_out_controllers: []promql_scale_out_controller required) List of scale out controllers.
* @schema (promql_scale_in_controller.query_string: string required) The Prometheus query to be run. Must return a scalar or a vector with a single element.
* @schema (promql_scale_in_controller.threshold: float64 required) Threshold for the controller.
* @schema (promql_scale_in_controller.gradient: aperture.spec.v1.DecreasingGradientParameters required) Gradient parameters for the controller.
* @param (policy.promql_scale_in_controllers: []promql_scale_in_controller required) List of scale in controllers.
* @param (policy.scaling_parameters: aperture.spec.v1.AutoScalerScalingParameters required) Parameters that define the scaling behavior.
* @param (policy.scaling_backend: aperture.spec.v1.AutoScalerScalingBackend required) Scaling backend for the policy.
* @param (policy.scaling_backend.kubernetes_replicas: aperture.spec.v1.AutoScalerScalingBackendKubernetesReplicas required) Kubernetes replicas scaling backend.
* @param (policy.scaling_backend.kubernetes_replicas.kubernetes_object_selector: aperture.spec.v1.KubernetesObjectSelector required) Kubernetes object selector.
* @param (policy.scaling_backend.kubernetes_replicas.min_replicas: string required) Minimum number of replicas.
* @param (policy.scaling_backend.kubernetes_replicas.max_replicas: string required) Maximum number of replicas.
* @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
* @param (policy.resources: aperture.spec.v1.Resources) List of additional resources.
* @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
* @param (policy.dry_run: bool) Dry run mode ensures that no scaling is invoked by this auto scaler.
* @param (policy.dry_run_config_key: string) Configuration key for overriding dry run setting through dynamic configuration.
*/
autoScalingDefaults {
  policy+: {
    scaling_backend: {
      kubernetes_replicas: {
        kubernetes_object_selector: '__REQUIRED_FIELD__',
        min_replicas: '__REQUIRED_FIELD__',
        max_replicas: '__REQUIRED_FIELD__',
      },
    },
  },

  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) Time from of dashboard.
  * @param (dashboard.time_to: string) Time to of dashboard.
  */
  dashboard+: {
    refresh_interval: '15s',
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
