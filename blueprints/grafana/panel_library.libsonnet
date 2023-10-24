{
  RateLimiter: import './panels/rate_limiter.libsonnet',
  PromQL: import './panels/query.libsonnet',

  // Grouped panels
  AutoScale: import './panels/grouped/auto_scale.libsonnet',
  Sampler: import './panels/grouped/load_ramp.libsonnet',
  QuotaScheduler: import './panels/grouped/quota_scheduler.libsonnet',
  Signals: import './panels/grouped/signals.libsonnet',
}
