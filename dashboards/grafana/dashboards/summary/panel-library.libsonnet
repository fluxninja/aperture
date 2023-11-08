{
  RateLimiter: import './rate-limiter/rows.libsonnet',
  PromQL: import './promql/rows.libsonnet',
  PodScaleActuator: import './pod-scale-actuator/rows.libsonnet',
  PodScaleReporter: import './pod-scale-reporter/rows.libsonnet',
  Sampler: import './sampler/rows.libsonnet',
  QuotaScheduler: import './quota-scheduler/rows.libsonnet',
  LoadActuator: import './load-actuator/rows.libsonnet',
}
