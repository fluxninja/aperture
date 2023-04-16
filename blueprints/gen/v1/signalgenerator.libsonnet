local signalgeneratorins = import './signalgeneratorins.libsonnet';
local signalgeneratorouts = import './signalgeneratorouts.libsonnet';
{
  new():: {
  },
  inPorts:: signalgeneratorins,
  outPorts:: signalgeneratorouts,
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
