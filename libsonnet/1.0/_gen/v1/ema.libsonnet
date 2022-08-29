local emains = import './emains.libsonnet';
local emaouts = import './emaouts.libsonnet';
{
  new():: {
    in_ports: {
      input: error 'Port input is missing',
      max_envelope: error 'Port max_envelope is missing',
      min_envelope: error 'Port min_envelope is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  inPorts:: emains,
  outPorts:: emaouts,
  withCorrectionFactorOnMaxEnvelopeViolation(correction_factor_on_max_envelope_violation):: {
    correction_factor_on_max_envelope_violation: correction_factor_on_max_envelope_violation,
  },
  withCorrectionFactorOnMaxEnvelopeViolationMixin(correction_factor_on_max_envelope_violation):: {
    correction_factor_on_max_envelope_violation+: correction_factor_on_max_envelope_violation,
  },
  withCorrectionFactorOnMinEnvelopeViolation(correction_factor_on_min_envelope_violation):: {
    correction_factor_on_min_envelope_violation: correction_factor_on_min_envelope_violation,
  },
  withCorrectionFactorOnMinEnvelopeViolationMixin(correction_factor_on_min_envelope_violation):: {
    correction_factor_on_min_envelope_violation+: correction_factor_on_min_envelope_violation,
  },
  withEmaWindow(ema_window):: {
    ema_window: ema_window,
  },
  withEmaWindowMixin(ema_window):: {
    ema_window+: ema_window,
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
  withWarmUpWindow(warm_up_window):: {
    warm_up_window: warm_up_window,
  },
  withWarmUpWindowMixin(warm_up_window):: {
    warm_up_window+: warm_up_window,
  },
}
