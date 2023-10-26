local utils = import '../common/utils.libsonnet';
local blueprint = import './average-latency.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params) {
  local c = std.mergePatch(config, params),

  local updated_cfg = utils.add_kubelet_overload_confirmations(c).updated_cfg,

  local p = policy(updated_cfg, params),
  policies: {
    [std.format('%s-cr.yaml', updated_cfg.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', updated_cfg.policy.policy_name)]: p.policyDef,
  },
}
