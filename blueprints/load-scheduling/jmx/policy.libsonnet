local averageLatencyFn = import '../average-latency/policy.libsonnet';
local jmxUtils = import './utils.libsonnet';

local config = blueprint.config;

function(cfg, metadata={}) {
  local averageLatencyPolicy = averageLatencyFn(cfg, metadata),
  local c = std.mergePatch(config, cfg),

  local policyDef = averageLatencyPolicy.policyDef {
    resources+: {
      infra_meters+: {
        jmxUtils(),
      },
    },
    circuit+: {
      components+: {
        flow_control+: {
          adaptive_load_scheduler+: {
            overload_confirmation+: [
              spec.v1.OverloadConfirmation.new()
              + spec.v1.OverloadConfirmation.withQueryString('avg(java_lang_OperatingSystem_CpuLoad{k8s_pod_name=~"service3-demo-app-.*"})')
              + spec.v1.OverloadConfirmation.withThreshold('0.35')
              + spec.v1.OverloadConfirmation.withOperator('gt'),
              spec.v1.OverloadConfirmation.new()
              + spec.v1.OverloadConfirmation.withQueryString('avg(java_lang_Copy_LastGcInfo_duration{k8s_pod_name=~"service3-demo-app-.*"})')
              + spec.v1.OverloadConfirmation.withThreshold('30')
              + spec.v1.OverloadConfirmation.withOperator('gt'),
            ],
            },
          },
        },
      },
    },

  policyResource: averageLatencyPolicy.policyResource {
    spec+: policyDef,
  },
  policyDef: policyDef,
}
