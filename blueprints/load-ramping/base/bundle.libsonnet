local creator = import '../../grafana/dashboard_group.libsonnet';
local utils = import '../../utils/utils.libsonnet';
local blueprint = import './load-ramping.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params, metadata={}) {
  local c = std.mergePatch(config, params),
  local metadataWrapper = metadata { values: std.toString(params) },

  local prepare_driver = function(metrics_name, criteria_name) {
    driver: {
      query_string: std.format('avg(%s{k8s_%s_name="%s",k8s_namespace_name="%s",policy_name="%s",infra_meter_name="%s"})', [
        metrics_name,
        std.asciiLower(c.policy.kubelet_metrics.infra_context.kind),
        c.policy.kubelet_metrics.infra_context.name,
        c.policy.kubelet_metrics.infra_context.namespace,
        c.policy.policy_name,
        'kubeletstats',
      ]),
      criteria: {
        ['%s' % criteria]: {
          threshold: c.policy.kubelet_metrics.criteria[criteria_name][criteria].threshold,
          operator: if criteria == 'forward' then 'lt' else 'gt',
        }
        for criteria in std.objectFields(c.policy.kubelet_metrics.criteria[criteria_name])
      },
    },
  },

  local pod_cpu_kubelet_driver =
    if std.objectHas(c.policy, 'kubelet_metrics') &&
       std.objectHas(c.policy.kubelet_metrics, 'criteria') &&
       std.objectHas(c.policy.kubelet_metrics.criteria, 'pod_cpu') then
      prepare_driver('k8s_pod_cpu_utilization_ratio', 'pod_cpu').driver
    else {},

  local pod_memory_kubelet_driver =
    if std.objectHas(c.policy, 'kubelet_metrics') &&
       std.objectHas(c.policy.kubelet_metrics, 'criteria') &&
       std.objectHas(c.policy.kubelet_metrics.criteria, 'pod_memory') then
      prepare_driver('k8s_pod_memory_usage_bytes', 'pod_memory').driver
    else {},

  local updated_cfg = c {
    policy+: {
      drivers+: {
        local promqlDrivers = if std.objectHas(c.policy, 'drivers') && std.objectHas(c.policy.drivers, 'promql_drivers') then c.policy.drivers.promql_drivers else [],
        promql_drivers+: promqlDrivers +
                         (if std.length(std.objectFields(pod_cpu_kubelet_driver)) > 0 then [pod_cpu_kubelet_driver] else []) +
                         (if std.length(std.objectFields(pod_memory_kubelet_driver)) > 0 then [pod_memory_kubelet_driver] else []),
      },
      resources+: {
        infra_meters:
          local infraMeters = if std.objectHas(c.policy.resources, 'infra_meters') then c.policy.resources.infra_meters else {};
          if std.objectHas(pod_cpu_kubelet_driver, 'query_string') ||
             std.objectHas(pod_memory_kubelet_driver, 'query_string') then
            assert !std.objectHas(infraMeters, 'kubeletstats') : 'An infra meter with name kubeletstats already exists. Please choose a different name.';
            utils.add_kubeletstats_infra_meter(
              infraMeters,
              if std.objectHas(c.policy.kubelet_metrics.infra_context, 'agent_group') then c.policy.kubelet_metrics.infra_context.agent_group else 'default',
              c.policy.kubelet_metrics.infra_context { agent_group:: '' },
            )
          else infraMeters,
      },
    },
  },

  local p = policy(updated_cfg, metadataWrapper),
  local d = creator(p.policyResource, updated_cfg),

  policies: {
    [std.format('%s-cr.yaml', updated_cfg.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', updated_cfg.policy.policy_name)]: p.policyDef { metadata: metadataWrapper },
  },
  dashboards: {
    [std.format('%s.json', updated_cfg.policy.policy_name)]: d.mainDashboard,
    [std.format('signals-%s.json', updated_cfg.policy.policy_name)]: d.signalsDashboard,
  } + d.receiverDashboards,
}
