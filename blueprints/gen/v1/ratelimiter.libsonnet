local ratelimiterins = import './ratelimiterins.libsonnet';
{
  new():: {
  },
  inPorts:: ratelimiterins,
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
  withSelectors(selectors):: {
    selectors:
      if std.isArray(selectors)
      then selectors
      else [selectors],
  },
  withSelectorsMixin(selectors):: {
    selectors+: selectors,
  },
  withTokensLabelKey(tokens_label_key):: {
    tokens_label_key: tokens_label_key,
  },
  withTokensLabelKeyMixin(tokens_label_key):: {
    tokens_label_key+: tokens_label_key,
  },
}
