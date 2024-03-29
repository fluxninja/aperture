{
  new():: {
  },
  withRego(rego):: {
    rego: rego,
  },
  withRegoMixin(rego):: {
    rego+: rego,
  },
  withRules(rules):: {
    rules: rules,
  },
  withRulesMixin(rules):: {
    rules+: rules,
  },
  withSelectors(selectors):: {
    selectors:
      if std.isArray(selectors)
      then selectors
      else [selectors],
  },
  withSelectorsMixin(selectors):: {
    selectors+: selectors,
  },
}
