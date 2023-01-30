local aimdconcurrencycontrollerins = import './aimdconcurrencycontrollerins.libsonnet';
local aimdconcurrencycontrollerouts = import './aimdconcurrencycontrollerouts.libsonnet';
{
  new():: {
    in_ports: {
      setpoint: error 'Port setpoint is missing',
      signal: error 'Port signal is missing',
    },
    out_ports: {
      is_overload: error 'Port is_overload is missing',
      load_multiplier: error 'Port load_multiplier is missing',
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
  withConcurrencyLimitMultiplier(concurrency_limit_multiplier):: {
    concurrency_limit_multiplier: concurrency_limit_multiplier,
  },
  withConcurrencyLimitMultiplierMixin(concurrency_limit_multiplier):: {
    concurrency_limit_multiplier+: concurrency_limit_multiplier,
  },
  withConcurrencyLinearIncrement(concurrency_linear_increment):: {
    concurrency_linear_increment: concurrency_linear_increment,
  },
  withConcurrencyLinearIncrementMixin(concurrency_linear_increment):: {
    concurrency_linear_increment+: concurrency_linear_increment,
  },
  withConcurrencySqrtIncrementMultiplier(concurrency_sqrt_increment_multiplier):: {
    concurrency_sqrt_increment_multiplier: concurrency_sqrt_increment_multiplier,
  },
  withConcurrencySqrtIncrementMultiplierMixin(concurrency_sqrt_increment_multiplier):: {
    concurrency_sqrt_increment_multiplier+: concurrency_sqrt_increment_multiplier,
  },
  withDryRunDynamicConfigKey(dry_run_dynamic_config_key):: {
    dry_run_dynamic_config_key: dry_run_dynamic_config_key,
  },
  withDryRunDynamicConfigKeyMixin(dry_run_dynamic_config_key):: {
    dry_run_dynamic_config_key+: dry_run_dynamic_config_key,
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
