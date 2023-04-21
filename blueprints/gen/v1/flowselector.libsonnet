{
  new():: {
  },
  withAgentGroup(agent_group):: {
    agent_group: agent_group,
  },
  withAgentGroupMixin(agent_group):: {
    agent_group+: agent_group,
  },
  withFlowMatcher(flow_matcher):: {
    flow_matcher: flow_matcher,
  },
  withFlowMatcherMixin(flow_matcher):: {
    flow_matcher+: flow_matcher,
  },
  withServiceSelector(service_selector):: {
    service_selector: service_selector,
  },
  withServiceSelectorMixin(service_selector):: {
    service_selector+: service_selector,
  },
}
