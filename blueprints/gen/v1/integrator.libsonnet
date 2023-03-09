local integratorins = import './integratorins.libsonnet';
local integratorouts = import './integratorouts.libsonnet';
{
  new():: {
  },
  inPorts:: integratorins,
  outPorts:: integratorouts,
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
