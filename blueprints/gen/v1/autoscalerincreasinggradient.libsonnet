local autoscalerincreasinggradientins = import './autoscalerincreasinggradientins.libsonnet';
{
  new():: {
    in_ports: {
      setpoint: error 'Port setpoint is missing',
      signal: error 'Port signal is missing',
    },
  },
  inPorts:: autoscalerincreasinggradientins,
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
