local aperture = import '../../../libsonnet/1.0/main.libsonnet';
local blueprint = import '../main.libsonnet';

local Override = aperture.v1.RateLimiterOverride;
local LazySync = aperture.v1.RateLimiterLazySync;

local config = {
  policy+: {
    policyName: 'example',
    selector+: {
      serviceSelector+: {
        service: 'service1-demo-app.demoapp.svc.cluster.local',
      },
    },
    rateLimit: '50.0',
    labelKey: 'http.request.header.user_type',
    limitResetInterval: '1s',
    overrides: [
      Override.new()
      + Override.withLabelValue('gold')
      + Override.withLimitScaleFactor(1),
    ],
  },
};

blueprint { _config:: config }
