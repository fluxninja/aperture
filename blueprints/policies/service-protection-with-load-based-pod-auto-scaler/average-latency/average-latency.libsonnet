{
  config: import 'config.libsonnet',
  policy: import '../../service-protection/average-latency/policy.libsonnet',
  dashboard: import '../../service-protection/average-latency/dashboard.libsonnet',
  bundle: import 'bundle.libsonnet',
}
