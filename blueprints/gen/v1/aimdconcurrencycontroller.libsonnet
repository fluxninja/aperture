local aimdconcurrencycontrollerins = import './aimdconcurrencycontrollerins.libsonnet';
local aimdconcurrencycontrollerouts = import './aimdconcurrencycontrollerouts.libsonnet';
{
  new():: {
    in_ports: {
      setpoint: error 'Port setpoint is missing',
      signal: error 'Port signal is missing',
    },
    out_ports: {
      accepted_concurrency: error 'Port accepted_concurrency is missing',
      desired_load_multiplier: error 'Port desired_load_multiplier is missing',
      incoming_concurrency: error 'Port incoming_concurrency is missing',
      is_overload: error 'Port is_overload is missing',
      observed_load_multiplier: error 'Port observed_load_multiplier is missing',
    },
  },
  inPorts:: aimdconcurrencycontrollerins,
  outPorts:: aimdconcurrencycontrollerouts,
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
