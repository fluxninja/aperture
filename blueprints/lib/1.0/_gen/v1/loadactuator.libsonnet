local loadactuatorins = import './loadactuatorins.libsonnet';
{
  new():: {
    in_ports: {
      load_multiplier: error 'Port load_multiplier is missing',
    },
  },
  inPorts:: loadactuatorins,
  withAlerterConfig(alerter_config):: {
    alerter_config: alerter_config,
  },
  withAlerterConfigMixin(alerter_config):: {
    alerter_config+: alerter_config,
  },
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
}
