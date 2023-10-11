local commonPolicyFn = import '../common/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config + cfg,

  local policyDef = commonPolicyFn(cfg).policyDef,

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,

  policyDef: policyDef,
}
