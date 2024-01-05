local commonConfig = import '../../common/config-defaults.libsonnet';

commonConfig {
  policy+: {
    /**
    * @param (policy.quota_scheduler.bucket_capacity: float64) Bucket capacity.
    * @param (policy.quota_scheduler.fill_amount: float64) Fill amount.
    * @param (policy.quota_scheduler.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against.
    * @param (policy.quota_scheduler.rate_limiter: aperture.spec.v1.RateLimiterParameters) Rate Limiter Parameters.
    * @param (policy.quota_scheduler.scheduler: aperture.spec.v1.Scheduler) Scheduler configuration.
    * @param (policy.quota_scheduler.alerter: aperture.spec.v1.AlerterParameters) Alerter.
    */
    quota_scheduler: {
      bucket_capacity: '__REQUIRED_FIELD__',
      fill_amount: '__REQUIRED_FIELD__',
      selectors: commonConfig.selectors_defaults,
      rate_limiter: {
        interval: '__REQUIRED_FIELD__',
      },
      scheduler: {},
      alerter: {
        alert_name: 'More than 90% of requests are being rate limited',
      },
    },
  },
}
