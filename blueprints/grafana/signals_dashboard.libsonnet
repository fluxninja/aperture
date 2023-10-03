local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local queryVariable = import './utils/variable_template.libsonnet';
local signals = import './panels/grouped/signals.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, datasource, extraFilters={}) {
  local dashboard = base('Aperture Signals - %s' % policyName, defaultConfig.dashboard.refresh_interval),
  local signalsPanels = signals(datasource, policyName, {}, extraFilters).panels,

  local signalNameVar = queryVariable('signal_name',
                                      'label_values(signal_reading{policy_name="%(policy_name)s"}, signal_name)' % { policy_name: policyName },
                                      '${datasource}',
                                      'Signal Name').variable,

  local subCircuitIDVar = queryVariable('sub_circuit_id',
                                        'label_values(signal_reading{policy_name="%(policy_name)s",signal_name="${signal_name}"}, sub_circuit_id)' % { policy_name: policyName },
                                        '${datasource}',
                                        'Sub Circuit ID').variable,

  local final = dashboard.baseDashboard
                + g.dashboard.withPanels(signalsPanels)
                + g.dashboard.withVariables([signalNameVar, subCircuitIDVar]),
  dashboard: final,
}
