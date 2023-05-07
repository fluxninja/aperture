{
  new():: {
  },
  withFlowControl(flow_control):: {
    flow_control: flow_control,
  },
  withFlowControlMixin(flow_control):: {
    flow_control+: flow_control,
  },
  withTelemetryCollector(telemetry_collector):: {
    telemetry_collector:
      if std.isArray(telemetry_collector)
      then telemetry_collector
      else [telemetry_collector],
  },
  withTelemetryCollectorMixin(telemetry_collector):: {
    telemetry_collector+: telemetry_collector,
  },
}
