{
  new():: {
  },
  withLabelMatcher(label_matcher):: {
    label_matcher: label_matcher,
  },
  withLabelMatcherMixin(label_matcher):: {
    label_matcher+: label_matcher,
  },
  withParameters(parameters):: {
    parameters: parameters,
  },
  withParametersMixin(parameters):: {
    parameters+: parameters,
  },
}
