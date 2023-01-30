local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local latencyAIMDPolicy = aperture.policies.LatencyAIMDConcurrencyLimiting.policy;

local flowSelector = aperture.spec.v1.FlowSelector;
local fluxMeter = aperture.spec.v1.FluxMeter;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local controlPoint = aperture.spec.v1.ControlPoint;

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

local policyResource = latencyAIMDPolicy({
  policy_name: 'service1-demo-app',
  flux_meter: fluxMeter.new() + fluxMeter.withFlowSelector(svcSelector),
  concurrency_controller+: {
    flow_selector: svcSelector,
    dynamic_config: {
      dry_run: false,
    },
  },
}).policyResource;

policyResource
