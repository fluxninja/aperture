local spec = import '../../../spec.libsonnet';
local basePolicyFn = import '../base/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg, params={}, metadata={}) {
  local updatedConfig = config + cfg,

  local basePolicy = basePolicyFn(cfg, params, metadata),

  // Add new components to basePolicy
  local policyDef = basePolicy.policyDef {
    circuit+: {
      components+: [
        spec.v1.Component.withQuery(
          spec.v1.Query.new()
          + spec.v1.Query.withPromql(
            local q = 'sum(increase(flux_meter_sum{flow_status="OK", flux_meter_name="%(policy_name)s", policy_name="%(policy_name)s"}[30s]))/sum(increase(flux_meter_count{flow_status="OK", flux_meter_name="%(policy_name)s", policy_name="%(policy_name)s"}[30s]))' % { policy_name: updatedConfig.policy.policy_name };
            spec.v1.PromQL.new()
            + spec.v1.PromQL.withQueryString(q)
            + spec.v1.PromQL.withEvaluationInterval(evaluation_interval=updatedConfig.policy.evaluation_interval)
            + spec.v1.PromQL.withOutPorts({ output: spec.v1.Port.withSignalName('SIGNAL') }),
          ),
        ),
        spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.mul(
          spec.v1.Port.withSignalName('SIGNAL'),
          spec.v1.Port.withConstantSignal(updatedConfig.policy.latency_baseliner.latency_ema_limit_multiplier),
          output=spec.v1.Port.withSignalName('MAX_EMA')
        )),
        spec.v1.Component.withEma(
          spec.v1.EMA.withParameters(updatedConfig.policy.latency_baseliner.ema)
          + spec.v1.EMA.withInPortsMixin(
            spec.v1.EMA.inPorts.withInput(spec.v1.Port.withSignalName('SIGNAL'))
            + spec.v1.EMA.inPorts.withMaxEnvelope(spec.v1.Port.withSignalName('MAX_EMA'))
          )
          + spec.v1.EMA.withOutPortsMixin(spec.v1.EMA.outPorts.withOutput(spec.v1.Port.withSignalName('SIGNAL_EMA')))
        ),
        spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.mul(
          spec.v1.Port.withSignalName('SIGNAL_EMA'),
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

  policyResource: basePolicy.policyResource {
    spec+: policyDef,
  },
  policyDef: policyDef,
}
