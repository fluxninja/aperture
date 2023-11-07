local timeSeriesPanel = import '../panels/time-series.libsonnet';
local promUtils = import '../utils/prometheus.libsonnet';
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

{
  getInPort(component, portName)::
    if 'in_ports' in component then self._findPort(component.in_ports, portName) else null,

  getOutPort(component, portName)::
    if 'out_ports' in component then self._findPort(component.out_ports, portName) else null,

  getPort(inOrOutPorts, portName)::
    self._findPort(inOrOutPorts, portName),

  _findPort(ports, portName):: (
    local matchingPorts = [port for port in ports if port.port_name == portName];
    if std.length(matchingPorts) > 0 then matchingPorts[0] else null
  ),

  processInPort(datasourceName, component, portName, policyName, extraFilters):: (
    local inPort = self.getInPort(component, portName);
    if inPort == null
    then
      {
        targets: [],
        port: null,
      }
    else
      local targets = self._targetsForPort(datasourceName, component.in_ports, portName, component.component_name, policyName, extraFilters);
      {
        targets: targets,
        port: inPort,
      }
  ),

  targetsForInPort(datasourceName, component, portName, policyName, extraFilters):: (
    local processedPort = self.processInPort(datasourceName, component, portName, policyName, extraFilters);
    processedPort.targets
  ),

  processOutPort(datasourceName, component, portName, policyName, extraFilters):: (
    local outPort = self.getOutPort(component, portName);
    if outPort == null
    then
      {
        targets: [],
        port: null,
      }
    else
      local targets = self._targetsForPort(datasourceName, component.out_ports, portName, component.component_name, policyName, extraFilters);
      {
        targets: targets,
        port: outPort,
      }
  ),

  targetsForOutPort(datasourceName, component, portName, policyName, extraFilters):: (
    local processedPort = self.processOutPort(datasourceName, component, portName, policyName, extraFilters);
    processedPort.targets
  ),

  _targetsForPort(datasourceName, ports, portName, componentName, policyName, extraFilters):: (
    local port = self._findPort(ports, portName);
    if port == null
    then
      []
    else if 'signal_name' in port
    then
      local signalName = port.signal_name;
      local subCircuitId = port.sub_circuit_id;
      local signalFilters = extraFilters { policy_name: policyName, signal_name: signalName, sub_circuit_id: subCircuitId };
      local stringFilters = promUtils.dictToPrometheusFilter(signalFilters);
      [
        g.query.prometheus.new(datasourceName, 'increase(signal_reading_sum{%(filters)s}[$__rate_interval]) / increase(signal_reading_count{%(filters)s}[$__rate_interval])' % { filters: stringFilters })
        + g.query.prometheus.withIntervalFactor(1)
        + g.query.prometheus.withLegendFormat('Signal at component:%(componentName)s, port:%(portName)s' % { portName: portName, componentName: componentName }),
      ]
    else
      [
        {
          refId: 'A',
          type: 'timeseries',
          expr: port.constant_value,
          legendFormat: 'Constant at component:%(componentName)s, port:%(portName)s' % { portName: portName, componentName: componentName },
        },
      ]
  ),

  panelsForInPort(title, datasourceName, component, portName, policyName, extraFilters, x=0, h=10, w=24, description=''):: (
    local processedPort = self.processInPort(datasourceName, component, portName, policyName, extraFilters);
    self._panelsForTargets(processedPort.port, title, datasourceName, processedPort.targets, x, h, w, description)
  ),

  panelsForOutPort(title, datasourceName, component, portName, policyName, extraFilters, x=0, h=10, w=24, description=''):: (
    local processedPort = self.processOutPort(datasourceName, component, portName, policyName, extraFilters);
    self._panelsForTargets(processedPort.port, title, datasourceName, processedPort.targets, x, h, w, description)
  ),

  _panelsForTargets(port, title, datasourceName, targets, x, h, w, description):: (
    if std.length(targets) == 0
    then
      []
    else
      local processedTitle =
        if port == null
        then
          title
        else
          if 'signal_name' in port
          then
            '%(signalName)s - %(title)s' % { title: title, signalName: port.signal_name }
          else
            '%(title)s (constant:%(constantValue)s)' % { title: title, constantValue: port.constant_value };
      [timeSeriesPanel(
        processedTitle,
        datasourceName,
        targets=targets,
        x=x,
        h=h,
        w=w,
        description=description,
      )]
  ),
}
