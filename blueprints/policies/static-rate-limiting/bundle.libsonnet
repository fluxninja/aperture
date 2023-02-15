local blueprint = import './static-rate-limiting.libsonnet';

local policy = blueprint.policy;
local dashboard = blueprint.dashboard;
local config = blueprint.config;

function(params) {
  // validate params within this function
  // check whether common.policy_name is set
  // check whether policy.rate_limiter.rate_limit is set
  // check whether policy.rate_limiter.flow_selector is set
  // check whether policy.rate_limiter.parameters.label_key is set
  validate:: if !std.objectHas(params, 'common') then
    error 'common is not set'
  else if !std.objectHas(params.common, 'policy_name') then
    error 'common.policy_name is not set'
  else if !std.objectHas(params, 'policy') then
    error 'policy is not set'
  else if !std.objectHas(params.policy, 'rate_limiter') then
    error 'policy.rate_limiter is not set'
  else if !std.objectHas(params.policy.rate_limiter, 'rate_limit') then
    error 'policy.rate_limiter.rate_limit is not set'
  else if !std.objectHas(params.policy.rate_limiter, 'flow_selector') then
    error 'policy.rate_limiter.flow_selector is not set'
  else if !std.objectHas(params.policy.rate_limiter, 'parameters') then
    error 'policy.rate_limiter.parameters is not set'
  else if !std.objectHas(params.policy.rate_limiter.parameters, 'label_key') then
    error 'policy.rate_limiter.parameters.label_key is not set'
  else
    true,

  // make sure param object contains fields that are in config
  local extra_keys = std.setDiff(std.objectFields(params), std.objectFields(config)),
  assert std.length(extra_keys) == 0 : 'Unknown keys in params: ' + extra_keys,

  local c = std.mergePatch(config, params),

  local p = policy(c.common + c.policy),
  local d = dashboard(c.common + c.dashboard),

  dashboards: {
    [std.format('%s.json', c.common.policy_name)]: d.dashboard,
  },
  policies: {
    [std.format('%s-cr.yaml', c.common.policy_name)]: p.policyResource,
    [std.format('%s.yaml', c.common.policy_name)]: p.policyDef,
  },
}
