{
  RateLimiter: import './panels/rate_limiter.libsonnet',
  PromQL: import './panels/query.libsonnet',

  // Grouped panels
  AdaptiveLoadScheduler: import './panels/grouped/load_scheduler.libsonnet',
  AIMDLoadScheduler: import './panels/grouped/load_scheduler.libsonnet',
  RangeDrivenLoadScheduler: import './panels/grouped/load_scheduler.libsonnet',
  AIADLoadScheduler: import './panels/grouped/load_scheduler.libsonnet',
  AutoScale: import './panels/grouped/auto_scale.libsonnet',
  LoadRamp: import './panels/grouped/load_ramp.libsonnet',
  QuotaScheduler: import './panels/grouped/quota_scheduler.libsonnet',
}
