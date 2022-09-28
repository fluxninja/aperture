local aperture = import '../../../lib/1.0/main.libsonnet';
local blueprint = import '../main.libsonnet';

local WorkloadParameters = aperture.spec.v1.SchedulerWorkloadParameters;
local LabelMatcher = aperture.spec.v1.LabelMatcher;
local Workload = aperture.spec.v1.SchedulerWorkload;
local classifier = aperture.spec.v1.Classifier;
local fluxMeter = aperture.spec.v1.FluxMeter;
local extractor = aperture.spec.v1.Extractor;
local rule = aperture.spec.v1.Rule;
local selector = aperture.spec.v1.Selector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowSelector = aperture.spec.v1.FlowSelector;
local controlPoint = aperture.spec.v1.ControlPoint;
local staticBuckets = aperture.spec.v1.FluxMeterStaticBuckets;

local svcSelector = {
  sSelector: {
    service: 'service1-demo-app.demoapp.svc.cluster.local',
  },
  fSelector: {
    controlPoint: {
      traffic: 'ingress',
    },
  },
};

local svcSelector = selector.new()
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


local config = {
  common+: {
    policyName: 'example',
  },
  policy+: {
    policyName: $.common.policyName,
    fluxMeterSelector: svcSelector,
    fluxMeters: {
      [$.common.policyName]: fluxMeter.new()
                             + fluxMeter.withSelector(svcSelector)
                             + fluxMeter.withAttributeKey('workload_duration_ms')
                             + fluxMeter.withStaticBuckets(
                               staticBuckets.new()
                               + staticBuckets.withBuckets([5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2500.0, 5000.0, 10000.0])
                             ),
    },
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
  dashboard+: {
    policyName: $.common.policyName,
  },
};

blueprint { _config:: config }
