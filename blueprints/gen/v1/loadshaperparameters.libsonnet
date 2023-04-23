{
  new():: {
  },
  withFlowRegulatorParameters(flow_regulator_parameters):: {
    flow_regulator_parameters: flow_regulator_parameters,
  },
  withFlowRegulatorParametersMixin(flow_regulator_parameters):: {
    flow_regulator_parameters+: flow_regulator_parameters,
  },
  withLoadRegulatorParameters(load_regulator_parameters):: {
    load_regulator_parameters: load_regulator_parameters,
  },
  withLoadRegulatorParametersMixin(load_regulator_parameters):: {
    load_regulator_parameters+: load_regulator_parameters,
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
