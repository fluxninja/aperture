{
  StaticRateLimiting: import 'static-rate-limiting/static-rate-limiting.libsonnet',
  RabbitMQQueueBuildup: import 'rabbitmq-queue-buildup/rabbitmq-queue-buildup.libsonnet',
  FeatureRollout: import 'feature-rollout/base/feature-rollout.libsonnet',
  ServiceProtectionAverageLatency: import 'service-protection/average-latency.libsonnet',
}
