local polynomialrangefunctionins = import './polynomialrangefunctionins.libsonnet';
local polynomialrangefunctionouts = import './polynomialrangefunctionouts.libsonnet';
{
  new():: {
  },
  inPorts:: polynomialrangefunctionins,
  outPorts:: polynomialrangefunctionouts,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withParameters(parameters):: {
    parameters: parameters,
  },
  withParametersMixin(parameters):: {
    parameters+: parameters,
  },
}
