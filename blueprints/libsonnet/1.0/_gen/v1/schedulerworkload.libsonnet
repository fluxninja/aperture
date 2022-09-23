{
  new():: {
  },
  withLabelMatcher(label_matcher):: {
    label_matcher: label_matcher,
  },
  withLabelMatcherMixin(label_matcher):: {
    label_matcher+: label_matcher,
  },
  withWorkloadParameters(workload_parameters):: {
    workload_parameters: workload_parameters,
  },
  withWorkloadParametersMixin(workload_parameters):: {
    workload_parameters+: workload_parameters,
  },
}
