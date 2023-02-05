local andins = import './andins.libsonnet';
local andouts = import './andouts.libsonnet';
{
  new():: {
    in_ports: {
      inputs: error 'Port inputs is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: andins,
  outPorts:: andouts,
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
