local integratorins = import './integratorins.libsonnet';
local integratorouts = import './integratorouts.libsonnet';
{
  new():: {
    in_ports: {
      input: error 'Port input is missing',
      max: error 'Port max is missing',
      min: error 'Port min is missing',
      reset: error 'Port reset is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: integratorins,
  outPorts:: integratorouts,
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
}
