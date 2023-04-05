{
  new():: {
  },
  withMaxGradient(max_gradient):: {
    max_gradient: max_gradient,
  },
  withMaxGradientMixin(max_gradient):: {
    max_gradient+: max_gradient,
  },
  withSlope(slope):: {
    slope: slope,
  },
  withSlopeMixin(slope):: {
    slope+: slope,
  },
}
