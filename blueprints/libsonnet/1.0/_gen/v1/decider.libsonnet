local deciderins = import './deciderins.libsonnet';
local deciderouts = import './deciderouts.libsonnet';
{
  new():: {
    in_ports: {
      lhs: error 'Port lhs is missing',
      rhs: error 'Port rhs is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: deciderins,
  outPorts:: deciderouts,
  withFalseFor(false_for):: {
    false_for: false_for,
  },
  withFalseForMixin(false_for):: {
    false_for+: false_for,
  },
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withOperator(operator):: {
    operator: operator,
  },
  withOperatorMixin(operator):: {
    operator+: operator,
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
