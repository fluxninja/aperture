local commonConfig = import '../../common/config-defaults.libsonnet';

commonConfig {
  policy+: {
    /**
    * @param (policy.rate_limiter.bucket_capacity: float64) Bucket capacity.
    * @param (policy.rate_limiter.fill_amount: float64) Fill amount.
    * @param (policy.rate_limiter.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against
    * @param (policy.rate_limiter.parameters: aperture.spec.v1.RateLimiterParameters) Parameters.
    * @param (policy.rate_limiter.parameters.label_key: string) Flow label key to use for rate limiting. If not specified, then all requests matching the selector will be rate limited.
    * @param (policy.rate_limiter.parameters.interval: string) Fill interval e.g. "1s".
    */
    rate_limiter: {
      bucket_capacity: '__REQUIRED_FIELD__',
      fill_amount: '__REQUIRED_FIELD__',
      selectors: commonConfig.selectors_defaults,
      parameters: {
        label_key: '',
        interval: '__REQUIRED_FIELD__',
      },
    },
  },

  dashboard+: {
    title: 'Aperture Rate Limiter',
  },
}
