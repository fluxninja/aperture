local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local policy = aperture.spec.v1.Policy;
local component = aperture.spec.v1.Component;
local selector = aperture.spec.v1.Selector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowSelector = aperture.spec.v1.FlowSelector;
local circuit = aperture.spec.v1.Circuit;
local port = aperture.spec.v1.Port;
local resources = aperture.spec.v1.Resources;
local fluxMeter = aperture.spec.v1.FluxMeter;
local promQL = aperture.spec.v1.PromQL;
local ema = aperture.spec.v1.EMA;
local combinator = aperture.spec.v1.ArithmeticCombinator;
local decider = aperture.spec.v1.Decider;
local sink = aperture.spec.v1.Sink;

local svcSelector =
  selector.new()
  + selector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + selector.withFlowSelector(
    flowSelector.new()
    + flowSelector.withControlPoint({ traffic: 'ingress' })
  );

local policyDef =
  policy.new()
  + policy.withResources(resources.new()
                         + resources.withFluxMetersMixin({ test: fluxMeter.new() + fluxMeter.withSelector(svcSelector) }))
  + policy.withCircuit(
    circuit.new()
    + circuit.withEvaluationInterval('0.5s')
    + circuit.withComponents([
      component.withPromql(
        local q = 'sum(increase(flux_meter_sum{decision_type!="DECISION_TYPE_REJECTED", response_status="OK", flux_meter_name="test"}[5s]))/sum(increase(flux_meter_count{decision_type!="DECISION_TYPE_REJECTED", response_status="OK", flux_meter_name="test"}[5s]))';
        promQL.new()
        + promQL.withQueryString(q)
        + promQL.withEvaluationInterval('1s')
        + promQL.withOutPorts({ output: port.withSignalName('LATENCY') }),
      ),
      component.withEma(
        ema.withEmaWindow('1500s')
        + ema.withWarmUpWindow('10s')
        + ema.withInPortsMixin(ema.inPorts.withInput(port.withSignalName('LATENCY')))
        + ema.withOutPortsMixin(ema.outPorts.withOutput(port.withSignalName('LATENCY_EMA')))
      ),
      component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY_EMA'),
                                                        port.withConstantValue('1.1'),
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
      component.withSink(
        sink.new()
        + sink.withInPorts({ inputs: [port.withSignalName('IS_OVERLOAD_SWITCH')] })
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
