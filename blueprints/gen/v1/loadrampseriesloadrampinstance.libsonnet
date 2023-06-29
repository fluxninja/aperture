local loadrampouts = import './loadrampouts.libsonnet';
{
  new():: {
  },
  outPorts:: loadrampouts,
  withLoadRamp(load_ramp):: {
    load_ramp: load_ramp,
  },
  withLoadRampMixin(load_ramp):: {
    load_ramp+: load_ramp,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
