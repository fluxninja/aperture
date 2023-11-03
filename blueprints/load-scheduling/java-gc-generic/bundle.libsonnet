local utils = import '../common/utils.libsonnet';
local blueprint = import './java-gc.libsonnet';
local jmxUtils = import './utils.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params) {
  local c = std.mergePatch(config, params),
  local policyName = c.policy.policy_name,
  local promqlQuery = 'avg(java_lang_G1_Young_Generation_LastGcInfo_duration{policy_name="%(policy_name)s", infra_meter_name="jmx_inframeter"})' % { policy_name: policyName },

  local updated_cfg = utils.add_kubelet_overload_confirmations(c).updated_cfg {
    policy+: {
      promql_query: promqlQuery,
      setpoint: c.policy.load_scheduling_core.setpoint,
      overload_condition: 'gt',
    },
  },

  local infraMeters = if std.objectHas(c.policy.resources, 'infra_meters') then c.policy.resources.infra_meters else {},
  assert !std.objectHas(infraMeters, 'jmx_inframeter') : 'An infra meter with name jmx_inframeter already exists. Please choose a different name.',
  local config_with_jmx_infra_meter = updated_cfg {
    policy+: {
      resources+: {
        infra_meters+: jmxUtils(c),
      },
    },
  },
  local p = policy(config_with_jmx_infra_meter, params),
  policies: {
    [std.format('%s-cr.yaml', config_with_jmx_infra_meter.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', config_with_jmx_infra_meter.policy.policy_name)]: p.policyDef,
  },
}
