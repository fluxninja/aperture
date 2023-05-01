{
  new():: {
  },
  withLoadRamps(load_ramps):: {
    load_ramps:
      if std.isArray(load_ramps)
      then load_ramps
      else [load_ramps],
  },
  withLoadRampsMixin(load_ramps):: {
    load_ramps+: load_ramps,
  },
}
