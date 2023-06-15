{
  new():: {
  },
  withGradient(gradient):: {
    gradient: gradient,
  },
  withGradientMixin(gradient):: {
    gradient+: gradient,
  },
  withPeriodic(periodic):: {
    periodic: periodic,
  },
  withPeriodicMixin(periodic):: {
    periodic+: periodic,
  },
}
