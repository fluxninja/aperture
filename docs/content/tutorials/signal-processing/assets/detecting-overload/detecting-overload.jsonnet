local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local policy = aperture.spec.v1.Policy;
local component = aperture.spec.v1.Component;
local query = aperture.spec.v1.Query;
local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local circuit = aperture.spec.v1.Circuit;
local port = aperture.spec.v1.Port;
local resources = aperture.spec.v1.Resources;
local flowControlResources = aperture.spec.v1.FlowControlResources;
local fluxMeter = aperture.spec.v1.FluxMeter;
local promQL = aperture.spec.v1.PromQL;
local ema = aperture.spec.v1.EMA;
local emaParameters = aperture.spec.v1.EMAParameters;
local combinator = aperture.spec.v1.ArithmeticCombinator;
local decider = aperture.spec.v1.Decider;
local alerter = aperture.spec.v1.Alerter;
local alerterParameters = aperture.spec.v1.AlerterParameters;
local constantSignal = aperture.spec.v1.ConstantSignal;

local svcSelector =
  flowSelector.new()
  + flowSelector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + flowSelector.withFlowMatcher(
    flowMatcher.new()
    + flowMatcher.withControlPoint('ingress')
  );

local policyDef =
  policy.new()
  + policy.withResources(
    resources.new()
    + resources.withFlowControl(flowControlResources.new()
                                + flowControlResources.withFluxMetersMixin(
                                  { test: fluxMeter.new() + fluxMeter.withFlowSelector(svcSelector) }
                                ))
  )
  + policy.withCircuit(
    circuit.new()
    + circuit.withEvaluationInterval('0.5s')
    + circuit.withComponents([
      component.withQuery(
        query.new()
        + query.withPromql(
          local q = 'sum(increase(flux_meter_sum{decision_type!="DECISION_TYPE_REJECTED", flow_status="OK", flux_meter_name="test"}[5s]))/sum(increase(flux_meter_count{decision_type!="DECISION_TYPE_REJECTED", flow_status="OK", flux_meter_name="test"}[5s]))';
          promQL.new()
          + promQL.withQueryString(q)
          + promQL.withEvaluationInterval('1s')
          + promQL.withOutPorts({ output: port.withSignalName('LATENCY') }),
        ),
      ),
      component.withEma(
        ema.withParameters(
          emaParameters.new()
          + emaParameters.withEmaWindow('1500s')
          + emaParameters.withWarmupWindow('10s')
        )
        + ema.withInPortsMixin(ema.inPorts.withInput(port.withSignalName('LATENCY')))
        + ema.withOutPortsMixin(ema.outPorts.withOutput(port.withSignalName('LATENCY_EMA')))
      ),
      component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY_EMA'),
                                                        port.withConstantSignal(1.1),
                                                        output=port.withSignalName('LATENCY_SETPOINT'))),
      component.withDecider(
        decider.new()
        + decider.withOperator('gt')
        + decider.withInPortsMixin(
          decider.inPorts.withLhs(port.withSignalName('LATENCY'))
          + decider.inPorts.withRhs(port.withSignalName('LATENCY_SETPOINT'))
        )
        + decider.withOutPortsMixin(decider.outPorts.withOutput(port.withSignalName('IS_OVERLOAD_SWITCH')))
      ),
      component.withAlerter(
        alerter.new()
        + alerter.withInPorts({ signal: port.withSignalName('IS_OVERLOAD_SWITCH') })
        + alerter.withParameters(
          alerterParameters.new()
          + alerterParameters.withAlertName('overload')
          + alerterParameters.withSeverity('crit')
        )
      ),
    ]),
  );

local policyResource = {
  kind: 'Policy',
  apiVersion: 'fluxninja.com/v1alpha1',
  metadata: {
    name: 'signal-processing',
    labels: {
      'fluxninja.com/validate': 'true',
    },
  },
  spec: policyDef,
};

policyResource
