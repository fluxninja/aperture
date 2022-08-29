{
  new():: {
  },
  withRules(rules):: {
    rules: rules,
  },
  withRulesMixin(rules):: {
    rules+: rules,
  },
  withSelector(selector):: {
    selector: selector,
  },
  withSelectorMixin(selector):: {
    selector+: selector,
  },
}
