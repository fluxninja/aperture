{
  new():: {
  },
  withLoadShapers(load_shapers):: {
    load_shapers:
      if std.isArray(load_shapers)
      then load_shapers
      else [load_shapers],
  },
  withLoadShapersMixin(load_shapers):: {
    load_shapers+: load_shapers,
  },
}
