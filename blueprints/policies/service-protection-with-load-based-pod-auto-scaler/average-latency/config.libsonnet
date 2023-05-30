local averageLatencyServiceProtection = import '../../service-protection/average-latency/config.libsonnet';
local baseServiceProtectionDefaults = import '../../service-protection/base/config-defaults.libsonnet';

averageLatencyServiceProtection {
  /**
  * @param (policy: policies/service-protection/average-latency:param:policy required) Configuration for the Service Protection policy.
  */
  policy+: {
    auto_scaling: baseServiceProtectionDefaults.auto_scaling_pods,
  },
  /**
  * @param (dashboard: policies/service-protection/average-latency:param:dashboard) Configuration for the Grafana dashboard accompanying this policy.
  */
  dashboard: averageLatencyServiceProtection.dashboard,
}
