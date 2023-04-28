local loadrampseriesloadrampinstanceouts = import './loadrampseriesloadrampinstanceouts.libsonnet';
{
  new():: {
  },
  outPorts:: loadrampseriesloadrampinstanceouts,
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
