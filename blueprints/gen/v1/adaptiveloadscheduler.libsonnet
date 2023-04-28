local adaptiveloadschedulerins = import './adaptiveloadschedulerins.libsonnet';
local adaptiveloadschedulerouts = import './adaptiveloadschedulerouts.libsonnet';
{
  new():: {
  },
  inPorts:: adaptiveloadschedulerins,
  outPorts:: adaptiveloadschedulerouts,
  withAlerterParameters(alerter_parameters):: {
    alerter_parameters: alerter_parameters,
  },
  withAlerterParametersMixin(alerter_parameters):: {
    alerter_parameters+: alerter_parameters,
  },
  withDefaultConfig(default_config):: {
    default_config: default_config,
  },
  withDefaultConfigMixin(default_config):: {
    default_config+: default_config,
  },
  withDynamicConfigKey(dynamic_config_key):: {
    dynamic_config_key: dynamic_config_key,
  },
  withDynamicConfigKeyMixin(dynamic_config_key):: {
    dynamic_config_key+: dynamic_config_key,
  },
  withFlowSelector(flow_selector):: {
    flow_selector: flow_selector,
  },
  withFlowSelectorMixin(flow_selector):: {
    flow_selector+: flow_selector,
  },
  withGradientParameters(gradient_parameters):: {
    gradient_parameters: gradient_parameters,
  },
  withGradientParametersMixin(gradient_parameters):: {
    gradient_parameters+: gradient_parameters,
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
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withSchedulerParameters(scheduler_parameters):: {
    scheduler_parameters: scheduler_parameters,
  },
  withSchedulerParametersMixin(scheduler_parameters):: {
    scheduler_parameters+: scheduler_parameters,
  },
}
