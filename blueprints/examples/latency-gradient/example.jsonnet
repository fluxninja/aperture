local aperture = import '../../lib/1.0/main.libsonnet';
local bundle = aperture.blueprints.LatencyGradient.bundle;

local WorkloadParameters = aperture.spec.v1.SchedulerWorkloadParameters;
local LabelMatcher = aperture.spec.v1.LabelMatcher;
local Workload = aperture.spec.v1.SchedulerWorkload;
local classifier = aperture.spec.v1.Classifier;
local fluxMeter = aperture.spec.v1.FluxMeter;
local extractor = aperture.spec.v1.Extractor;
local rule = aperture.spec.v1.Rule;
local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local controlPoint = aperture.spec.v1.ControlPoint;
local staticBuckets = aperture.spec.v1.FluxMeterStaticBuckets;

local svcSelector = flowSelector.new()
                    + flowSelector.withServiceSelector(
                      serviceSelector.new()
                      + serviceSelector.withAgentGroup('default')
                      + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
                    )
                    + flowSelector.withFlowMatcher(
                      flowMatcher.new()
                      + flowMatcher.withControlPoint('ingress')
                    );


local config = {
  common+: {
    policyName: 'example',
  },
  policy+: {
    fluxMeter: fluxMeter.new() + fluxMeter.withFlowSelector(svcSelector),
    concurrencyLimiterFlowSelector: svcSelector,
    classifiers: [
      classifier.new()
      + classifier.withFlowSelector(svcSelector)
      + classifier.withRules({
        user_type: rule.new()
                   + rule.withExtractor(extractor.new()
                                        + extractor.withFrom('request.http.headers.user-type')),
      }),
    ],
    concurrencyLimiter+: {
      defaultWorkloadParameters: {
        priority: 20,
      },
      workloads: [
        Workload.new()
        + Workload.withWorkloadParameters(WorkloadParameters.withPriority(50))
        // match the label extracted by classifier
        + Workload.withLabelMatcher(LabelMatcher.withMatchLabels({ user_type: 'guest' })),
        Workload.new()
        + Workload.withWorkloadParameters(WorkloadParameters.withPriority(200))
        // match the http header directly
        + Workload.withLabelMatcher(LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
      ],
    },
  },
};

bundle { _config+:: config }
