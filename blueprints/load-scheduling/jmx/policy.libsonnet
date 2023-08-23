local spec = import '../../spec.libsonnet';
local averageLatencyFn = import '../average-latency/policy.libsonnet';
local blueprint = import './jmx.libsonnet';
local jmxUtils = import './utils.libsonnet';

local config = blueprint.config;

function(cfg, params={}, metadata={}) {
  local c = std.mergePatch(config, cfg),
  local averageLatencyPolicy = averageLatencyFn(c, params, metadata),

  policyResource: averageLatencyPolicy.policyResource {
    spec+: averageLatencyPolicy.policyDef,
  },
  policyDef: averageLatencyPolicy.policyDef,
}
