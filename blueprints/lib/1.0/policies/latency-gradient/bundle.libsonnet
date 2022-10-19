local lib = import '../../main.libsonnet';
local blueprint = import './latency-gradient.libsonnet';

local policy = blueprint.policy;
local dashboard = blueprint.dashboard;
local config = blueprint.config;

local signalsDashboard = lib.blueprints.dashboards.Signals;

{
  policies: {
    [std.format('%s.yaml', $._config.policy.policyName)]: policy($._config.policy).policyResource,
  },
  dashboards: {
    [std.format('%s-signals.json', $._config.policy.policyName)]: signalsDashboard($._config.signalsDashboard).dashboard,
    [std.format('%s.json', $._config.policy.policyName)]: dashboard($._config.dashboard).dashboard,
  },
} +
{
  _config:: config,
}
