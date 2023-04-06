local scale_criteria = {
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
};

{
  /**
  * @param (common.policy_name: string required) Name of the policy.
  */
  common: {
    policy_name: '__REQUIRED_FIELD__',
  },
  /**
  * @param (policy.min_replicas: string) Minimum number of replicas.
  * @param (policy.max_replicas: string) Maximum number of replicas.
  * @param (policy.scale_in_cooldown: string) The amount of time to wait after a scale-in operation for another scale-in operation.
  * @param (policy.scale_out_cooldown: string) The amount of time to wait after a scale-out operation for another scale-out or scale-in operation.
  * @param (policy.cooldown_override_percentage: float64) Cooldown override percentage defines a threshold change in scale-out beyond which previous cooldown is overridden.
  * @param (policy.max_scale_in_percentage: float64) The maximum decrease of replicas (e.g. pods) at one time.
  * @param (policy.max_scale_out_percentage: float64) The maximum increase of replicas (e.g. pods) at one time.
  * @param (policy.scale_in_alerter_parameters: aperture.spec.v1.AlerterParameters) Configuration for scale-in alerter.
  * @param (policy.scale_out_alerter_parameters: aperture.spec.v1.AlerterParameters) Cooldown override percentage.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  */
  policy: {
    /**
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
    min_replicas: '1',
    max_replicas: '10',
    scale_in_cooldown: '40s',
    scale_out_cooldown: '30s',
    cooldown_override_percentage: 50,
    max_scale_in_percentage: 1,
    max_scale_out_percentage: 10,
    scale_in_alerter_parameters: {
      alert_name: 'Kubernetes Auto Scaler Scale In Event',
    },
    scale_out_alerter_parameters: {
      alert_name: 'Kubernetes Auto Scaler Scale Out Event',
    },
    components: [],
    /**
    * @param (policy.scale_in_criteria: []scale_criteria required) List of scale-in criteria.
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
    * @param (policy.scale_out_criteria: []scale_criteria required) List of scale-out criteria.
    */
    scale_out_criteria: [
      scale_criteria,
    ],
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
  /**
  * @schema (scale_criteria.query: aperture.spec.v1.Query required) Query.
  * @schema (scale_criteria.query.promql: aperture.spec.v1.PromQL required) PromQL query.
  * @schema (scale_criteria.query.promql.query_string: string required) PromQL query string.
  * @schema (scale_criteria.query.promql.evaluation_interval: string) Evaluation interval.
  * @schema (scale_criteria.query.promql.out_ports: aperture.spec.v1.PromQLOuts required) PromQL query execution output.
  * @schema (scale_criteria.query.promql.out_ports.output: aperture.spec.v1.OutPort required) PromQL query execution output port.
  * @schema (scale_criteria.query.promql.out_ports.output.signal_name: string required) Output Signal name.
  * @schema (scale_criteria.set_point: float64) Set point.
  * @schema (scale_criteria.parameters: aperture.spec.v1.IncreasingGradientParameters) Parameters.
  * @schema (scale_criteria.parameters.slope: float64) Slope.
  */
  scale_criteria: scale_criteria,
}
