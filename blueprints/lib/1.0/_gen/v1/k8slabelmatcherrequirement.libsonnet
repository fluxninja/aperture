{
  new():: {
  },
  withKey(key):: {
    key: key,
  },
  withKeyMixin(key):: {
    key+: key,
  },
  withOperator(operator):: {
    operator: operator,
  },
  withOperatorMixin(operator):: {
    operator+: operator,
  },
  withValues(values):: {
    values:
      if std.isArray(values)
      then values
      else [values],
  },
  withValuesMixin(values):: {
    values+: values,
  },
}
