local quotaschedulerins = import './quotaschedulerins.libsonnet';
{
  new():: {
  },
  inPorts:: quotaschedulerins,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withRateLimiter(rate_limiter):: {
    rate_limiter: rate_limiter,
  },
  withRateLimiterMixin(rate_limiter):: {
    rate_limiter+: rate_limiter,
  },
  withScheduler(scheduler):: {
    scheduler: scheduler,
  },
  withSchedulerMixin(scheduler):: {
    scheduler+: scheduler,
  },
  withSelectors(selectors):: {
    selectors:
      if std.isArray(selectors)
      then selectors
      else [selectors],
  },
  withSelectorsMixin(selectors):: {
    selectors+: selectors,
  },
}
