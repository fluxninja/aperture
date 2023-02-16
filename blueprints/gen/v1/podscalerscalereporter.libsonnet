local podscalerscalereporterouts = import './podscalerscalereporterouts.libsonnet';
{
  new():: {
    out_ports: {
      actual_replicas: error 'Port actual_replicas is missing',
      configured_replicas: error 'Port configured_replicas is missing',
    },
  },
  outPorts:: podscalerscalereporterouts,
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
