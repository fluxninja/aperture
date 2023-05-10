{
  StaticRateLimiting: import 'static-rate-limiting/static-rate-limiting.libsonnet',
  FeatureRollout: import 'feature-rollout/base/feature-rollout.libsonnet',
  PromQLServiceProtection: import 'service-protection/promql/promql.libsonnet',
  ServiceProtectionAverageLatency: import 'service-protection/average-latency/average-latency.libsonnet',
  PodAutoScaler: import 'pod-auto-scaler/pod-auto-scaler.libsonnet',
}
