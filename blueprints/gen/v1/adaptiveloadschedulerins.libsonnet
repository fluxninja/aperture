{
  new():: {
  },
  withEnabled(enabled):: {
    enabled: enabled,
  },
  withEnabledMixin(enabled):: {
    enabled+: enabled,
  },
  withSetpoint(setpoint):: {
    setpoint: setpoint,
  },
  withSetpointMixin(setpoint):: {
    setpoint+: setpoint,
  },
  withSignal(signal):: {
    signal: signal,
  },
  withSignalMixin(signal):: {
    signal+: signal,
  },
}
