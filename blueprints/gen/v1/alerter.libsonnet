local alerterins = import './alerterins.libsonnet';
{
  new():: {
    in_ports: {
      signal: error 'Port signal is missing',
    },
  },
  inPorts:: alerterins,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withParameters(parameters):: {
    parameters: parameters,
  },
  withParametersMixin(parameters):: {
    parameters+: parameters,
  },
}
