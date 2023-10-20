{
  RateLimiter: import './panels/rate_limiter.libsonnet',
  PromQL: import './panels/query.libsonnet',

  // Grouped panels
  Actuator: import './panels/grouped/load_scheduler.libsonnet',
  AutoScale: import './panels/grouped/auto_scale.libsonnet',
  Sampler: import './panels/grouped/load_ramp.libsonnet',
  QuotaScheduler: import './panels/grouped/quota_scheduler.libsonnet',
}
