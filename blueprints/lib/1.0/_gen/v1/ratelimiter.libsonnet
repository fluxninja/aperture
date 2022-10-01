local ratelimiterins = import './ratelimiterins.libsonnet';
{
  new():: {
    in_ports: {
      limit: error 'Port limit is missing',
    },
  },
  inPorts:: ratelimiterins,
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
  withInitConfig(init_config):: {
    init_config: init_config,
  },
  withInitConfigMixin(init_config):: {
    init_config+: init_config,
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
