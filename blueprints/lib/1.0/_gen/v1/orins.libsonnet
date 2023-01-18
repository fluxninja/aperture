{
  new():: {
  },
  withInputs(inputs):: {
    inputs:
      if std.isArray(inputs)
      then inputs
      else [inputs],
  },
  withInputsMixin(inputs):: {
    inputs+: inputs,
  },
}
