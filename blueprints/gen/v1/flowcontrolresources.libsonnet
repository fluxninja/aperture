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
  withFluxMeters(flux_meters):: {
    flux_meters: flux_meters,
  },
  withFluxMetersMixin(flux_meters):: {
    flux_meters+: flux_meters,
  },
}
