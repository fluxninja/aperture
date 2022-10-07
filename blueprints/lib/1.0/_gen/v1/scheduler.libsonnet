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
  withDefaultWorkloadParameters(default_workload_parameters):: {
    default_workload_parameters: default_workload_parameters,
  },
  withDefaultWorkloadParametersMixin(default_workload_parameters):: {
    default_workload_parameters+: default_workload_parameters,
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
