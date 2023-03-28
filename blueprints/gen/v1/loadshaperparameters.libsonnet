{
  new():: {
  },
  withEndBehavior(end_behavior):: {
    end_behavior: end_behavior,
  },
  withEndBehaviorMixin(end_behavior):: {
    end_behavior+: end_behavior,
  },
  withFluxRegulatorParameters(flux_regulator_parameters):: {
    flux_regulator_parameters: flux_regulator_parameters,
  },
  withFluxRegulatorParametersMixin(flux_regulator_parameters):: {
    flux_regulator_parameters+: flux_regulator_parameters,
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
