local spec = import '../../spec.libsonnet';
local commonPolicyFn = import '../common/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg, params={}, metadata={}) {
  local updatedConfig = config + cfg,

  local commonPolicy = commonPolicyFn(cfg, params, metadata),

  local createQuery = function(policy_name, interval) 'sum(increase(flux_meter_sum{flow_status="OK", flux_meter_name="%(policy_name)s", policy_name="%(policy_name)s"}[%(interval)s]))/sum(increase(flux_meter_count{flow_status="OK", flux_meter_name="%(policy_name)s", policy_name="%(policy_name)s"}[%(interval)s]))' % { policy_name: policy_name, interval: interval },

  // Add new components to commonPolicy
  local policyDef = commonPolicy.policyDef {
    circuit+: {
      components+: [
        spec.v1.Component.withQuery(
          spec.v1.Query.new()
          + spec.v1.Query.withPromql(
            local q = createQuery(updatedConfig.policy.policy_name, '30s');
            spec.v1.PromQL.new()
            + spec.v1.PromQL.withQueryString(q)
            + spec.v1.PromQL.withEvaluationInterval(evaluation_interval=updatedConfig.policy.evaluation_interval)
            + spec.v1.PromQL.withOutPorts({ output: spec.v1.Port.withSignalName('SIGNAL') }),
          ),
        ),
        spec.v1.Component.withQuery(
          spec.v1.Query.new()
          + spec.v1.Query.withPromql(
            local q = createQuery(updatedConfig.policy.policy_name, updatedConfig.policy.latency_baseliner.long_term_query_interval);
            spec.v1.PromQL.new()
            + spec.v1.PromQL.withQueryString(q)
            + spec.v1.PromQL.withEvaluationInterval(evaluation_interval=updatedConfig.policy.latency_baseliner.long_term_query_periodic_interval)
            + spec.v1.PromQL.withOutPorts({ output: spec.v1.Port.withSignalName('SIGNAL_LONG_TERM') }),
          ),
        ),
        spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.mul(
          spec.v1.Port.withSignalName('SIGNAL_LONG_TERM'),
          spec.v1.Port.withConstantSignal(updatedConfig.policy.latency_baseliner.latency_tolerance_multiplier),
          output=spec.v1.Port.withSignalName('SETPOINT')
        )),
      ],
    },
    resources+: {
      flow_control+: {
        flux_meters+: { [updatedConfig.policy.policy_name]: updatedConfig.policy.latency_baseliner.flux_meter },
      },
    },
  },

  policyResource: commonPolicy.policyResource {
    spec+: policyDef,
  },
  policyDef: policyDef,
}
