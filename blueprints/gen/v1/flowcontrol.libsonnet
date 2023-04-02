{
  new():: {
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
