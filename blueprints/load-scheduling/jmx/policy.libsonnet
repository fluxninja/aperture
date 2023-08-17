local spec = import '../../spec.libsonnet';
local averageLatencyFn = import '../average-latency/policy.libsonnet';
local blueprint = import './jmx.libsonnet';
local jmxUtils = import './utils.libsonnet';

local config = blueprint.config;

function(cfg, params={}, metadata={}) {
  local averageLatencyPolicy = averageLatencyFn(cfg, params, metadata),
  local c = std.mergePatch(config, cfg),

  local policyDef = averageLatencyPolicy.policyDef {
    resources+: {
      infra_meters+: jmxUtils(),
    },
  },

  policyResource: averageLatencyPolicy.policyResource {
    spec+: policyDef,
  },
  policyDef: policyDef,
}
