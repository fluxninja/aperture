local blueprint = import './adaptive-load-scheduler.libsonnet';
local dashboard = blueprint.dashboard;
local config = blueprint.config;

function(params) {
  // validate params within this function
  // check whether common.policy_name is set
  validate:: if !std.objectHas(params, 'common') then
    error 'common is not set'
  else if !std.objectHas(params.common, 'policy_name') then
    error 'common.policy_name is not set'
  else
    true,

  // make sure param object contains fields that are in config
  local extra_keys = std.setDiff(std.objectFields(params), std.objectFields(config)),
  assert std.length(extra_keys) == 0 : 'Unknown keys in params: ' + extra_keys,

  local c = std.mergePatch(config, params),

  local d = dashboard(c.common + c.dashboard),

  dashboards: {
    [std.format('%s.json', c.common.policy_name)]: d.dashboard,
  },
}
