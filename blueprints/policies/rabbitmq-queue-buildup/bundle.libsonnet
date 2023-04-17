local blueprint = import './rabbitmq-queue-buildup.libsonnet';

local policy = blueprint.policy;
local dashboard = blueprint.dashboard;
local config = blueprint.config;

function(params) {
  // make sure param object contains fields that are in config
  local extra_keys = std.setDiff(std.objectFields(params), std.objectFields(config)),
  assert std.length(extra_keys) == 0 : 'Unknown keys in params: ' + extra_keys,

  local c = std.mergePatch(config, params),

  local p = policy(c.common + c.policy),
  local d = dashboard(c.common + c.dashboard),

  policies: {
    [std.format('%s-cr.yaml', c.common.policy_name)]: p.policyResource,
    [std.format('%s.yaml', c.common.policy_name)]: p.policyDef,
  },
  dashboards: {
    [std.format('%s.json', c.common.policy_name)]: d.dashboard,
  },
}
