local creator = import '../../grafana/dashboard_group.libsonnet';
local utils = import '../common/utils.libsonnet';
local blueprint = import './promql.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params, metadata={}) {
  // make sure param object contains fields that are in config
  local extra_keys = std.setDiff(std.objectFields(params), std.objectFields(config)),
  assert std.length(extra_keys) == 0 : 'Unknown keys in params: ' + extra_keys,

  local c = std.mergePatch(config, params),
  local metadataWrapper = metadata { values: std.toString(params) },

  local updated_cfg = utils.add_kubelet_overload_confirmations(c).updated_cfg,

  local p = policy(updated_cfg, params, metadataWrapper),
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
