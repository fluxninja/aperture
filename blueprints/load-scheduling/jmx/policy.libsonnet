local spec = import '../../spec.libsonnet';
local promqlFn = import '../promql/policy.libsonnet';
local blueprint = import './jmx.libsonnet';
local jmxUtils = import './utils.libsonnet';

local config = blueprint.config;

function(cfg, params={}, metadata={}) {
  local updated_cfg = cfg {
    policy+: {
      promql_query: 'avg(java_lang_G1_Young_Generation_LastGcInfo_duration{k8s_pod_name=~"service3-demo-app-.*"})',
      setpoint: 20,
    },
  },
  local c = std.mergePatch(config, updated_cfg),
  local promqlPolicy = promqlFn(c, params, metadata),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
