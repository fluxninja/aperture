local loadshedactuatorins = import './loadshedactuatorins.libsonnet';
{
  new():: {
    in_ports: {
      load_shed_factor: error 'Port load_shed_factor is missing',
    },
  },
  inPorts:: loadshedactuatorins,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
}
