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
    * @param (policy.quota_scheduler.bucket_capacity: float64 required) Bucket capacity.
    * @param (policy.quota_scheduler.fill_amount: float64 required) Fill amount.
    * @param (policy.quota_scheduler.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against
    * @param (policy.quota_scheduler.rate_limiter: aperture.spec.v1.RateLimiterParameters) Rate Limiter Parameters.
    * @param (policy.quota_scheduler.scheduler: aperture.spec.v1.Scheduler) Scheduler configuration.
    */
    quota_scheduler: {
      bucket_capacity: '__REQUIRED_FIELD__',
      fill_amount: '__REQUIRED_FIELD__',
      selectors: [{
        service: '__REQUIRED_FIELD__',
        control_point: '__REQUIRED_FIELD__',
      }],
      rate_limiter: {
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
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  */
  dashboard: {
    refresh_interval: '10s',
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
  },
}
