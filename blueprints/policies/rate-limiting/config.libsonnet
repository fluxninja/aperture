{
  policy: {
    /**
    * @param (policy.policy_name: string required) Name of the policy.
    * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
    * @param (policy.resources: aperture.spec.v1.Resources) Additional resources.
    */
    policy_name: '__REQUIRED_FIELD__',
    components: [],
    resources: {
      flow_control: {
        classifiers: [],
      },
    },
    /**
    * @param (policy.rate_limiter.bucket_capacity: float64 required) Bucket capacity.
    * @param (policy.rate_limiter.fill_amount: float64 required) Fill amount.
    * @param (policy.rate_limiter.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against
    * @param (policy.rate_limiter.parameters: aperture.spec.v1.RateLimiterParameters) Parameters.
    * @param (policy.rate_limiter.parameters.label_key: string) Flow label key to use for rate limiting. If not specified, then all requests matching the selector will be rate limited.
    * @param (policy.rate_limiter.parameters.interval: string required) Fill interval e.g. "1s".
    */
    rate_limiter: {
      bucket_capacity: '__REQUIRED_FIELD__',
      fill_amount: '__REQUIRED_FIELD__',
      selectors: [{
        service: '__REQUIRED_FIELD__',
        control_point: '__REQUIRED_FIELD__',
      }],
      parameters: {
        label_key: '',
        interval: '__REQUIRED_FIELD__',
      },
    },
  },
  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  */
  dashboard: {
    refresh_interval: '10s',
    extra_filters: {},
    title: 'Aperture Rate Limiter',
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
