local switcherins = import './switcherins.libsonnet';
local switcherouts = import './switcherouts.libsonnet';
{
  new():: {
    in_ports: {
      off_signal: error 'Port off_signal is missing',
      on_signal: error 'Port on_signal is missing',
      switch: error 'Port switch is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: switcherins,
  outPorts:: switcherouts,
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
