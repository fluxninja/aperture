local blueprint = import './rate-limiter.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

{
  dashboards: {},
  policies: {
    [std.format('%s.yaml', $._config.common.policyName)]: policy($._config.common + $._config.policy).policyResource,
  },
} + { _config:: config }
