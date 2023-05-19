{
  new():: {
  },
  withAdaptiveLoadScheduler(adaptive_load_scheduler):: {
    adaptive_load_scheduler: adaptive_load_scheduler,
  },
  withAdaptiveLoadSchedulerMixin(adaptive_load_scheduler):: {
    adaptive_load_scheduler+: adaptive_load_scheduler,
  },
  withLeakyBucketRateLimiter(leaky_bucket_rate_limiter):: {
    leaky_bucket_rate_limiter: leaky_bucket_rate_limiter,
  },
  withLeakyBucketRateLimiterMixin(leaky_bucket_rate_limiter):: {
    leaky_bucket_rate_limiter+: leaky_bucket_rate_limiter,
  },
  withLoadRamp(load_ramp):: {
    load_ramp: load_ramp,
  },
  withLoadRampMixin(load_ramp):: {
    load_ramp+: load_ramp,
  },
  withLoadRampSeries(load_ramp_series):: {
    load_ramp_series: load_ramp_series,
  },
  withLoadRampSeriesMixin(load_ramp_series):: {
    load_ramp_series+: load_ramp_series,
  },
  withLoadScheduler(load_scheduler):: {
    load_scheduler: load_scheduler,
  },
  withLoadSchedulerMixin(load_scheduler):: {
    load_scheduler+: load_scheduler,
  },
  withRateLimiter(rate_limiter):: {
    rate_limiter: rate_limiter,
  },
  withRateLimiterMixin(rate_limiter):: {
    rate_limiter+: rate_limiter,
  },
  withRegulator(regulator):: {
    regulator: regulator,
  },
  withRegulatorMixin(regulator):: {
    regulator+: regulator,
  },
}
