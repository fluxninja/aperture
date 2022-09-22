local dashboard = import '../../lib/1.0/dashboards/latency-gradient.libsonnet';
local policy = import '../../lib/1.0/policies/latency-gradient.libsonnet';

local config = import './config.libsonnet';

{
  policies: {
    [std.format('%s.yaml', $._config.policy.policyName)]: policy($._config.policy),
  },
  dashboards: {
    [std.format('%s.json', $._config.common.policyName)]: dashboard($._config.dashboard),
  },
} +
{
  _config:: config,
}
