{
  new():: {
  },
  withValid(valid):: {
    valid: valid,
  },
  withValidMixin(valid):: {
    valid+: valid,
  },
  withValue(value):: {
    value: value,
  },
  withValueMixin(value):: {
    value+: value,
  },
}
