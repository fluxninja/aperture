{
  new():: {
  },
  withRegulator(regulator):: {
    regulator: regulator,
  },
  withRegulatorMixin(regulator):: {
    regulator+: regulator,
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
