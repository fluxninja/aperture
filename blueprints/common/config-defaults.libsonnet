local selectors_defaults = [{
  service: '__REQUIRED_FIELD__',
  control_point: '__REQUIRED_FIELD__',
}];

local kubeletstats_infra_meter_label_filter = {
  key: '__REQUIRED_FIELD__',
  value: '__REQUIRED_FIELD__',
  op: '__REQUIRED_FIELD__',
};

local kubeletstats_infra_meter_filter = {
  node: '',
  node_from_env_var: '',
  namespace: '',
  fields: [],
  labels: [],
};

{
  /**
  * @param (policy.policy_name: string) Name of the policy.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  * @param (policy.resources: aperture.spec.v1.Resources) Additional resources.
  */
  policy: {
    policy_name: '__REQUIRED_FIELD__',
    components: [],
    resources: {
      flow_control: {
        classifiers: [],
      },
    },
  },

  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) From time of dashboard.
  * @param (dashboard.time_to: string) To time of dashboard.
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  */
  dashboard: {
    refresh_interval: '15s',
    time_from: 'now-15m',
    time_to: 'now',
    extra_filters: {},
    title: '',
    /**
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
    variant_name: '',
  },

  // defaults for the schemas
  selectors_defaults: selectors_defaults,
  kubeletstats_infra_meter_label_filter: kubeletstats_infra_meter_label_filter,
  kubeletstats_infra_meter_filter: kubeletstats_infra_meter_filter,
}
