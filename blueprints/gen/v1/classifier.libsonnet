{
  new():: {
  },
  withFlowSelector(flow_selector):: {
    flow_selector: flow_selector,
  },
  withFlowSelectorMixin(flow_selector):: {
    flow_selector+: flow_selector,
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
}
