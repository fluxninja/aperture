{
  RateLimiting: import 'rate-limiting/base/rate-limiting.libsonnet',
  QuotaScheduler: import 'quota-scheduler/base/quota-scheduler.libsonnet',
  FeatureRollout: import 'feature-rollout/base/feature-rollout.libsonnet',
  ServiceProtectionPromQL: import 'service-protection/promql/promql.libsonnet',
  ServiceProtectionPostgreSQL: import 'service-protection/postgresql/postgresql.libsonnet',
  ServiceProtectionAverageLatency: import 'service-protection/average-latency/average-latency.libsonnet',
  PodAutoScaler: import 'auto-scaling/pod-auto-scaler/pod-auto-scaler.libsonnet',
}
