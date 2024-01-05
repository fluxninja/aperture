local k = import 'k.libsonnet';

local container = k.core.v1.container;
local containerPort = k.core.v1.containerPort;
local namespace = k.core.v1.namespace;
local deployment = k.apps.v1.deployment;
local service = k.core.v1.service;
local servicePort = k.core.v1.servicePort;
local serviceAccount = k.core.v1.serviceAccount;

local defaults = {
  environment:: {
    namespace: error 'environment.namespace is required',
    name: error 'environment.name is required',
  },
  values:: {
    image: {
      repository: error 'values.image.repository is required',
      tag: 'latest',
      pullPolicy: 'IfNotPresent',
    },
    service: {
      type: 'ClusterIP',
      port: 80,
    },
    labels: {
      'app.kubernetes.io/name': 'example',
      'app.kubernetes.io/instance': $.environment.name,
    },
    app_port: 8080,
    agent: {
      address: 'aperture-agent.aperture-agent.svc.cluster.local:8080',
    },
    extraEnv: {},
  },
};

function(values={}, environment={}) {
  local _merged = defaults { environment+: environment, values+: values },
  local _environment = _merged.environment,
  local _values = _merged.values,
  deployment:
    deployment.new(name=_environment.name, containers=[
      container.new(_environment.name, image='%(repository)s:%(tag)s' % _values.image)
      + container.withImagePullPolicy(_values.image.pullPolicy)
      + container.withPorts([
        containerPort.newNamed(_values.app_port, 'srvhttp'),
      ])
      + container.withEnvMap({
        APERTURE_APP_PORT: std.toString(_values.app_port),
        APERTURE_AGENT_ADDRESS: _values.agent.address,
        APERTURE_AGENT_INSECURE: 'true',
      } + _values.extraEnv),
    ])
    + deployment.metadata.withLabels(_values.labels)
    + deployment.metadata.withNamespace(_environment.namespace)
    + deployment.spec.selector.withMatchLabels(_values.labels)
    + deployment.spec.template.metadata.withLabels(_values.labels)
    + deployment.spec.template.spec.withServiceAccountName(_environment.name),
  service:
    service.new($.deployment.metadata.name, selector=_values.labels, ports=[
      servicePort.newNamed(name='http', port=_values.service.port, targetPort='srvhttp'),
    ])
    + service.spec.withSelector(_values.labels)
    + service.metadata.withNamespace(_environment.namespace),
  serviceAccount:
    serviceAccount.new($.deployment.metadata.name)
    + serviceAccount.metadata.withNamespace(_environment.namespace)
    + serviceAccount.metadata.withLabels(_values.labels),
}
