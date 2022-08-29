{
  new():: {
  },
  withLabelMatcher(label_matcher):: {
    label_matcher: label_matcher,
  },
  withLabelMatcherMixin(label_matcher):: {
    label_matcher+: label_matcher,
  },
  withWorkload(workload):: {
    workload: workload,
  },
  withWorkloadMixin(workload):: {
    workload+: workload,
  },
}
