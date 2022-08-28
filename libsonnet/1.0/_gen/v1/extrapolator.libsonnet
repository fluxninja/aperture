local extrapolatorins = import './extrapolatorins.libsonnet';
local extrapolatorouts = import './extrapolatorouts.libsonnet';
{
  new():: {
    in_ports: {
      input: error 'Port input is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: extrapolatorins,
  outPorts:: extrapolatorouts,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withMaxExtrapolationInterval(max_extrapolation_interval):: {
    max_extrapolation_interval: max_extrapolation_interval,
  },
  withMaxExtrapolationIntervalMixin(max_extrapolation_interval):: {
    max_extrapolation_interval+: max_extrapolation_interval,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
