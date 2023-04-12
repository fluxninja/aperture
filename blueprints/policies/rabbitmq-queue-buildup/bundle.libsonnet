local blueprint = import './rabbitmq-queue-buildup.libsonnet';

local policy = blueprint.policy;
local dashboard = blueprint.dashboard;
local config = blueprint.config;

function(params) {
  // validate params within this function
  // check whether common.policy_name is set
  // check whether common.queue_name is set
  // check whether policy.concurrency_controller.flow_selector is set
  // check whether policy.concurrency_controller.queue_buildup_setpoint is set
  validate:: if !std.objectHas(params, 'common') then
    error 'must provide common object'
  else if !std.objectHas(params.common, 'policy_name') then
    error 'must provide common.policy_name'
  else if !std.objectHas(params.common, 'queue_name') then
    error 'must provide common.queue_name'
  else if !std.objectHas(params, 'policy') then
    error 'must provide policy object'
  else if !std.objectHas(params.policy, 'concurrency_controller') then
    error 'must provide policy.concurrency_controller'
  else if !std.objectHas(params.policy.concurrency_controller, 'flow_selector') then
    error 'must provide policy.concurrency_controller.flow_selector'
  else if !std.objectHas(params.policy.concurrency_controller, 'queue_buildup_setpoint') then
    error 'must provide policy.concurrency_controller.queue_buildup_setpoint'
  else true,

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
