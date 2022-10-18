local sinkins = import './sinkins.libsonnet';
{
  new():: {
    in_ports: {
      inputs: error 'Port inputs is missing',
    },
  },
  inPorts:: sinkins,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
}
