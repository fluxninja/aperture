local concurrencyactuatorins = import './concurrencyactuatorins.libsonnet';
{
  new():: {
    in_ports: {
      desired_concurrency: error 'Port desired_concurrency is missing',
    },
  },
  inPorts:: concurrencyactuatorins,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
}
