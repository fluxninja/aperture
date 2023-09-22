local creator = import '../../grafana/creator.libsonnet';
local utils = import '../../utils/utils.libsonnet';
local blueprint = import './pod-auto-scaler.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params, metadata={}) {
  // make sure param object contains fields that are in config
  local extra_keys = std.setDiff(std.objectFields(params), std.objectFields(config)),
  assert std.length(extra_keys) == 0 : 'Unknown keys in params: ' + extra_keys,

  local c = std.mergePatch(config, params),
  local metadataWrapper = metadata { values: std.toString(params) },

  local prepare_controller = function(metrics_name, gradient_slope, threshold, alert_name) {
    controller: [{
      query_string: std.format('avg(%s{k8s_%s_name="%s",k8s_namespace_name="%s",policy_name="%s",infra_meter_name="%s"})', [
        metrics_name,
        std.asciiLower(c.policy.scaling_backend.kubernetes_replicas.kubernetes_object_selector.kind),
        c.policy.scaling_backend.kubernetes_replicas.kubernetes_object_selector.name,
        c.policy.scaling_backend.kubernetes_replicas.kubernetes_object_selector.namespace,
        c.policy.policy_name,
        'kubeletstats',
      ]),
      setpoint: threshold,
      gradient: {
        slope: gradient_slope,
      },
      alerter: {
        alert_name: alert_name,
      },
    }],
  },

  local pod_cpu_scale_in_controllers =
    if std.objectHas(c.policy, 'pod_cpu') &&
       std.objectHas(c.policy.pod_cpu, 'scale_in') &&
       std.objectHas(c.policy.pod_cpu.scale_in, 'enabled') &&
       c.policy.pod_cpu.scale_in.enabled then
      prepare_controller('k8s_pod_cpu_utilization_ratio', -1, c.policy.pod_cpu.scale_in.threshold, 'Pod CPU based scale in intended').controller
    else [],

  local pod_cpu_scale_out_controllers =
    if std.objectHas(c.policy, 'pod_cpu') &&
       std.objectHas(c.policy.pod_cpu, 'scale_out') &&
       std.objectHas(c.policy.pod_cpu.scale_out, 'enabled') &&
       c.policy.pod_cpu.scale_out.enabled then
      prepare_controller('k8s_pod_cpu_utilization_ratio', 1, c.policy.pod_cpu.scale_out.threshold, 'Pod CPU based scale out intended').controller
    else [],

  local pod_memory_scale_in_controllers =
    if std.objectHas(c.policy, 'pod_memory') &&
       std.objectHas(c.policy.pod_memory, 'scale_in') &&
       std.objectHas(c.policy.pod_memory.scale_in, 'enabled') &&
       c.policy.pod_memory.scale_in.enabled then
      prepare_controller('k8s_pod_memory_usage_bytes', -1, c.policy.pod_memory.scale_in.threshold, 'Pod Memory based scale in intended').controller
    else [],

  local pod_memory_scale_out_controllers =
    if std.objectHas(c.policy, 'pod_memory') &&
       std.objectHas(c.policy.pod_memory, 'scale_out') &&
       std.objectHas(c.policy.pod_memory.scale_out, 'enabled') &&
       c.policy.pod_memory.scale_out.enabled then
      prepare_controller('k8s_pod_memory_usage_bytes', 1, c.policy.pod_memory.scale_out.threshold, 'Pod Memory based scale out intended').controller
    else [],

  local updated_cfg = c {
    policy+: {
      local promqlScaleOutControllers = if std.objectHas(c.policy, 'promql_scale_out_controllers') then c.policy.promql_scale_out_controllers else [],
      promql_scale_out_controllers: promqlScaleOutControllers +
                                    (if std.length(pod_cpu_scale_out_controllers) > 0 then pod_cpu_scale_out_controllers else []) +
                                    (if std.length(pod_memory_scale_out_controllers) > 0 then pod_memory_scale_out_controllers else []),
      local promqlScaleInControllers = if std.objectHas(c.policy, 'promql_scale_in_controllers') then c.policy.promql_scale_in_controllers else [],
      promql_scale_in_controllers+: promqlScaleInControllers +
                                    (if std.length(pod_cpu_scale_in_controllers) > 0 then pod_cpu_scale_in_controllers else []) +
                                    (if std.length(pod_memory_scale_in_controllers) > 0 then pod_memory_scale_in_controllers else []),
      resources+: {
        infra_meters:
          local infraMeters = if std.objectHas(c.policy.resources, 'infra_meters') then c.policy.resources.infra_meters else {};
          if std.length(pod_cpu_scale_in_controllers) > 0 ||
             std.length(pod_cpu_scale_out_controllers) > 0 ||
             std.length(pod_memory_scale_in_controllers) > 0 ||
             std.length(pod_memory_scale_out_controllers) > 0 then
            assert !std.objectHas(infraMeters, 'kubeletstats') : 'An infra meter with name kubeletstats already exists. Please choose a different name.';
            utils.add_kubeletstats_infra_meter(
              infraMeters,
              c.policy.scaling_backend.kubernetes_replicas.kubernetes_object_selector.agent_group,
              c.policy.scaling_backend.kubernetes_replicas.kubernetes_object_selector { agent_group:: '' },
            )
          else infraMeters,
      },
    },
  },

  local p = policy(updated_cfg, metadataWrapper),
  local d = creator(p.policyResource, updated_cfg),

  dashboards: {
    [std.format('%s.json', updated_cfg.policy.policy_name)]: d.dashboard,
  } + d.receiverDashboards,

  policies: {
    [std.format('%s-cr.yaml', updated_cfg.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', updated_cfg.policy.policy_name)]: p.policyDef { metadata: metadataWrapper },
  },
}
