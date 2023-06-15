{
  new():: {
  },
  withSampler(sampler):: {
    sampler: sampler,
  },
  withSamplerMixin(sampler):: {
    sampler+: sampler,
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
