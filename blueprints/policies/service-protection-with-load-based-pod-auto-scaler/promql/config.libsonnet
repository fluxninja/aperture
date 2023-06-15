local baseServiceProtectionDefaults = import '../../service-protection/base/config-defaults.libsonnet';
local promqlServiceProtection = import '../../service-protection/promql/config.libsonnet';

promqlServiceProtection {
  /**
  * @param (policy: policies/service-protection/promql:param:policy required) Configuration for the Service Protection policy.
  */
  policy+: {
    auto_scaling: baseServiceProtectionDefaults.auto_scaling_pods,
  },
  /**
  * @param (dashboard: policies/service-protection/promql:param:dashboard) Configuration for the Grafana dashboard accompanying this policy.
  */
  dashboard: promqlServiceProtection.dashboard,
}
