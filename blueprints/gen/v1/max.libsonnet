local maxins = import './maxins.libsonnet';
local maxouts = import './maxouts.libsonnet';
{
  new():: {
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
