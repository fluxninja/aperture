{
  new():: {
  },
  withComponents(components):: {
    components:
      if std.isArray(components)
      then components
      else [components],
  },
  withComponentsMixin(components):: {
    components+: components,
  },
  withEvaluationInterval(evaluation_interval):: {
    evaluation_interval: evaluation_interval,
  },
  withEvaluationIntervalMixin(evaluation_interval):: {
    evaluation_interval+: evaluation_interval,
  },
}
