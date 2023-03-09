local inverterins = import './inverterins.libsonnet';
local inverterouts = import './inverterouts.libsonnet';
{
  new():: {
  },
  inPorts:: inverterins,
  outPorts:: inverterouts,
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
