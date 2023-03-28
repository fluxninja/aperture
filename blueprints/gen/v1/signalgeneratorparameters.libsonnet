{
  new():: {
  },
  withEndBehavior(end_behavior):: {
    end_behavior: end_behavior,
  },
  withEndBehaviorMixin(end_behavior):: {
    end_behavior+: end_behavior,
  },
  withSteps(steps):: {
    steps:
      if std.isArray(steps)
      then steps
      else [steps],
  },
  withStepsMixin(steps):: {
    steps+: steps,
  },
}
