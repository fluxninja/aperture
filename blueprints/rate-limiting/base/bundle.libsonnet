local blueprint = import './rate-limiting.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params) {
  local c = std.mergePatch(config, params),

  local p = policy(c),
  policies: {
    [std.format('%s-cr.yaml', c.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', c.policy.policy_name)]: p.policyDef,
  },
}
