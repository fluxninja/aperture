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
    * @param (policy.quota_scheduler.bucket_capacity: float64 required) Bucket capacity.
    * @param (policy.quota_scheduler.fill_amount: float64 required) Fill amount.
    * @param (policy.quota_scheduler.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against
    * @param (policy.quota_scheduler.parameters: aperture.spec.v1.RateLimiterParameters) Parameters.
    * @param (policy.quota_scheduler.parameters.label_key: string) Flow label key to use for quota limiting. If no label key is specified, then all requests matching the selectors will be scheduled based on the global quota.
    * @param (policy.quota_scheduler.parameters.interval: string required) Fill interval e.g. "1s".
    * @param (policy.quota_scheduler.scheduler: aperture.spec.v1.Scheduler) Scheduler configuration.
    */
    quota_scheduler: {
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
      scheduler: {},
    },
  },
  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) Time from of dashboard.
  * @param (dashboard.time_to: string) Time to of dashboard.
  */
  dashboard: {
    refresh_interval: '10s',
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
