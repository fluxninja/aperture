local lib = import '../../lib/1.0/main.libsonnet';

local config = import './config.libsonnet';

local policy = lib.blueprints.policies.LatencyGradient;
local dashboard = lib.blueprints.dashboards.LatencyGradient;
local signalsDashboard = lib.blueprints.dashboards.Signals;

{
  policies: {
    [std.format('%s.yaml', $._config.policy.policyName)]: policy($._config.policy).policyResource,
  },
  dashboards: {
    [std.format('%s.json', $._config.common.policyName)]: dashboard($._config.dashboard),
    [std.format('%s-signals.json', $._config.common.policyName)]: signalsDashboard($._config.dashboard),
  },
} +
{
  _config:: config,
}
