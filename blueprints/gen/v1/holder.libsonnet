local holderins = import './holderins.libsonnet';
local holderouts = import './holderouts.libsonnet';
{
  new():: {
  },
  inPorts:: holderins,
  outPorts:: holderouts,
  withHoldFor(hold_for):: {
    hold_for: hold_for,
  },
  withHoldForMixin(hold_for):: {
    hold_for+: hold_for,
  },
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
