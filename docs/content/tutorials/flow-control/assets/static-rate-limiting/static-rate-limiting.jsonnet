local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local policy = aperture.spec.v1.Policy;
local component = aperture.spec.v1.Component;
local rateLimiter = aperture.spec.v1.RateLimiter;
local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local circuit = aperture.spec.v1.Circuit;
local port = aperture.spec.v1.Port;

local rateLimitPort = port.new() + port.withSignalName('RATE_LIMIT');

local svcSelector =
  flowSelector.new()
  + flowSelector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + flowSelector.withFlowMatcher(
    flowMatcher.new()
    + flowMatcher.withControlPoint({ traffic: 'ingress' })
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
        + rateLimiter.withFlowSelector(svcSelector)
        + rateLimiter.withLimitResetInterval('60s')
        + rateLimiter.withLabelKey('http.request.header.user_id')
      ),
    ]),
  );

local policyResource = {
  kind: 'Policy',
  apiVersion: 'fluxninja.com/v1alpha1',
  metadata: {
    name: 'static-rate-limiting',
    labels: {
      'fluxninja.com/validate': 'true',
    },
  },
  spec: policyDef,
};

policyResource
