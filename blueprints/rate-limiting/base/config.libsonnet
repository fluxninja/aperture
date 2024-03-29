local commonConfig = import '../../common/config-defaults.libsonnet';

commonConfig {
  policy+: {
    /**
    * @param (policy.rate_limiter.bucket_capacity: float64) Bucket capacity.
    * @param (policy.rate_limiter.fill_amount: float64) Fill amount.
    * @param (policy.rate_limiter.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against
    * @param (policy.rate_limiter.parameters: aperture.spec.v1.RateLimiterParameters) Parameters.
    * @param (policy.rate_limiter.request_parameters: aperture.spec.v1.RateLimiterRequestParameters) Request Parameters.
    * @param (policy.rate_limiter.alerter: aperture.spec.v1.AlerterParameters) Alerter.
    */
    rate_limiter: {
      bucket_capacity: '__REQUIRED_FIELD__',
      fill_amount: '__REQUIRED_FIELD__',
      selectors: commonConfig.selectors_defaults,
      parameters: {
        interval: '__REQUIRED_FIELD__',
      },
      request_parameters: {},
      alerter: {
        alert_name: 'More than 90% of requests are being rate limited',
      },
    },
  },
}
