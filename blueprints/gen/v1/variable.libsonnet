local variableouts = import './variableouts.libsonnet';
{
  new():: {
  },
  outPorts:: variableouts,
  withConstantSignal(constant_signal):: {
    constant_signal: constant_signal,
  },
  withConstantSignalMixin(constant_signal):: {
    constant_signal+: constant_signal,
  },
  withConstantSignalConfigKey(constant_signal_config_key):: {
    constant_signal_config_key: constant_signal_config_key,
  },
  withConstantSignalConfigKeyMixin(constant_signal_config_key):: {
    constant_signal_config_key+: constant_signal_config_key,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
