{
  StaticRateLimiting: import 'static-rate-limiting/static-rate-limiting.libsonnet',
  RabbitMQQueueBuildup: import 'service-protection/average-latency/service-protection.libsonnet',
  FeatureRollout: import 'feature-rollout/base/feature-rollout.libsonnet',
  ServiceProtectionAverageLatency: import 'service-protection/average-latency/service-protection.libsonnet',
}
