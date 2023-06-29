local creator = import '../../grafana/dashboard_group.libsonnet';
local blueprint = import './quota-scheduling.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params, metadata={}) {
  // make sure param object contains fields that are in config
  local extra_keys = std.setDiff(std.objectFields(params), std.objectFields(config)),
  assert std.length(extra_keys) == 0 : 'Unknown keys in params: ' + extra_keys,

  local c = std.mergePatch(config, params),

  local metadataWrapper = metadata { values: std.toString(params) },
  local p = policy(c, metadataWrapper),
  local d = creator(p.policyResource, c),

  dashboards: {
    [std.format('%s.json', c.policy.policy_name)]: d.mainDashboard,
    [std.format('signals-%s.json', c.policy.policy_name)]: d.signalsDashboard,
  },
  policies: {
    [std.format('%s-cr.yaml', c.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', c.policy.policy_name)]: p.policyDef { metadata: metadataWrapper },
  },
}
