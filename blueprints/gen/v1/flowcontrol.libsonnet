{
  new():: {
  },
  withAdaptiveLoadScheduler(adaptive_load_scheduler):: {
    adaptive_load_scheduler: adaptive_load_scheduler,
  },
  withAdaptiveLoadSchedulerMixin(adaptive_load_scheduler):: {
    adaptive_load_scheduler+: adaptive_load_scheduler,
  },
  withAimdConcurrencyController(aimd_concurrency_controller):: {
    aimd_concurrency_controller: aimd_concurrency_controller,
  },
  withAimdConcurrencyControllerMixin(aimd_concurrency_controller):: {
    aimd_concurrency_controller+: aimd_concurrency_controller,
  },
  withConcurrencyLimiter(concurrency_limiter):: {
    concurrency_limiter: concurrency_limiter,
  },
  withConcurrencyLimiterMixin(concurrency_limiter):: {
    concurrency_limiter+: concurrency_limiter,
  },
  withFlowRegulator(flow_regulator):: {
    flow_regulator: flow_regulator,
  },
  withFlowRegulatorMixin(flow_regulator):: {
    flow_regulator+: flow_regulator,
  },
  withLoadRegulator(load_regulator):: {
    load_regulator: load_regulator,
  },
  withLoadRegulatorMixin(load_regulator):: {
    load_regulator+: load_regulator,
  },
  withLoadScheduler(load_scheduler):: {
    load_scheduler: load_scheduler,
  },
  withLoadSchedulerMixin(load_scheduler):: {
    load_scheduler+: load_scheduler,
  },
  withLoadShaper(load_shaper):: {
    load_shaper: load_shaper,
  },
  withLoadShaperMixin(load_shaper):: {
    load_shaper+: load_shaper,
  },
  withLoadShaperSeries(load_shaper_series):: {
    load_shaper_series: load_shaper_series,
  },
  withLoadShaperSeriesMixin(load_shaper_series):: {
    load_shaper_series+: load_shaper_series,
  },
  withRateLimiter(rate_limiter):: {
    rate_limiter: rate_limiter,
  },
  withRateLimiterMixin(rate_limiter):: {
    rate_limiter+: rate_limiter,
  },
}
