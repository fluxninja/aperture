local concurrencylimiterins = import './concurrencylimiterins.libsonnet';
local concurrencylimiterouts = import './concurrencylimiterouts.libsonnet';
{
  new():: {
  },
  inPorts:: concurrencylimiterins,
  outPorts:: concurrencylimiterouts,
  withConcurrencyLimiter(concurrency_limiter):: {
    concurrency_limiter: concurrency_limiter,
  },
  withConcurrencyLimiterMixin(concurrency_limiter):: {
    concurrency_limiter+: concurrency_limiter,
  },
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
