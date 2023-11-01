local kubeletUtils = import '../../utils/utils.libsonnet';
local utils = import '../common/utils.libsonnet';
local blueprint = import './cpu-overload-protection.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params) {
  local c = std.mergePatch(config, params),

  local policyName = c.policy.policy_name,
  local k8s_pod_name = c.policy.kubernetes_object_selector.name,
  local k8s_namespace_name = c.policy.kubernetes_object_selector.namespace,
  local promqlQuery = 'avg(k8s_pod_cpu_utilization_ratio{k8s_deployment_name="%(k8s_pod_name)s",k8s_namespace_name="%(k8s_namespace_name)s",policy_name="%(policyName)s", infra_meter_name="kubeletstats"}) * 100' % {
    k8s_pod_name: k8s_pod_name,
    k8s_namespace_name: k8s_namespace_name,
    policyName: policyName,
  },

  local updated_cfg = utils.add_kubelet_overload_confirmations(c).updated_cfg {
    policy+: {
      promql_query: promqlQuery,
      setpoint: c.policy.load_scheduling_core.setpoint,
      overload_condition: 'gt',
    },
  },

  local infraMeters = if std.objectHas(c.policy.resources, 'infra_meters') then c.policy.resources.infra_meters else {},
  assert !std.objectHas(infraMeters, 'kubeletstats') : 'An infra meter with name kubeletstats already exists. Please choose a different name.',

  local agent_group = if std.objectHas(c.policy, 'kubernetes_object_selector') && std.objectHas(c.policy.kubernetes_object_selector, 'agent_group') then c.policy.kubernetes_object_selector.agent_group else 'default',

  local kubeletstatsInfraMeter = kubeletUtils.add_kubeletstats_infra_meter({}, agent_group, updated_cfg.policy.kubernetes_object_selector),

  local config_with_kubeletstats_infra_meter = updated_cfg {
    policy+: {
      resources+: {
        infra_meters+: kubeletstatsInfraMeter,
      },
    },
  },

  local p = policy(config_with_kubeletstats_infra_meter),
  policies: {
    [std.format('%s-cr.yaml', config_with_kubeletstats_infra_meter.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', config_with_kubeletstats_infra_meter.policy.policy_name)]: p.policyDef,
  },
}
