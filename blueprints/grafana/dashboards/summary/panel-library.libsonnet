{
  RateLimiter: import './rate-limiter/panels.libsonnet',
  PromQL: import './promql/panels.libsonnet',
  PodScaleActuator: import './pod-scale-actuator/panels.libsonnet',
  PodScaleReporter: import './pod-scale-reporter/panels.libsonnet',
  Sampler: import './sampler/panels.libsonnet',
  QuotaScheduler: import './quota-scheduler/panels.libsonnet',
  LoadActuator: import './load-actuator/panels.libsonnet',
}
