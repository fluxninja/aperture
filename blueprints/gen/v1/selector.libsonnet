{
  new():: {
  },
  withAgentGroup(agent_group):: {
    agent_group: agent_group,
  },
  withAgentGroupMixin(agent_group):: {
    agent_group+: agent_group,
  },
  withControlPoint(control_point):: {
    control_point: control_point,
  },
  withControlPointMixin(control_point):: {
    control_point+: control_point,
  },
  withLabelMatcher(label_matcher):: {
    label_matcher: label_matcher,
  },
  withLabelMatcherMixin(label_matcher):: {
    label_matcher+: label_matcher,
  },
  withService(service):: {
    service: service,
  },
  withServiceMixin(service):: {
    service+: service,
  },
}
