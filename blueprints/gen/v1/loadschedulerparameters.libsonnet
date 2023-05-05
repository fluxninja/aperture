{
  new():: {
  },
  withScheduler(scheduler):: {
    scheduler: scheduler,
  },
  withSchedulerMixin(scheduler):: {
    scheduler+: scheduler,
  },
  withSelectors(selectors):: {
    selectors:
      if std.isArray(selectors)
      then selectors
      else [selectors],
  },
  withSelectorsMixin(selectors):: {
    selectors+: selectors,
  },
  withWorkloadLatencyBasedTokens(workload_latency_based_tokens):: {
    workload_latency_based_tokens: workload_latency_based_tokens,
  },
  withWorkloadLatencyBasedTokensMixin(workload_latency_based_tokens):: {
    workload_latency_based_tokens+: workload_latency_based_tokens,
  },
}
