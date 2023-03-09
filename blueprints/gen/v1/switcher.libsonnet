local switcherins = import './switcherins.libsonnet';
local switcherouts = import './switcherouts.libsonnet';
{
  new():: {
  },
  inPorts:: switcherins,
  outPorts:: switcherouts,
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
