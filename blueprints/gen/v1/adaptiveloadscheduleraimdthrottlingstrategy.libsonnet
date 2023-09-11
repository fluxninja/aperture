local adaptiveloadscheduleraimdthrottlingstrategyins = import './adaptiveloadscheduleraimdthrottlingstrategyins.libsonnet';
{
  new():: {
  },
  inPorts:: adaptiveloadscheduleraimdthrottlingstrategyins,
  withGradient(gradient):: {
    gradient: gradient,
  },
  withGradientMixin(gradient):: {
    gradient+: gradient,
  },
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withLoadMultiplierLinearIncrement(load_multiplier_linear_increment):: {
    load_multiplier_linear_increment: load_multiplier_linear_increment,
  },
  withLoadMultiplierLinearIncrementMixin(load_multiplier_linear_increment):: {
    load_multiplier_linear_increment+: load_multiplier_linear_increment,
  },
  withMaxLoadMultiplier(max_load_multiplier):: {
    max_load_multiplier: max_load_multiplier,
  },
  withMaxLoadMultiplierMixin(max_load_multiplier):: {
    max_load_multiplier+: max_load_multiplier,
  },
}
