local base = import '../../dashboard/base.libsonnet';
local defaultConfig = import '../../dashboard/default-config.libsonnet';
local queryVariable = import '../../dashboard/variable-template.libsonnet';
local panels = import './panels.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(policyName, datasource, extraFilters={})
  local dashboard = base('Signals Dashboard - %s' % policyName, defaultConfig.refresh_interval);
  local signalsPanels = panels(datasource, policyName, extraFilters);

  local signalNameVar = queryVariable('signal_name',
                                      'label_values(signal_reading{policy_name="%(policy_name)s"}, signal_name)' % { policy_name: policyName },
                                      '%(datasource)s' % { datasource: datasource },
                                      'Signal Name').variable;

  local subCircuitIDVar = queryVariable('sub_circuit_id',
                                        'label_values(signal_reading{policy_name="%(policy_name)s",signal_name="${signal_name}"}, sub_circuit_id)' % { policy_name: policyName },
                                        '%(datasource)s' % { datasource: datasource },
                                        'Sub Circuit ID').variable;

  dashboard.baseDashboard
  + g.dashboard.withPanels(signalsPanels)
  + g.dashboard.withVariables([signalNameVar, subCircuitIDVar])
