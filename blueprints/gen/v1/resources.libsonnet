{
  new():: {
  },
  withFlowControl(flow_control):: {
    flow_control: flow_control,
  },
  withFlowControlMixin(flow_control):: {
    flow_control+: flow_control,
  },
  withInfraMeters(infra_meters):: {
    infra_meters: infra_meters,
  },
  withInfraMetersMixin(infra_meters):: {
    infra_meters+: infra_meters,
  },
  withTelemetryCollectors(telemetry_collectors):: {
    telemetry_collectors:
      if std.isArray(telemetry_collectors)
      then telemetry_collectors
      else [telemetry_collectors],
  },
  withTelemetryCollectorsMixin(telemetry_collectors):: {
    telemetry_collectors+: telemetry_collectors,
  },
}
