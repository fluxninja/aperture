local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

local Workload = aperture.v1.SchedulerWorkload;
local LabelMatcher = aperture.v1.LabelMatcher;
local WorkloadWithLabelMatcher = aperture.v1.SchedulerWorkloadAndLabelMatcher;
local classifier = aperture.v1.Classifier;
local extractor = aperture.v1.Extractor;
local rule = aperture.v1.Rule;
local selector = aperture.v1.Selector;
local controlPoint = aperture.v1.ControlPoint;

local svcSelector = {
  service: 'service1-demo-app.demoapp.svc.cluster.local',
  controlPoint: {
    traffic: 'ingress',
  },
};

local svcSelector = selector.new()
                    + selector.withAgentGroup('default')
                    + selector.withService('service1-demo-app.demoapp.svc.cluster.local')
                    + selector.withControlPoint(controlPoint.new()
                                                + controlPoint.withTraffic('ingress'));


{
  common+: {
    policyName: 'service1-latency-gradient',
  },
  policy+: {
    policyName: $.common.policyName,
    fluxMeterSelector: svcSelector,
    concurrencyLimiterSelector: svcSelector,
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
      defaultWorkload: {
        priority: 20,
      },
      workloads: [
        WorkloadWithLabelMatcher.new()
        + WorkloadWithLabelMatcher.withWorkload(Workload.withPriority(50))
        // match the label extracted by classifier
        + WorkloadWithLabelMatcher.withLabelMatcher(LabelMatcher.withMatchLabels({ user_type: 'guest' })),
        WorkloadWithLabelMatcher.new()
        + WorkloadWithLabelMatcher.withWorkload(Workload.withPriority(200))
        // match the http header directly
        + WorkloadWithLabelMatcher.withLabelMatcher(LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
      ],
    },
  },
  dashboard+: {
    policyName: $.common.policyName,
  },
}
