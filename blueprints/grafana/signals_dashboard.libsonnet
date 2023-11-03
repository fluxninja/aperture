local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrap = import './utils/unwrap_panels.libsonnet';
local queryVariable = import './utils/variable_template.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(componentsList, policyName, datasource, extraFilters={}) {
  local dashboard = base('Signals Dashboard - %s' % policyName, defaultConfig.dashboard.refresh_interval),
  local signalsPanels = unwrap(datasource, policyName, { component_name: 'Signals' }, false, extraFilters).panel,

  local signalNameVar = queryVariable('signal_name',
                                      'label_values(signal_reading{policy_name="%(policy_name)s"}, signal_name)' % { policy_name: policyName },
                                      '%(datasource)s' % { datasource: datasource },
                                      'Signal Name').variable,

  local subCircuitIDVar = queryVariable('sub_circuit_id',
                                        'label_values(signal_reading{policy_name="%(policy_name)s",signal_name="${signal_name}"}, sub_circuit_id)' % { policy_name: policyName },
                                        '%(datasource)s' % { datasource: datasource },
                                        'Sub Circuit ID').variable,

  local final = dashboard.baseDashboard
                + g.dashboard.withPanels(signalsPanels)
                + g.dashboard.withVariables([signalNameVar, subCircuitIDVar]),
  dashboard: final,
}
