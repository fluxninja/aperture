local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local latencyGradientPolicy = aperture.blueprints.LatencyGradient.policy;

local selector = aperture.spec.v1.Selector;
local fluxMeter = aperture.spec.v1.FluxMeter;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowSelector = aperture.spec.v1.FlowSelector;
local controlPoint = aperture.spec.v1.ControlPoint;

local svcSelector =
  selector.new()
  + selector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + selector.withFlowSelector(
    flowSelector.new()
    + flowSelector.withControlPoint(controlPoint.new()
                                    + controlPoint.withTraffic('ingress'))
  );

local policyResource = latencyGradientPolicy({
  policyName: 'service1-demo-app',
  fluxMeter: fluxMeter.new() + fluxMeter.withSelector(svcSelector),
  concurrencyLimiterSelector: svcSelector,
  dynamicConfig: {
    dryRun: false,
  },
}).policyResource;

policyResource
