local concurrencycontrollerins = import './concurrencycontrollerins.libsonnet';
local concurrencycontrollerouts = import './concurrencycontrollerouts.libsonnet';
{
  new():: {
    in_ports: {
      signal: error 'Port signal is missing',
    },
    out_ports: {
      accepted_concurrency: error 'Port accepted_concurrency is missing',
      desired_concurrency: error 'Port desired_concurrency is missing',
      incoming_concurrency: error 'Port incoming_concurrency is missing',
      is_overload: error 'Port is_overload is missing',
      load_multiplier: error 'Port load_multiplier is missing',
    },
  },
  inPorts:: concurrencycontrollerins,
  outPorts:: concurrencycontrollerouts,
  withAlerterConfig(alerter_config):: {
    alerter_config: alerter_config,
  },
  withAlerterConfigMixin(alerter_config):: {
    alerter_config+: alerter_config,
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
  withConcurrencySquareRootIncrementMultiplier(concurrency_square_root_increment_multiplier):: {
    concurrency_square_root_increment_multiplier: concurrency_square_root_increment_multiplier,
  },
  withConcurrencySquareRootIncrementMultiplierMixin(concurrency_square_root_increment_multiplier):: {
    concurrency_square_root_increment_multiplier+: concurrency_square_root_increment_multiplier,
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
