local loadactuatorins = import './loadactuatorins.libsonnet';
{
  new():: {
    in_ports: {
      load_multiplier: error 'Port load_multiplier is missing',
    },
  },
  inPorts:: loadactuatorins,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
}
