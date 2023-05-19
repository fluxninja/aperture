/**
* @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
* @param (policy.resources: aperture.spec.v1.Resources) Additional resources.
*/
{
  policy: {
    /**
    * @param (policy.policy_name: string required) Name of the policy.
    */
    policy_name: '__REQUIRED_FIELD__',
    components: [],
    resources: {
      flow_control: {
        classifiers: [],
      },
    },
    /**
    * @param (policy.rate_limit: float64 required) Number of requests per `policy.rate_limiter.limit_reset_interval` to accept
    */
    rate_limit: '__REQUIRED_FIELD__',
    /**
    * @param (policy.rate_limiter: aperture.spec.v1.RateLimiterParameters) Parameters for _Rate Limiter_.
    */
    rate_limiter: {
      selectors: [{
        service: '__REQUIRED_FIELD__',
        control_point: '__REQUIRED_FIELD__',
      }],
      limit_reset_interval: '__REQUIRED_FIELD__',
      label_key: '__REQUIRED_FIELD__',
    },
    /**
    * @param (policy.custom_limits: []aperture.spec.v1.RateLimiterCustomLimit) Allows to specify different limits for particular label values. This setting can be updated at runtime without restarting the policy via dynamic config.
    */
    custom_limits: [],
  },
  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  */
  dashboard: {
    refresh_interval: '10s',
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
