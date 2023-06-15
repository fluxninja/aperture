{
  new():: {
  },
  withAdaptiveLoadScheduler(adaptive_load_scheduler):: {
    adaptive_load_scheduler: adaptive_load_scheduler,
  },
  withAdaptiveLoadSchedulerMixin(adaptive_load_scheduler):: {
    adaptive_load_scheduler+: adaptive_load_scheduler,
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
  withQuotaScheduler(quota_scheduler):: {
    quota_scheduler: quota_scheduler,
  },
  withQuotaSchedulerMixin(quota_scheduler):: {
    quota_scheduler+: quota_scheduler,
  },
  withRateLimiter(rate_limiter):: {
    rate_limiter: rate_limiter,
  },
  withRateLimiterMixin(rate_limiter):: {
    rate_limiter+: rate_limiter,
  },
  withSampler(sampler):: {
    sampler: sampler,
  },
  withSamplerMixin(sampler):: {
    sampler+: sampler,
  },
}
