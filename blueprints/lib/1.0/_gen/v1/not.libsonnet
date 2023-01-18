local notins = import './notins.libsonnet';
local notouts = import './notouts.libsonnet';
{
  new():: {
    in_ports: {
      input: error 'Port input is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: notins,
  outPorts:: notouts,
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
