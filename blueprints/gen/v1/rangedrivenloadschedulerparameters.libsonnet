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
  withHighWatermark(high_watermark):: {
    high_watermark: high_watermark,
  },
  withHighWatermarkMixin(high_watermark):: {
    high_watermark+: high_watermark,
  },
  withLoadScheduler(load_scheduler):: {
    load_scheduler: load_scheduler,
  },
  withLoadSchedulerMixin(load_scheduler):: {
    load_scheduler+: load_scheduler,
  },
  withLowWatermark(low_watermark):: {
    low_watermark: low_watermark,
  },
  withLowWatermarkMixin(low_watermark):: {
    low_watermark+: low_watermark,
  },
}
