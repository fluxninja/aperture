local lib = import '../../lib/1.0/main.libsonnet';

local config = import './config.libsonnet';

local policy = lib.blueprints.policies.RateLimiter;

{
  dashboards: {},
  policies: {
    [std.format('%s.yaml', $._config.policy.policyName)]: policy($._config.policy),
  },
} + { _config:: config }
