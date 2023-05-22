local gradientcontrollerins = import './gradientcontrollerins.libsonnet';
local gradientcontrollerouts = import './gradientcontrollerouts.libsonnet';
{
  new():: {
  },
  inPorts:: gradientcontrollerins,
  outPorts:: gradientcontrollerouts,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withManualMode(manual_mode):: {
    manual_mode: manual_mode,
  },
  withManualModeMixin(manual_mode):: {
    manual_mode+: manual_mode,
  },
  withManualModeConfigKey(manual_mode_config_key):: {
    manual_mode_config_key: manual_mode_config_key,
  },
  withManualModeConfigKeyMixin(manual_mode_config_key):: {
    manual_mode_config_key+: manual_mode_config_key,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withParameters(parameters):: {
    parameters: parameters,
  },
  withParametersMixin(parameters):: {
    parameters+: parameters,
  },
}
