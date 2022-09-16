local schedulerouts = import './schedulerouts.libsonnet';
{
  new():: {
    out_ports: {
      accepted_concurrency_ms: error 'Port accepted_concurrency_ms is missing',
      incoming_concurrency_ms: error 'Port incoming_concurrency_ms is missing',
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
  withMaxTimeout(max_timeout):: {
    max_timeout: max_timeout,
  },
  withMaxTimeoutMixin(max_timeout):: {
    max_timeout+: max_timeout,
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
  withTimeoutFactor(timeout_factor):: {
    timeout_factor: timeout_factor,
  },
  withTimeoutFactorMixin(timeout_factor):: {
    timeout_factor+: timeout_factor,
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
