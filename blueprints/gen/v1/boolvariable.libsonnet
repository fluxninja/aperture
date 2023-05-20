local boolvariableouts = import './boolvariableouts.libsonnet';
{
  new():: {
  },
  outPorts:: boolvariableouts,
  withConstantBool(constant_bool):: {
    constant_bool: constant_bool,
  },
  withConstantBoolMixin(constant_bool):: {
    constant_bool+: constant_bool,
  },
  withConstantBoolConfigKey(constant_bool_config_key):: {
    constant_bool_config_key: constant_bool_config_key,
  },
  withConstantBoolConfigKeyMixin(constant_bool_config_key):: {
    constant_bool_config_key+: constant_bool_config_key,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
