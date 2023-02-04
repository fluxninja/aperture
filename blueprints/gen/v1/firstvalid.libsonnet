local firstvalidins = import './firstvalidins.libsonnet';
local firstvalidouts = import './firstvalidouts.libsonnet';
{
  new():: {
    in_ports: {
      inputs: error 'Port inputs is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: firstvalidins,
  outPorts:: firstvalidouts,
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
