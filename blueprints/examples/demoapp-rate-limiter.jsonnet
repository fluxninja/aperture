local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

local Override = aperture.v1.RateLimiterOverride;
local LazySync = aperture.v1.RateLimiterLazySync;

{
  policy+: {
    policyName: 'service1-rate-limiter',
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
}
