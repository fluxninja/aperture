local loadshaperseriesloadshaperinstanceouts = import './loadshaperseriesloadshaperinstanceouts.libsonnet';
{
  new():: {
  },
  outPorts:: loadshaperseriesloadshaperinstanceouts,
  withLoadShaper(load_shaper):: {
    load_shaper: load_shaper,
  },
  withLoadShaperMixin(load_shaper):: {
    load_shaper+: load_shaper,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
