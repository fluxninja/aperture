local smains = import './smains.libsonnet';
local smaouts = import './smaouts.libsonnet';
{
  new():: {
    in_ports: {
      signal: error 'Port signal is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: smains,
  outPorts:: smaouts,
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
