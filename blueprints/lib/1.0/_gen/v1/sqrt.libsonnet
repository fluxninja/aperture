local sqrtins = import './sqrtins.libsonnet';
local sqrtouts = import './sqrtouts.libsonnet';
{
  new():: {
    in_ports: {
      input: error 'Port input is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: sqrtins,
  outPorts:: sqrtouts,
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
  withScale(scale):: {
    scale: scale,
  },
  withScaleMixin(scale):: {
    scale+: scale,
  },
}
