{
  RateLimiting: import 'rate-limiting/rate-limiting.libsonnet',
  FeatureRollout: import 'feature-rollout/base/feature-rollout.libsonnet',
  PromQLServiceProtection: import 'service-protection/promql/promql.libsonnet',
  ServiceProtectionAverageLatency: import 'service-protection/average-latency/average-latency.libsonnet',
  PodAutoScaler: import 'pod-auto-scaler/pod-auto-scaler.libsonnet',
  ServiceProtectionAverageLatencyWithPodAutoScaler: import 'service-protection-with-load-based-pod-auto-scaler/average-latency/average-latency.libsonnet',
  PromQLServiceProtectionWithPodAutoScaler: import 'service-protection-with-load-based-pod-auto-scaler/promql/promql.libsonnet',
}
