local loadshaperins = import './loadshaperins.libsonnet';
local loadshaperouts = import './loadshaperouts.libsonnet';
{
  new():: {
  },
  inPorts:: loadshaperins,
  outPorts:: loadshaperouts,
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
