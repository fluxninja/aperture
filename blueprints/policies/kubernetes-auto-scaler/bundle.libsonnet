local blueprint = import './kubernetes-auto-scaler.libsonnet';

local policy = blueprint.policy;
local dashboard = blueprint.dashboard;
local config = blueprint.config;

function(params) {
  // validate params within this function
  // check whether common.policy_name is set
  // check whether policy.kubernetes_object_selector is set
  // check whether policy.scale_in_criteria is set and not empty
  // check whether policy.scale_out_criteria is set and not empty
  // check whether policy.scale_in_criteria[].query.promql.query_string is set
  // check whether policy.scale_out_criteria[].query.promql.query_string is set

  validate:: if !std.objectHas(params, 'common') then
    error 'must provide common object'
  else if !std.objectHas(params.common, 'policy_name') then
    error 'must provide common.policy_name'
  else if !std.objectHas(params, 'policy') then
    error 'must provide policy object'
  else if !std.objectHas(params.policy, 'kubernetes_object_selector') then
    error 'must provide policy.kubernetes_object_selector'
  else if !std.objectHas(params.policy, 'scale_in_criteria') || std.length(params.policy.scale_in_criteria) == 0 then
    error 'must provide policy.scale_in_criteria'
  else if !std.objectHas(params.policy, 'scale_out_criteria') || std.length(params.policy.scale_out_criteria) == 0 then
    error 'must provide policy.scale_out_criteria'
  else if std.foldl(function(a, b) a || !std.objectHas(b, 'query') || !std.objectHas(b.query, 'promql') || !std.objectHas(b.query.promql, 'query_string'), false, params.policy.scale_in_criteria) then
    error 'must provide policy.scale_in_criteria[].query.promql.query_string'
  else if std.foldl(function(a, b) a || !std.objectHas(b, 'query') || !std.objectHas(b.query, 'promql') || !std.objectHas(b.query.promql, 'query_string'), false, params.policy.scale_out_criteria) then
    error 'must provide policy.scale_out_criteria[].query.promql.query_string'
  else if std.foldl(function(a, b) a || !std.objectHas(b, 'query') || !std.objectHas(b.query, 'promql') || !std.objectHas(b.query.promql, 'out_ports') || !std.objectHas(b.query.promql.out_ports, 'output') || !std.objectHas(b.query.promql.out_ports.output, 'signal_name'), false, params.policy.scale_in_criteria) then
    error 'must provide policy.scale_in_criteria[].query.promql.out_ports.output.signal_name'
  else if std.foldl(function(a, b) a || !std.objectHas(b, 'query') || !std.objectHas(b.query, 'promql') || !std.objectHas(b.query.promql, 'out_ports') || !std.objectHas(b.query.promql.out_ports, 'output') || !std.objectHas(b.query.promql.out_ports.output, 'signal_name'), false, params.policy.scale_out_criteria) then
    error 'must provide policy.scale_out_criteria[].query.promql.out_ports.output.signal_name'
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
