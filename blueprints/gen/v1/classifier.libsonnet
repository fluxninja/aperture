{
  new():: {
  },
  withFlowSelector(flow_selector):: {
    flow_selector: flow_selector,
  },
  withFlowSelectorMixin(flow_selector):: {
    flow_selector+: flow_selector,
  },
  withRules(rules):: {
    rules: rules,
  },
  withRulesMixin(rules):: {
    rules+: rules,
  },
}
