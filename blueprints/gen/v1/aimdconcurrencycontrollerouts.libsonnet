{
  new():: {
  },
  withAcceptedConcurrency(accepted_concurrency):: {
    accepted_concurrency: accepted_concurrency,
  },
  withAcceptedConcurrencyMixin(accepted_concurrency):: {
    accepted_concurrency+: accepted_concurrency,
  },
  withDesiredLoadMultiplier(desired_load_multiplier):: {
    desired_load_multiplier: desired_load_multiplier,
  },
  withDesiredLoadMultiplierMixin(desired_load_multiplier):: {
    desired_load_multiplier+: desired_load_multiplier,
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
  withObservedLoadMultiplier(observed_load_multiplier):: {
    observed_load_multiplier: observed_load_multiplier,
  },
  withObservedLoadMultiplierMixin(observed_load_multiplier):: {
    observed_load_multiplier+: observed_load_multiplier,
  },
}
