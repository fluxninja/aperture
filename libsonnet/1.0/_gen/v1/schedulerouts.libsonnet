{
  new():: {
  },
  withAcceptedConcurrency(accepted_concurrency):: {
    accepted_concurrency: accepted_concurrency,
  },
  withAcceptedConcurrencyMixin(accepted_concurrency):: {
    accepted_concurrency+: accepted_concurrency,
  },
  withIncomingConcurrency(incoming_concurrency):: {
    incoming_concurrency: incoming_concurrency,
  },
  withIncomingConcurrencyMixin(incoming_concurrency):: {
    incoming_concurrency+: incoming_concurrency,
  },
}
