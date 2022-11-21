local alerterins = import './alerterins.libsonnet';
{
  new():: {
    in_ports: {
      alert: error 'Port alert is missing',
    },
  },
  inPorts:: alerterins,
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
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
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
