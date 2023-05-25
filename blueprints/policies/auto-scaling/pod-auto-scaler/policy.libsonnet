local spec = import '../../../spec.libsonnet';
local basePolicyFn = import '../base/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg, metadata={}) {
  local params = config + cfg,

  local policyDef = basePolicyFn(cfg).policyDef,

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
      annotations: {
        [if std.objectHas(metadata, 'values') then 'fluxninja.com/values']: metadata.values,
        [if std.objectHas(metadata, 'blueprints_uri') then 'fluxninja.com/blueprint-uri']: metadata.blueprints_uri,
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,

  policyDef: policyDef,
}
