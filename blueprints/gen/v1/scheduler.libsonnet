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
  withPrioritiesLabelKey(priorities_label_key):: {
    priorities_label_key: priorities_label_key,
  },
  withPrioritiesLabelKeyMixin(priorities_label_key):: {
    priorities_label_key+: priorities_label_key,
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
