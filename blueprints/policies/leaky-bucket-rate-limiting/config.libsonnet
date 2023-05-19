{
  policy: {
    /**
    * @param (policy.policy_name: string required) Name of the policy.
    */
    policy_name: '__REQUIRED_FIELD__',
    /**
    * @param (policy.classifiers: []aperture.spec.v1.Classifier) List of classification rules.
    */
    classifiers: [],
    /**
    * @param (policy.rate_limiter.bucket_capacity: float64 required) Bucket capacity.
    * @param (policy.rate_limiter.leak_amount: float64 required) Leak amount.
    * @param (policy.rate_limiter.leak_interval_ms: int64 required) Leak interval in milliseconds.
    * @param (policy.rate_limiter.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against
    * @param (policy.rate_limiter.parameters: aperture.spec.v1.RateLimiterParameters) Parameters.
    * @param (policy.rate_limiter.parameters.label_key: string required) Flow label to use for rate limiting.
    */
    rate_limiter: {
      bucket_capacity: '__REQUIRED_FIELD__',
      leak_amount: '__REQUIRED_FIELD__',
      leak_interval_ms: '__REQUIRED_FIELD__',
      selectors: [{
        service: '__REQUIRED_FIELD__',
        control_point: '__REQUIRED_FIELD__',
      }],
      parameters: {
        label_key: '__REQUIRED_FIELD__',
      },
    },
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
