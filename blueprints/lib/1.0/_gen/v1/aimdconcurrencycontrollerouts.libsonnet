{
  new():: {
  },
  withAcceptedConcurrency(accepted_concurrency):: {
    accepted_concurrency: accepted_concurrency,
  },
  withAcceptedConcurrencyMixin(accepted_concurrency):: {
    accepted_concurrency+: accepted_concurrency,
  },
  withDesiredConcurrency(desired_concurrency):: {
    desired_concurrency: desired_concurrency,
  },
  withDesiredConcurrencyMixin(desired_concurrency):: {
    desired_concurrency+: desired_concurrency,
  },
  withIncomingConcurrency(incoming_concurrency):: {
    incoming_concurrency: incoming_concurrency,
  },
  withIncomingConcurrencyMixin(incoming_concurrency):: {
    incoming_concurrency+: incoming_concurrency,
  },
  withIsOverload(is_overload):: {
    is_overload: is_overload,
  },
  withIsOverloadMixin(is_overload):: {
    is_overload+: is_overload,
  },
  withLoadMultiplier(load_multiplier):: {
    load_multiplier: load_multiplier,
  },
  withLoadMultiplierMixin(load_multiplier):: {
    load_multiplier+: load_multiplier,
  },
}
