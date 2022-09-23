{
  new():: {
  },
  withAgentGroup(agent_group):: {
    agent_group: agent_group,
  },
  withAgentGroupMixin(agent_group):: {
    agent_group+: agent_group,
  },
  withService(service):: {
    service: service,
  },
  withServiceMixin(service):: {
    service+: service,
  },
}
