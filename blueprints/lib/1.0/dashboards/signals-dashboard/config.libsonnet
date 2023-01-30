{
  /**
  * @section Common
  *
  * @param (common.policy_name: string required) Name of the policy.
  */
  common: {
    policy_name: error 'policy_name must be set',
  },
  /**
  * @section Dashboard
  *
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) From time of dashboard.
  * @param (dashboard.time_to: string) To time of dashboard.
  */
  dashboard: {
    refresh_interval: '10s',
    time_from: 'now-30m',
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
