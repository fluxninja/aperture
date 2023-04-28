{
  new():: {
  },
  withRegulatorParameters(regulator_parameters):: {
    regulator_parameters: regulator_parameters,
  },
  withRegulatorParametersMixin(regulator_parameters):: {
    regulator_parameters+: regulator_parameters,
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
