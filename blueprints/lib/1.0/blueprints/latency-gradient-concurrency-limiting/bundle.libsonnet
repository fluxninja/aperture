local blueprint = import './latency-gradient-concurrency-limiting.libsonnet';

local policy = blueprint.policy;
local dashboard = blueprint.dashboard;
local config = blueprint.config;

{
  policies: {
    [std.format('%s.yaml', $._config.common.policyName)]: policy($._config.common + $._config.policy).policyResource,
  },
  dashboards: {
    [std.format('%s.json', $._config.common.policyName)]: dashboard($._config.common + $._config.dashboard).dashboard,
  },
} +
{
  _config:: config,
}
