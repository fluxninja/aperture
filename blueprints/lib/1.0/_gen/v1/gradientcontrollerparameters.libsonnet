{
  new():: {
  },
  withMaxGradient(max_gradient):: {
    max_gradient: max_gradient,
  },
  withMaxGradientMixin(max_gradient):: {
    max_gradient+: max_gradient,
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
