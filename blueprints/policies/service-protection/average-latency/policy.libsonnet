local spec = import '../../../spec.libsonnet';
local basePolicyFn = import '../base/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config.common + config.policy + cfg,

  local basePolicy = basePolicyFn(cfg).policyDef,

  // Add new components to basePolicy
  local policyDef = basePolicy {
    circuit+: {
      components+: [
        spec.v1.Component.withQuery(
          spec.v1.Query.new()
          + spec.v1.Query.withPromql(
            local q = 'sum(increase(flux_meter_sum{flow_status="OK", flux_meter_name="%(policy_name)s"}[5s]))/sum(increase(flux_meter_count{flow_status="OK", flux_meter_name="%(policy_name)s"}[5s]))' % { policy_name: params.policy_name };
            spec.v1.PromQL.new()
            + spec.v1.PromQL.withQueryString(q)
            + spec.v1.PromQL.withEvaluationInterval('1s')
            + spec.v1.PromQL.withOutPorts({ output: spec.v1.Port.withSignalName('SIGNAL') }),
          ),
        ),
        spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.mul(
          spec.v1.Port.withSignalName('SIGNAL'),
          spec.v1.Port.withConstantSignal(params.latency_baseliner.latency_ema_limit_multiplier),
          output=spec.v1.Port.withSignalName('MAX_EMA')
        )),
        spec.v1.Component.withEma(
          spec.v1.EMA.withParameters(params.latency_baseliner.ema)
          + spec.v1.EMA.withInPortsMixin(
            spec.v1.EMA.inPorts.withInput(spec.v1.Port.withSignalName('SIGNAL'))
            + spec.v1.EMA.inPorts.withMaxEnvelope(spec.v1.Port.withSignalName('MAX_EMA'))
          )
          + spec.v1.EMA.withOutPortsMixin(spec.v1.EMA.outPorts.withOutput(spec.v1.Port.withSignalName('SIGNAL_EMA')))
        ),
        spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.mul(
          spec.v1.Port.withSignalName('SIGNAL_EMA'),
          spec.v1.Port.withConstantSignal(params.latency_baseliner.latency_tolerance_multiplier),
          output=spec.v1.Port.withSignalName('SETPOINT')
        )),
      ],
    },
    resources+: {
      flow_control+: {
        flux_meters+: { [params.policy_name]: params.latency_baseliner.flux_meter },
      },
    },
  },

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,

  policyDef: policyDef,
}
