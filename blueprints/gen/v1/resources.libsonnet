{
  new():: {
  },
  withClassifiers(classifiers):: {
    classifiers:
      if std.isArray(classifiers)
      then classifiers
      else [classifiers],
  },
  withClassifiersMixin(classifiers):: {
    classifiers+: classifiers,
  },
  withFlowControl(flow_control):: {
    flow_control: flow_control,
  },
  withFlowControlMixin(flow_control):: {
    flow_control+: flow_control,
  },
  withFluxMeters(flux_meters):: {
    flux_meters: flux_meters,
  },
  withFluxMetersMixin(flux_meters):: {
    flux_meters+: flux_meters,
  },
}
