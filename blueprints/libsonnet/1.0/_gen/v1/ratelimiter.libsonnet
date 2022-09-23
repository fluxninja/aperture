local ratelimiterins = import './ratelimiterins.libsonnet';
{
  new():: {
    in_ports: {
      limit: error 'Port limit is missing',
    },
  },
  inPorts:: ratelimiterins,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withLabelKey(label_key):: {
    label_key: label_key,
  },
  withLabelKeyMixin(label_key):: {
    label_key+: label_key,
  },
  withLazySync(lazy_sync):: {
    lazy_sync: lazy_sync,
  },
  withLazySyncMixin(lazy_sync):: {
    lazy_sync+: lazy_sync,
  },
  withLimitResetInterval(limit_reset_interval):: {
    limit_reset_interval: limit_reset_interval,
  },
  withLimitResetIntervalMixin(limit_reset_interval):: {
    limit_reset_interval+: limit_reset_interval,
  },
  withOverrides(overrides):: {
    overrides:
      if std.isArray(overrides)
      then overrides
      else [overrides],
  },
  withOverridesMixin(overrides):: {
    overrides+: overrides,
  },
  withSelector(selector):: {
    selector: selector,
  },
  withSelectorMixin(selector):: {
    selector+: selector,
  },
}
