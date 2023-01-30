local emains = import './emains.libsonnet';
local emaouts = import './emaouts.libsonnet';
{
  new():: {
    in_ports: {
      input: error 'Port input is missing',
      max_envelope: error 'Port max_envelope is missing',
      min_envelope: error 'Port min_envelope is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: emains,
  outPorts:: emaouts,
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
