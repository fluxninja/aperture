local variableouts = import './variableouts.libsonnet';
{
  new():: {
  },
  outPorts:: variableouts,
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
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
