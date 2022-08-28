local maxins = import './maxins.libsonnet';
local maxouts = import './maxouts.libsonnet';
{
  new():: {
    in_ports: {
      inputs: error 'Port inputs is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: maxins,
  outPorts:: maxouts,
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
}
