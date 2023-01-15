local horizontalpodscalerscalereporterouts = import './horizontalpodscalerscalereporterouts.libsonnet';
{
  new():: {
    out_ports: {
      actual_replicas: error 'Port actual_replicas is missing',
      configured_replicas: error 'Port configured_replicas is missing',
    },
  },
  outPorts:: horizontalpodscalerscalereporterouts,
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
