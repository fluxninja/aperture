local creator = import '../../grafana/dashboard_group.libsonnet';
local blueprint = import './rate-limiting.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params, metadata={}) {
  local c = std.mergePatch(config, params),

  local metadataWrapper = metadata { values: std.toString(params) },
  local p = policy(c, metadataWrapper),
  local d = creator(p.policyResource, c),

  dashboards: {
    [std.format('%s.json', c.policy.policy_name)]: d.mainDashboard,
    [std.format('signals-%s.json', c.policy.policy_name)]: d.signalsDashboard,
  } + d.receiverDashboards,
  policies: {
    [std.format('%s-cr.yaml', c.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', c.policy.policy_name)]: p.policyDef { metadata: metadataWrapper },
  },
}
