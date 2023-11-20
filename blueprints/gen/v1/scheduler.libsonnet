{
  new():: {
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
  withDeniedResponseStatusCode(denied_response_status_code):: {
    denied_response_status_code: denied_response_status_code,
  },
  withDeniedResponseStatusCodeMixin(denied_response_status_code):: {
    denied_response_status_code+: denied_response_status_code,
  },
  withPriorityLabelKey(priority_label_key):: {
    priority_label_key: priority_label_key,
  },
  withPriorityLabelKeyMixin(priority_label_key):: {
    priority_label_key+: priority_label_key,
  },
  withTokensLabelKey(tokens_label_key):: {
    tokens_label_key: tokens_label_key,
  },
  withTokensLabelKeyMixin(tokens_label_key):: {
    tokens_label_key+: tokens_label_key,
  },
  withWorkloadLabelKey(workload_label_key):: {
    workload_label_key: workload_label_key,
  },
  withWorkloadLabelKeyMixin(workload_label_key):: {
    workload_label_key+: workload_label_key,
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
