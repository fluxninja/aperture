local orins = import './orins.libsonnet';
local orouts = import './orouts.libsonnet';
{
  new():: {
  },
  inPorts:: orins,
  outPorts:: orouts,
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
