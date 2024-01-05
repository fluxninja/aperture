local commonConfig = import '../../common/config-defaults.libsonnet';

commonConfig {
  policy+: {
    /**
    * @param (policy.concurrency_limiter.max_concurrency: float64) Max concurrency.
    * @param (policy.concurrency_limiter.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against
    * @param (policy.concurrency_limiter.parameters: aperture.spec.v1.ConcurrencyLimiterParameters) Parameters.
    * @param (policy.concurrency_limiter.request_parameters: aperture.spec.v1.ConcurrencyLimiterRequestParameters) Request Parameters.
    * @param (policy.concurrency_limiter.alerter: aperture.spec.v1.AlerterParameters) Alerter.
    */
    concurrency_limiter: {
      max_concurrency: '__REQUIRED_FIELD__',
      selectors: commonConfig.selectors_defaults,
      parameters: {
        max_inflight_duration: '__REQUIRED_FIELD__',
      },
      request_parameters: {},
      alerter: {
        alert_name: 'Too many inflight requests',
      },
    },
  },
}
