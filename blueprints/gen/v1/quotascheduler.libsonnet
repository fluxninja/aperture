local quotaschedulerouts = import './quotaschedulerouts.libsonnet';
local ratelimiterins = import './ratelimiterins.libsonnet';
{
  new():: {
  },
  inPorts:: ratelimiterins,
  outPorts:: quotaschedulerouts,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
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
