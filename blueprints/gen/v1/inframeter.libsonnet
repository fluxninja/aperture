{
  new():: {
  },
  withAgentGroup(agent_group):: {
    agent_group: agent_group,
  },
  withAgentGroupMixin(agent_group):: {
    agent_group+: agent_group,
  },
  withPerAgentGroup(per_agent_group):: {
    per_agent_group: per_agent_group,
  },
  withPerAgentGroupMixin(per_agent_group):: {
    per_agent_group+: per_agent_group,
  },
  withPipeline(pipeline):: {
    pipeline: pipeline,
  },
  withPipelineMixin(pipeline):: {
    pipeline+: pipeline,
  },
  withProcessors(processors):: {
    processors: processors,
  },
  withProcessorsMixin(processors):: {
    processors+: processors,
  },
  withReceivers(receivers):: {
    receivers: receivers,
  },
  withReceiversMixin(receivers):: {
    receivers+: receivers,
  },
}
