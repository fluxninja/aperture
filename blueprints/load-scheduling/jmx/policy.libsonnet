local spec = import '../../spec.libsonnet';
local promqlFn = import '../promql/policy.libsonnet';
local blueprint = import './jmx.libsonnet';
local jmxUtils = import './utils.libsonnet';

local config = blueprint.config;

function(cfg, params={}, metadata={}) {
  local c = std.mergePatch(config, cfg),
  local promqlPolicy = promqlFn(c, params, metadata),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
