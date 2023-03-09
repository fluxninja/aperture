local andins = import './andins.libsonnet';
local andouts = import './andouts.libsonnet';
{
  new():: {
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
