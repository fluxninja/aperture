{
  RateLimiting: import 'rate-limiting/base/rate-limiting.libsonnet',
  QuotaScheduling: import 'quota-scheduling/base/quota-scheduling.libsonnet',
  LoadRamping: import 'load-ramping/base/load-ramping.libsonnet',
  LoadSchedulingPromQL: import 'load-scheduling/promql/promql.libsonnet',
  LoadSchedulingPostgreSQL: import 'load-scheduling/postgresql/postgresql.libsonnet',
  LoadSchedulingAverageLatency: import 'load-scheduling/average-latency/average-latency.libsonnet',
  PodAutoScaler: import 'auto-scaling/pod-auto-scaler/pod-auto-scaler.libsonnet',
}
