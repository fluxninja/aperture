local blueprint = import './latency-aimd-concurrency-limiting.libsonnet';

local policy = blueprint.policy;
local dashboard = blueprint.dashboard;
local config = blueprint.config;

{
  policies: {
    [std.format('%s.yaml', $._config.common.policy_name)]: policy($._config.common + $._config.policy).policyResource,
  },
  dashboards: {
    [std.format('%s.json', $._config.common.policy_name)]: dashboard($._config.common + $._config.dashboard).dashboard,
  },
} +
{
  _config:: config,
}
