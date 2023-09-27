{
  new():: {
  },
  withAlerter(alerter):: {
    alerter: alerter,
  },
  withAlerterMixin(alerter):: {
    alerter+: alerter,
  },
  withDegree(degree):: {
    degree: degree,
  },
  withDegreeMixin(degree):: {
    degree+: degree,
  },
  withHighThrottleThreshold(high_throttle_threshold):: {
    high_throttle_threshold: high_throttle_threshold,
  },
  withHighThrottleThresholdMixin(high_throttle_threshold):: {
    high_throttle_threshold+: high_throttle_threshold,
  },
  withLoadScheduler(load_scheduler):: {
    load_scheduler: load_scheduler,
  },
  withLoadSchedulerMixin(load_scheduler):: {
    load_scheduler+: load_scheduler,
  },
  withLowThrottleThreshold(low_throttle_threshold):: {
    low_throttle_threshold: low_throttle_threshold,
  },
  withLowThrottleThresholdMixin(low_throttle_threshold):: {
    low_throttle_threshold+: low_throttle_threshold,
  },
}
