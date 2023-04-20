local firstvalidins = import './firstvalidins.libsonnet';
local firstvalidouts = import './firstvalidouts.libsonnet';
{
  new():: {
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
