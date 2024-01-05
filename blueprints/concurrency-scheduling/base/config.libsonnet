local commonConfig = import '../../common/config-defaults.libsonnet';

commonConfig {
  policy+: {
    /**
    * @param (policy.concurrency_scheduler.max_concurrency: float64) Max concurrency.
    * @param (policy.concurrency_scheduler.selectors: []aperture.spec.v1.Selector) Flow selectors to match requests against.
    * @param (policy.concurrency_scheduler.concurrency_limiter: aperture.spec.v1.ConcurrencyLimiterParameters) Concurrency Limiter Parameters.
    * @param (policy.concurrency_scheduler.scheduler: aperture.spec.v1.Scheduler) Scheduler configuration.
    * @param (policy.concurrency_scheduler.alerter: aperture.spec.v1.AlerterParameters) Alerter.
    */
    concurrency_scheduler: {
      max_concurrency: '__REQUIRED_FIELD__',
      selectors: commonConfig.selectors_defaults,
      concurrency_limiter: {
        max_inflight_duration: '__REQUIRED_FIELD__',
      },
      scheduler: {},
      alerter: {
        alert_name: 'Too many inflight requests',
      },
    },
  },
}
