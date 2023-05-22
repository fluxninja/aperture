local variableouts = import './variableouts.libsonnet';
{
  new():: {
  },
  outPorts:: variableouts,
  withConfigKey(config_key):: {
    config_key: config_key,
  },
  withConfigKeyMixin(config_key):: {
    config_key+: config_key,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withValue(value):: {
    value: value,
  },
  withValueMixin(value):: {
    value+: value,
  },
}
