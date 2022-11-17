local ratelimiterins = import './ratelimiterins.libsonnet';
{
  new():: {
    in_ports: {
      limit: error 'Port limit is missing',
    },
  },
  inPorts:: ratelimiterins,
  withDefaultConfig(default_config):: {
    default_config: default_config,
  },
  withDefaultConfigMixin(default_config):: {
    default_config+: default_config,
  },
  withDynamicConfigKey(dynamic_config_key):: {
    dynamic_config_key: dynamic_config_key,
  },
  withDynamicConfigKeyMixin(dynamic_config_key):: {
    dynamic_config_key+: dynamic_config_key,
  },
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
  withSelector(selector):: {
    selector: selector,
  },
  withSelectorMixin(selector):: {
    selector+: selector,
  },
}
