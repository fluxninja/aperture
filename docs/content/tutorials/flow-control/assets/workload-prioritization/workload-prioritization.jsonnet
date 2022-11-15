local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local latencyGradientPolicy = aperture.blueprints.LatencyGradient.policy;

local selector = aperture.spec.v1.Selector;
local fluxMeter = aperture.spec.v1.FluxMeter;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowSelector = aperture.spec.v1.FlowSelector;
local controlPoint = aperture.spec.v1.ControlPoint;
local classifier = aperture.spec.v1.Classifier;
local extractor = aperture.spec.v1.Extractor;
local rule = aperture.spec.v1.Rule;
local workloadParameters = aperture.spec.v1.SchedulerWorkloadParameters;
local labelMatcher = aperture.spec.v1.LabelMatcher;
local workload = aperture.spec.v1.SchedulerWorkload;


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
  // highlight-start
  classifiers: [
    classifier.new()
    + classifier.withSelector(svcSelector)
    + classifier.withRules({
      user_type: rule.new()
                 + rule.withExtractor(extractor.new()
                                      + extractor.withFrom('request.http.headers.user-type')),
    }),
  ],
  concurrencyLimiter+: {
    timeoutFactor: 0.5,
    defaultWorkloadParameters: {
      priority: 20,
    },
    workloads: [
      workload.new()
      + workload.withWorkloadParameters(workloadParameters.withPriority(50))
      // match the label extracted by classifier
      + workload.withLabelMatcher(labelMatcher.withMatchLabels({ user_type: 'guest' })),
      workload.new()
      + workload.withWorkloadParameters(workloadParameters.withPriority(200))
      // match the http header directly
      + workload.withLabelMatcher(labelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
    ],
  },
  // highlight-end
}).policyResource;

policyResource
