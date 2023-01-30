{
  new():: {
  },
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
