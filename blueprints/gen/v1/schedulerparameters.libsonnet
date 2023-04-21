{
  new():: {
  },
  withAutoTokens(auto_tokens):: {
    auto_tokens: auto_tokens,
  },
  withAutoTokensMixin(auto_tokens):: {
    auto_tokens+: auto_tokens,
  },
  withDecisionDeadlineMargin(decision_deadline_margin):: {
    decision_deadline_margin: decision_deadline_margin,
  },
  withDecisionDeadlineMarginMixin(decision_deadline_margin):: {
    decision_deadline_margin+: decision_deadline_margin,
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
  withTokensLabelKey(tokens_label_key):: {
    tokens_label_key: tokens_label_key,
  },
  withTokensLabelKeyMixin(tokens_label_key):: {
    tokens_label_key+: tokens_label_key,
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
