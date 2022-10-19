{
  /**
  * @section Common
  *
  * @param (common.policyName: string required) Name of the policy.
  */
  common: {
    policyName: error 'policyName must be set',
  },
  /**
  * @section Dashboard
  *
  * @param (dashboard.refreshInterval: string) Refresh interval for dashboard panels.
  * @param (dashboard.timeFrom: string) From time of dashboard.
  * @param (dashboard.timeTo: string) To time of dashboard.
  */
  dashboard: {
    refreshInterval: '10s',
    timeFrom: 'now-30m',
    timeTo: 'now',
    /**
    * @section Dashboard
    * @subsection Datasource
    *
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filterRegex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filterRegex: '',
    },
  },
}
