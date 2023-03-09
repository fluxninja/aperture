local differentiatorins = import './differentiatorins.libsonnet';
local differentiatorouts = import './differentiatorouts.libsonnet';
{
  new():: {
  },
  inPorts:: differentiatorins,
  outPorts:: differentiatorouts,
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
  withWindow(window):: {
    window: window,
  },
  withWindowMixin(window):: {
    window+: window,
  },
}
