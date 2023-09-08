local pidcontrollerins = import './pidcontrollerins.libsonnet';
local pidcontrollerouts = import './pidcontrollerouts.libsonnet';
{
  new():: {
  },
  inPorts:: pidcontrollerins,
  outPorts:: pidcontrollerouts,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
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
