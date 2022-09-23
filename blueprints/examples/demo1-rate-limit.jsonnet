{
  _config+:: {
    common: {
      policyName: 'demo1-rate-limit',
    },
    policy+: {
      rateLimit: '250.0',
      labelKey: 'request_header_user-type',
      selector+: {
        serviceSelector+: {
          service: 'demo1-demo-app.demoapp.svc.cluster.local',
        },
      },
    },
  },
}
