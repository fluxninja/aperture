local spec = import '../../spec.libsonnet';
local averageLatencyFn = import '../average-latency/policy.libsonnet';
local blueprint = import './jmx.libsonnet';
local jmxUtils = import './utils.libsonnet';

local config = blueprint.config;

function(cfg, params={}, metadata={}) {
  local averageLatencyPolicy = averageLatencyFn(cfg, params, metadata),
  local c = std.mergePatch(config, cfg),

  // add overload confirmations
  local policyDef = averageLatencyPolicy.policyDef {
    service_protection_core+: {
      overload_confirmations+: [
        {
          query_string: c.policy.jmx.cpu_query,
          threshold: c.policy.jmx.cpu_threshold,
          operator: 'gt',
        },
        {
          query_string: c.policy.jmx.gc_query,
          threshold: c.policy.jmx.gc_threshold,
          operator: 'gt',
        },
      ],
    },
  },

  policyResource: averageLatencyPolicy.policyResource {
    spec+: policyDef,
  },
  policyDef: policyDef,
}
