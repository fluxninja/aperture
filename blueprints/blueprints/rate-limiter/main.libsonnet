local policy = import '../../lib/1.0/policies/rate-limiter.libsonnet';

local config = import './config.libsonnet';

{
  dashboards: {},
  policies: {
    [std.format('%s.yaml', $._config.policy.policyName)]: policy($._config.policy),
  },
} + { _config:: config }
