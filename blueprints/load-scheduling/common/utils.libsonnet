local utils = import '../../utils/utils.libsonnet';

local add_kubelet_overload_confirmations(c) = {
  local prepare_signal = function(metrics_name, threshold) {
    controller: {
      query_string: std.format('avg(%s{k8s_%s_name="%s",k8s_namespace_name="%s",policy_name="%s",infra_meter_name="%s"})', [
        metrics_name,
        std.asciiLower(c.policy.service_protection_core.kubelet_overload_confirmations.infra_context.kind),
        c.policy.service_protection_core.kubelet_overload_confirmations.infra_context.name,
        c.policy.service_protection_core.kubelet_overload_confirmations.infra_context.namespace,
        c.policy.policy_name,
        'kubeletstats',
      ]),
      threshold: threshold,
      operator: 'gt',
    },
  },

  local pod_cpu_overload_confirmation =
    if std.objectHas(c.policy, 'service_protection_core') &&
       std.objectHas(c.policy.service_protection_core, 'kubelet_overload_confirmations') &&
       std.objectHas(c.policy.service_protection_core.kubelet_overload_confirmations, 'criteria') &&
       std.objectHas(c.policy.service_protection_core.kubelet_overload_confirmations.criteria, 'pod_cpu') &&
       std.objectHas(c.policy.service_protection_core.kubelet_overload_confirmations.criteria.pod_cpu, 'enabled') &&
       c.policy.service_protection_core.kubelet_overload_confirmations.criteria.pod_cpu.enabled then
      prepare_signal('k8s_pod_cpu_utilization_ratio', c.policy.service_protection_core.kubelet_overload_confirmations.criteria.pod_cpu.threshold).controller
    else {},

  local pod_memory_overload_confirmation =
    if std.objectHas(c.policy, 'service_protection_core') &&
       std.objectHas(c.policy.service_protection_core, 'kubelet_overload_confirmations') &&
       std.objectHas(c.policy.service_protection_core.kubelet_overload_confirmations, 'criteria') &&
       std.objectHas(c.policy.service_protection_core.kubelet_overload_confirmations.criteria, 'pod_memory') &&
       std.objectHas(c.policy.service_protection_core.kubelet_overload_confirmations.criteria.pod_memory, 'enabled') &&
       c.policy.service_protection_core.kubelet_overload_confirmations.criteria.pod_memory.enabled then
      prepare_signal('k8s_pod_memory_usage_bytes', c.policy.service_protection_core.kubelet_overload_confirmations.criteria.pod_memory.threshold).controller
    else {},

  local updated_cfg = c {
    policy+: {
      service_protection_core+: {
        local overloadConfirmations = if std.objectHas(c.policy.service_protection_core, 'overload_confirmations') then
          c.policy.service_protection_core.overload_confirmations
        else [],
        overload_confirmations: overloadConfirmations + (if std.objectHas(pod_cpu_overload_confirmation, 'query_string') then [pod_cpu_overload_confirmation] else []) +
                                (if std.objectHas(pod_memory_overload_confirmation, 'query_string') then [pod_memory_overload_confirmation] else []),
      },
      resources+: {
        infra_meters:
          local infraMeters = if std.objectHas(c.policy.resources, 'infra_meters') then c.policy.resources.infra_meters else {};
          if std.objectHas(pod_cpu_overload_confirmation, 'query_string') ||
             std.objectHas(pod_memory_overload_confirmation, 'query_string') then
            assert !std.objectHas(infraMeters, 'kubeletstats') : 'An infra meter with name kubeletstats already exists. Please choose a different name.';
            utils.add_kubeletstats_infra_meter(
              infraMeters,
              if std.objectHas(c.policy.service_protection_core.kubelet_overload_confirmations.infra_context, 'agent_group') then
                c.policy.service_protection_core.kubelet_overload_confirmations.infra_context.agent_group else 'default',
              c.policy.service_protection_core.kubelet_overload_confirmations.infra_context { agent_group:: '' },
            )
          else infraMeters,
      },
    },
  },

  updated_cfg: updated_cfg,
};

{
  add_kubelet_overload_confirmations: add_kubelet_overload_confirmations,
}
