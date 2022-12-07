{
  new():: {
  },
  withAlertChannels(alert_channels):: {
    alert_channels:
      if std.isArray(alert_channels)
      then alert_channels
      else [alert_channels],
  },
  withAlertChannelsMixin(alert_channels):: {
    alert_channels+: alert_channels,
  },
  withAlertName(alert_name):: {
    alert_name: alert_name,
  },
  withAlertNameMixin(alert_name):: {
    alert_name+: alert_name,
  },
  withResolveTimeout(resolve_timeout):: {
    resolve_timeout: resolve_timeout,
  },
  withResolveTimeoutMixin(resolve_timeout):: {
    resolve_timeout+: resolve_timeout,
  },
  withSeverity(severity):: {
    severity: severity,
  },
  withSeverityMixin(severity):: {
    severity+: severity,
  },
}
