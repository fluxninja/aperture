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
  withEnd(end):: {
    end: end,
  },
  withEndMixin(end):: {
    end+: end,
  },
  withLoadScheduler(load_scheduler):: {
    load_scheduler: load_scheduler,
  },
  withLoadSchedulerMixin(load_scheduler):: {
    load_scheduler+: load_scheduler,
  },
  withStart(start):: {
    start: start,
  },
  withStartMixin(start):: {
    start+: start,
  },
}
