local schedulerouts = import './schedulerouts.libsonnet';
{
  new():: {
    out_ports: {
      accepted_concurrency: error 'Port accepted_concurrency is missing',
      incoming_concurrency: error 'Port incoming_concurrency is missing',
    },
  },
  outPorts:: schedulerouts,
  withAutoTokens(auto_tokens):: {
    auto_tokens: auto_tokens,
  },
  withAutoTokensMixin(auto_tokens):: {
    auto_tokens+: auto_tokens,
  },
  withDefaultWorkload(default_workload):: {
    default_workload: default_workload,
  },
  withDefaultWorkloadMixin(default_workload):: {
    default_workload+: default_workload,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withSelector(selector):: {
    selector: selector,
  },
  withSelectorMixin(selector):: {
    selector+: selector,
  },
  withWorkloads(workloads):: {
    workloads:
      if std.isArray(workloads)
      then workloads
      else [workloads],
  },
  withWorkloadsMixin(workloads):: {
    workloads+: workloads,
  },
}
