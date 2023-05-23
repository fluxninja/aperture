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
  withConstantOutput(constant_output):: {
    constant_output: constant_output,
  },
  withConstantOutputMixin(constant_output):: {
    constant_output+: constant_output,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
