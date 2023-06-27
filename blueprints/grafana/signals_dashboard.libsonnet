local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrap = import './utils/unwrap_panels.libsonnet';
local queryVariable = import './utils/variable_template.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';
local var = g.dashboard.variable;

function(cfg) {
  local config = defaultConfig + cfg,
  local dashboard = base(config),
  local signalsPanels = unwrap(std.toString('signals'), {}, config).panel,

  local signalNameVar = queryVariable('signal_name',
                                      'label_values(signal_reading{policy_name="%(policy_name)s"}, signal_name)' % { policy_name: cfg.policy.policy_name },
                                      '${datasource}',
                                      'Signal Name').variable,

  local subCircuitIDVar = queryVariable('sub_circuit_id',
                                        'label_values(signal_reading{policy_name="%(policy_name)s",signal_name="${signal_name}"}, sub_circuit_id)' % { policy_name: cfg.policy.policy_name },
                                        '${datasource}',
                                        'Sub Circuit ID').variable,

  local final = dashboard.baseDashboard
                + g.dashboard.withPanels(signalsPanels)
                + g.dashboard.withVariables([signalNameVar, subCircuitIDVar]),
  dashboard: final,
}
