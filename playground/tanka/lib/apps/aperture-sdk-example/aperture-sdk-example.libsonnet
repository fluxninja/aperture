local k = import 'k.libsonnet';

local container = k.core.v1.container;
local containerPort = k.core.v1.containerPort;
local deployment = k.apps.v1.deployment;
local service = k.core.v1.service;
local servicePort = k.core.v1.servicePort;

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
      'app.kubernetes.io/name': $.environment.name,
    },
    app_port: 8080,
    agent: {
      host: 'aperture-agent.aperture-agent.svc.cluster.local',
      port: 8080,
    },
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
        containerPort.newNamed(_values.app_port, 'http'),
      ])
      + container.withEnvMap({
        FN_APP_PORT: std.toString(_values.app_port),
        APERTURE_AGENT_HOST: _values.agent.host,
        APERTURE_AGENT_PORT: std.toString(_values.agent.port),
      }),
    ])
    + deployment.metadata.withNamespace(_environment.namespace)
    + deployment.spec.selector.withMatchLabels(_values.labels)
    + deployment.spec.template.metadata.withLabels(_values.labels),
  service:
    service.new($.deployment.metadata.name, selector=_values.labels, ports=[
      local portName = 'http';
      servicePort.newNamed(name=portName, port=_values.service.port, targetPort=portName),
    ])
    + service.metadata.withNamespace(_environment.namespace)
    + service.metadata.withLabels(_values.labels),
}
