{
  new():: {
  },
  withOf(of):: {
    of:
      if std.isArray(of)
      then of
      else [of],
  },
  withOfMixin(of):: {
    of+: of,
  },
}
