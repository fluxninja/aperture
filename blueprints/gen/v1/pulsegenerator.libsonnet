local pulsegeneratorouts = import './pulsegeneratorouts.libsonnet';
{
  new():: {
  },
  outPorts:: pulsegeneratorouts,
  withFalseFor(false_for):: {
    false_for: false_for,
  },
  withFalseForMixin(false_for):: {
    false_for+: false_for,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withTrueFor(true_for):: {
    true_for: true_for,
  },
  withTrueForMixin(true_for):: {
    true_for+: true_for,
  },
}
