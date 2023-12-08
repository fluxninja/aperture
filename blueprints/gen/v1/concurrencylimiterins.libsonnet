{
  new():: {
  },
  withMaxConcurrency(max_concurrency):: {
    max_concurrency: max_concurrency,
  },
  withMaxConcurrencyMixin(max_concurrency):: {
    max_concurrency+: max_concurrency,
  },
  withPassThrough(pass_through):: {
    pass_through: pass_through,
  },
  withPassThroughMixin(pass_through):: {
    pass_through+: pass_through,
  },
}
