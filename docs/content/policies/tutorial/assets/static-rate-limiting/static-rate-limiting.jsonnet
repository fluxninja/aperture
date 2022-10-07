local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local policy = aperture.spec.v1.Policy;
local component = aperture.spec.v1.Component;
local rateLimiter = aperture.spec.v1.RateLimiter;
local selector = aperture.spec.v1.Selector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowSelector = aperture.spec.v1.FlowSelector;
local circuit = aperture.spec.v1.Circuit;
local port = aperture.spec.v1.Port;

local rateLimitPort = port.new() + port.withSignalName('RATE_LIMIT');

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
  + policy.withCircuit(
    circuit.new()
    + circuit.withEvaluationInterval('300s')
    + circuit.withComponents([
      component.withRateLimiter(
        rateLimiter.new()
        + rateLimiter.withInPorts({ limit: port.withConstantValue(120) })
        + rateLimiter.withSelector(svcSelector)
        + rateLimiter.withLimitResetInterval('60s')
        + rateLimiter.withLabelKey('http.request.header.user_id')
      ),
    ]),
  );

local policyResource = {
  kind: 'Policy',
  apiVersion: 'fluxninja.com/v1alpha1',
  metadata: {
    name: 'service1-demo-app',
    labels: {
      'fluxninja.com/validate': 'true',
    },
  },
  spec: policyDef,
};

[
  policyResource,
]
