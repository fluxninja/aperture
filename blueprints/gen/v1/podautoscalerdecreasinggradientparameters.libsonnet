{
  new():: {
  },
  withMinGradient(min_gradient):: {
    min_gradient: min_gradient,
  },
  withMinGradientMixin(min_gradient):: {
    min_gradient+: min_gradient,
  },
  withSlope(slope):: {
    slope: slope,
  },
  withSlopeMixin(slope):: {
    slope+: slope,
  },
}
