{
  new():: {
  },
  withDesiredConcurrency(desired_concurrency):: {
    desired_concurrency: desired_concurrency,
  },
  withDesiredConcurrencyMixin(desired_concurrency):: {
    desired_concurrency+: desired_concurrency,
  },
}
