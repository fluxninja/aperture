{
  new():: {
  },
  withAlerter(alerter):: {
    alerter: alerter,
  },
  withAlerterMixin(alerter):: {
    alerter+: alerter,
  },
  withArithmeticCombinator(arithmetic_combinator):: {
    arithmetic_combinator: arithmetic_combinator,
  },
  withArithmeticCombinatorMixin(arithmetic_combinator):: {
    arithmetic_combinator+: arithmetic_combinator,
  },
  withConcurrencyLimiter(concurrency_limiter):: {
    concurrency_limiter: concurrency_limiter,
  },
  withConcurrencyLimiterMixin(concurrency_limiter):: {
    concurrency_limiter+: concurrency_limiter,
  },
  withConstant(constant):: {
    constant: constant,
  },
  withConstantMixin(constant):: {
    constant+: constant,
  },
  withDecider(decider):: {
    decider: decider,
  },
  withDeciderMixin(decider):: {
    decider+: decider,
  },
  withDifferentiator(differentiator):: {
    differentiator: differentiator,
  },
  withDifferentiatorMixin(differentiator):: {
    differentiator+: differentiator,
  },
  withEma(ema):: {
    ema: ema,
  },
  withEmaMixin(ema):: {
    ema+: ema,
  },
  withExtrapolator(extrapolator):: {
    extrapolator: extrapolator,
  },
  withExtrapolatorMixin(extrapolator):: {
    extrapolator+: extrapolator,
  },
  withFirstValid(first_valid):: {
    first_valid: first_valid,
  },
  withFirstValidMixin(first_valid):: {
    first_valid+: first_valid,
  },
  withGradientController(gradient_controller):: {
    gradient_controller: gradient_controller,
  },
  withGradientControllerMixin(gradient_controller):: {
    gradient_controller+: gradient_controller,
  },
  withHorizontalPodScaler(horizontal_pod_scaler):: {
    horizontal_pod_scaler: horizontal_pod_scaler,
  },
  withHorizontalPodScalerMixin(horizontal_pod_scaler):: {
    horizontal_pod_scaler+: horizontal_pod_scaler,
  },
  withIntegrator(integrator):: {
    integrator: integrator,
  },
  withIntegratorMixin(integrator):: {
    integrator+: integrator,
  },
  withMax(max):: {
    max: max,
  },
  withMaxMixin(max):: {
    max+: max,
  },
  withMin(min):: {
    min: min,
  },
  withMinMixin(min):: {
    min+: min,
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
  withSink(sink):: {
    sink: sink,
  },
  withSinkMixin(sink):: {
    sink+: sink,
  },
  withSqrt(sqrt):: {
    sqrt: sqrt,
  },
  withSqrtMixin(sqrt):: {
    sqrt+: sqrt,
  },
  withSwitcher(switcher):: {
    switcher: switcher,
  },
  withSwitcherMixin(switcher):: {
    switcher+: switcher,
  },
}
