local blueprint = import './signals.libsonnet';

local dashboard = blueprint.dashboard;
local config = blueprint.config;

{
  dashboards: {
    [std.format('%s.json', $._config.common.policyName)]: dashboard($._config.common + $._config.dashboard).dashboard,
  },
} +
{
  _config:: config,
}
