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
  withHorizontalPodScaler(horizontal_pod_scaler):: {
    horizontal_pod_scaler: horizontal_pod_scaler,
  },
  withHorizontalPodScalerMixin(horizontal_pod_scaler):: {
    horizontal_pod_scaler+: horizontal_pod_scaler,
  },
  withPromql(promql):: {
    promql: promql,
  },
  withPromqlMixin(promql):: {
    promql+: promql,
  },
  withRateLimiter(rate_limiter):: {
    rate_limiter: rate_limiter,
  },
  withRateLimiterMixin(rate_limiter):: {
    rate_limiter+: rate_limiter,
  },
}
