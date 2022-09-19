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
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withNegativeFor(negative_for):: {
    negative_for: negative_for,
  },
  withNegativeForMixin(negative_for):: {
    negative_for+: negative_for,
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
  withPositiveFor(positive_for):: {
    positive_for: positive_for,
  },
  withPositiveForMixin(positive_for):: {
    positive_for+: positive_for,
  },
}
