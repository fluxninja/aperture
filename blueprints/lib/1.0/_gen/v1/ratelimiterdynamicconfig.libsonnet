{
  new():: {
  },
  withOverrides(overrides):: {
    overrides:
      if std.isArray(overrides)
      then overrides
      else [overrides],
  },
  withOverridesMixin(overrides):: {
    overrides+: overrides,
  },
}
