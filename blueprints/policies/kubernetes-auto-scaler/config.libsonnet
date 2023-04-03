{
  /**
  * @section Common
  *
  * @param (common.policy_name: string required) Name of the policy.
  */
  common: {
    policy_name: '__REQUIRED_FIELD__',
  },
  /**
  * @section Policy
  *
  * @param (policy.min_replicas: number) Minimum number of replicas.
  * @param (policy.max_replicas: number) Maximum number of replicas.
  * @param (policy.scale_in_cooldown: string) Scale in cooldown.
  * @param (policy.scale_out_cooldown: string) Scale out cooldown.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  */
  policy: {
    /**
    * @section Policy
    * @subsection Kubernetes Object Selector
    *
    * @param (policy.kubernetes_object_selector.namespace: string required) Namespace.
    * @param (policy.kubernetes_object_selector.api_version: string required) API Version.
    * @param (policy.kubernetes_object_selector.kind: string required) Kind.
    * @param (policy.kubernetes_object_selector.name: string required) Name.
    */
    kubernetes_object_selector: {
      namespace: '__REQUIRED_FIELD__',
      api_version: '__REQUIRED_FIELD__',
      kind: '__REQUIRED_FIELD__',
      name: '__REQUIRED_FIELD__',
    },
    min_replicas: 1,
    max_replicas: 10,
    scale_in_cooldown: '40s',
    scale_out_cooldown: '30s',
    components: [],
    /**
    * @section Policy
    * @subsection Scale-in criteria
    *
    * @param (policy.scale_in_criteria: []object) List of scale-in criteria.
    * @param (policy.scale_in_criteria.query: aperture.spec.v1.Query required) Query.
    * @param (policy.scale_in_criteria.query.promql: aperture.spec.v1.PromQL required) PromQL query.
    * @param (policy.scale_in_criteria.query.promql.query_string: string required) PromQL query string.
    * @param (policy.scale_in_criteria.query.promql.evaluation_interval: string) Evaluation interval.
    * @param (policy.scale_in_criteria.query.promql.out_ports.output.signal_name: string required) Output Signal name.
    * @param (policy.scale_in_criteria.set_point: number) Set point.
    * @param (policy.scale_in_criteria.parameters: aperture.spec.v1.GradientControllerParameters) Parameters.
    * @param (policy.scale_in_criteria.parameters.slope: number) Slope.
    */
    scale_in_criteria: [
      {
        query: {
          promql: {
            query_string: '__REQUIRED_FIELD__',
            evaluation_interval: '10s',
            out_ports: {
              output: {
                signal_name: '__REQUIRED_FIELD__',
              },
            },
          },
        },
        set_point: 0.5,
        parameters: {
          slope: 1.0,
        },
      },
    ],
    /**
    * @section Policy
    * @subsection Scale-out criteria
    *
    * @param (policy.scale_out_criteria: []object) List of scale-out criteria.
    * @param (policy.scale_out_criteria.query: aperture.spec.v1.Query required) Query.
    * @param (policy.scale_out_criteria.query.promql: aperture.spec.v1.PromQL required) PromQL query.
    * @param (policy.scale_out_criteria.query.promql.query_string: string required) PromQL query string.
    * @param (policy.scale_out_criteria.query.promql.evaluation_interval: string) Evaluation interval.
    * @param (policy.scale_out_criteria.query.promql.out_ports.output.signal_name: string required) Output Signal name.
    * @param (policy.scale_out_criteria.set_point: number) Set point.
    * @param (policy.scale_out_criteria.parameters: aperture.spec.v1.GradientControllerParameters) Parameters.
    * @param (policy.scale_out_criteria.parameters.slope: number) Slope.
    */
    scale_out_criteria: [
      {
        query: {
          promql: {
            query_string: '__REQUIRED_FIELD__',
            evaluation_interval: '10s',
            out_ports: {
              output: {
                signal_name: '__REQUIRED_FIELD__',
              },
            },
          },
        },
        set_point: 1.0,
        parameters: {
          slope: -1.0,
        },
      },
    ],
  },
  /**
  * @section Dashboard
  *
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) From time of dashboard.
  * @param (dashboard.time_to: string) To time of dashboard.
  */
  dashboard: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    /**
    * @section Dashboard
    * @subsection Datasource
    *
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
  },
}
