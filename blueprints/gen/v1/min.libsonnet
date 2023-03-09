local minins = import './minins.libsonnet';
local minouts = import './minouts.libsonnet';
{
  new():: {
  },
  inPorts:: minins,
  outPorts:: minouts,
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
