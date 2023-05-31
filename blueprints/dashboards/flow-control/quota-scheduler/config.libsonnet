{
  /**
  * @param (policy.policy_name: string required) Name of the policy.
  */
  policy: {
    policy_name: '__REQUIRED_FIELD__',
  },
  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) From time of dashboard.
  * @param (dashboard.time_to: string) To time of dashboard.
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  * @param (dashboard.variant_name: string) Name of the dashboard variant.
  */
  dashboard: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    extra_filters: {},
    title: 'Aperture Quota Scheduler',
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
}
