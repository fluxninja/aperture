{
  RateLimiting: import 'rate-limiting/rate-limiting.libsonnet',
  QuotaScheduler: import 'quota-scheduler/quota-scheduler.libsonnet',
  FeatureRolloutAverageLatency: import 'feature-rollout/average-latency/average-latency.libsonnet',
  FeatureRolloutPercentileLatency: import 'feature-rollout/percentile-latency/percentile-latency.libsonnet',
  FeatureRolloutPromQL: import 'feature-rollout/promql/promql.libsonnet',
  ServiceProtectionPromQL: import 'service-protection/promql/promql.libsonnet',
  ServiceProtectionPostgreSQL: import 'service-protection/postgresql/postgresql.libsonnet',
  ServiceProtectionAverageLatency: import 'service-protection/average-latency/average-latency.libsonnet',
  PodAutoScaler: import 'auto-scaling/pod-auto-scaler/pod-auto-scaler.libsonnet',
}
