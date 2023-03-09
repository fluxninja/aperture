local podscalerscalereporterouts = import './podscalerscalereporterouts.libsonnet';
{
  new():: {
  },
  outPorts:: podscalerscalereporterouts,
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
