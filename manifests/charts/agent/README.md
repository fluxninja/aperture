# Aperture Agent

Aperture Agent

## Introduction

Aperture Agent

## Installation

## Parameters

### Global parameters

| Name                      | Description                                     | Value             |
| ------------------------- | ----------------------------------------------- | ----------------- |
| `global.imageRegistry`    | Global Docker image registry                    | `""`              |
| `global.imagePullSecrets` | Global Docker registry secret names as an array | `[]`              |
| `global.istioNamespace`   | Specifies namespace for Istio resources         | `aperture-system` |


### Common parameters

| Name                | Description                                       | Value           |
| ------------------- | ------------------------------------------------- | --------------- |
| `kubeVersion`       | Override Kubernetes version                       | `""`            |
| `nameOverride`      | String to partially override common.names.name    | `""`            |
| `fullnameOverride`  | String to fully override common.names.fullname    | `""`            |
| `namespaceOverride` | String to fully override common.names.namespace   | `""`            |
| `commonLabels`      | Labels to add to all deployed objects             | `{}`            |
| `commonAnnotations` | Annotations to add to all deployed objects        | `{}`            |
| `clusterDomain`     | Kubernetes cluster domain name                    | `cluster.local` |
| `extraDeploy`       | Array of extra objects to deploy with the release | `[]`            |


### FluxNinja cloud integration parameters

| Name                                                        | Description                                                                                                              | Value    |
| ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ | -------- |
| `fluxninjaPlugin.enabled`                                   | Boolean flag for enabling FluxNinja cloud connection from Agent and Controller                                           | `false`  |
| `fluxninjaPlugin.endpoint`                                  | FluxNinja cloud instance endpoint                                                                                        | `nil`    |
| `fluxninjaPlugin.heartbeatsInterval`                        | specifies how often to send heartbeats to the cloud                                                                      | `30s`    |
| `fluxninjaPlugin.tls.insecure`                              | specifies whether to communicate with FluxNinja cloud over TLS or in plain text                                          | `false`  |
| `fluxninjaPlugin.tls.insecureSkipVerify`                    | specifies whether to verify FluxNinja cloud certificate                                                                  | `false`  |
| `fluxninjaPlugin.tls.caFile`                                | specifies an alternative CA certificates bundle to use to validate FluxNinja cloud certificate                           | `nil`    |
| `fluxninjaPlugin.apiKeySecret.agent.create`                 | Whether to create Kubernetes Secret with provided Agent API Key                                                          | `true`   |
| `fluxninjaPlugin.apiKeySecret.agent.secretKeyRef.name`      | specifies a name of the Secret for Agent API Key to be used. This defaults to {{ .Release.Name }}-agent-apikey           | `nil`    |
| `fluxninjaPlugin.apiKeySecret.agent.secretKeyRef.key`       | specifies which key from the Secret for Agent API Key to use                                                             | `apiKey` |
| `fluxninjaPlugin.apiKeySecret.agent.value`                  | API Key to use when creating a new Agent API Key Secret                                                                  | `nil`    |
| `fluxninjaPlugin.apiKeySecret.controller.create`            | Whether to create Kubernetes Secret with provided Controller API Key                                                     | `true`   |
| `fluxninjaPlugin.apiKeySecret.controller.secretKeyRef.name` | specifies a name of the Secret for Controller API Key to be used. This defaults to {{ .Release.Name }}-controller-apikey | `nil`    |
| `fluxninjaPlugin.apiKeySecret.controller.secretKeyRef.key`  | specifies which key from the Secret for Controller API Key to use                                                        | `apiKey` |
| `fluxninjaPlugin.apiKeySecret.controller.value`             | API Key to use when creating a new Controller API Key Secret                                                             | `nil`    |


### Agent Parameters

| Name                                                    | Description                                                                                                                 | Value                       |
| ------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| `agent.image.registry`                                  | agent image registry                                                                                                        | `gcr.io/devel-309501/cf-fn` |
| `agent.image.repository`                                | agent image repository                                                                                                      | `aperture-agent`            |
| `agent.image.tag`                                       | agent image tag (immutable tags are recommended)                                                                            | `latest`                    |
| `agent.image.pullPolicy`                                | agent image pull policy                                                                                                     | `IfNotPresent`              |
| `agent.image.pullSecrets`                               | agent image pull secrets                                                                                                    | `[]`                        |
| `agent.serverPort`                                      | The agent's server port                                                                                                     | `80`                        |
| `agent.distributedCachePort`                            | Port for Agent's distributed cache endpoint                                                                                 | `3320`                      |
| `agent.memberListPort`                                  | Port for Agent's member list endpoint                                                                                       | `3322`                      |
| `agent.log.prettyConsole`                               | Additional log writer: pretty console (stdout) logging (not recommended for prod environments)                              | `false`                     |
| `agent.log.nonBlocking`                                 | Use non-blocking log writer (can lose logs at high throughput)                                                              | `true`                      |
| `agent.log.level`                                       | Log level. Keywords allowed - ["debug", "info", "warn", "fatal", "panic", "trace"]. Defaults to info                        | `info`                      |
| `agent.log.file`                                        | Output file for logs. Keywords allowed - ["stderr", "stderr", "default"]. "default" maps to /var/log/aperture/<service>.log | `stderr`                    |
| `agent.livenessProbe.enabled`                           | Enable livenessProbe on agent containers                                                                                    | `true`                      |
| `agent.livenessProbe.initialDelaySeconds`               | Initial delay seconds for livenessProbe                                                                                     | `15`                        |
| `agent.livenessProbe.periodSeconds`                     | Period seconds for livenessProbe                                                                                            | `15`                        |
| `agent.livenessProbe.timeoutSeconds`                    | Timeout seconds for livenessProbe                                                                                           | `5`                         |
| `agent.livenessProbe.failureThreshold`                  | Failure threshold for livenessProbe                                                                                         | `6`                         |
| `agent.livenessProbe.successThreshold`                  | Success threshold for livenessProbe                                                                                         | `1`                         |
| `agent.readinessProbe.enabled`                          | Enable readinessProbe on agent containers                                                                                   | `true`                      |
| `agent.readinessProbe.initialDelaySeconds`              | Initial delay seconds for readinessProbe                                                                                    | `15`                        |
| `agent.readinessProbe.periodSeconds`                    | Period seconds for readinessProbe                                                                                           | `15`                        |
| `agent.readinessProbe.timeoutSeconds`                   | Timeout seconds for readinessProbe                                                                                          | `5`                         |
| `agent.readinessProbe.failureThreshold`                 | Failure threshold for readinessProbe                                                                                        | `6`                         |
| `agent.readinessProbe.successThreshold`                 | Success threshold for readinessProbe                                                                                        | `1`                         |
| `agent.customLivenessProbe`                             | Custom livenessProbe that overrides the default one                                                                         | `{}`                        |
| `agent.customReadinessProbe`                            | Custom readinessProbe that overrides the default one                                                                        | `{}`                        |
| `agent.resources.limits`                                | The resources limits for the agent containers                                                                               | `{}`                        |
| `agent.resources.requests`                              | The requested resources for the agent containers                                                                            | `{}`                        |
| `agent.podSecurityContext.enabled`                      | Enabled agent pods' Security Context                                                                                        | `false`                     |
| `agent.podSecurityContext.fsGroup`                      | Set agent pod's Security Context fsGroup                                                                                    | `1001`                      |
| `agent.containerSecurityContext.enabled`                | Enabled agent containers' Security Context                                                                                  | `false`                     |
| `agent.containerSecurityContext.runAsUser`              | Set agent containers' Security Context runAsUser                                                                            | `1001`                      |
| `agent.containerSecurityContext.runAsNonRoot`           | Set agent containers' Security Context runAsNonRoot                                                                         | `false`                     |
| `agent.containerSecurityContext.readOnlyRootFilesystem` | Set agent containers' Security Context runAsNonRoot                                                                         | `false`                     |
| `agent.command`                                         | Override default container command (useful when using custom images)                                                        | `[]`                        |
| `agent.args`                                            | Override default container args (useful when using custom images)                                                           | `[]`                        |
| `agent.podLabels`                                       | Extra labels for agent pods                                                                                                 | `{}`                        |
| `agent.podAnnotations`                                  | Annotations for agent pods                                                                                                  | `{}`                        |
| `agent.affinity`                                        | Affinity for agent pods assignment                                                                                          | `{}`                        |
| `agent.nodeSelector`                                    | Node labels for agent pods assignment                                                                                       | `{}`                        |
| `agent.tolerations`                                     | Tolerations for agent pods assignment                                                                                       | `[]`                        |
| `agent.terminationGracePeriodSeconds`                   | configures how long kubelet gives agent chart to terminate cleanly                                                          | `""`                        |
| `agent.lifecycleHooks`                                  | for the agent container(s) to automate configuration before or after startup                                                | `{}`                        |
| `agent.extraEnvVars`                                    | Array with extra environment variables to add to agent nodes                                                                | `[]`                        |
| `agent.extraEnvVarsCM`                                  | Name of existing ConfigMap containing extra env vars for agent nodes                                                        | `""`                        |
| `agent.extraEnvVarsSecret`                              | Name of existing Secret containing extra env vars for agent nodes                                                           | `""`                        |
| `agent.extraVolumes`                                    | Optionally specify extra list of additional volumes for the agent pod(s)                                                    | `[]`                        |
| `agent.extraVolumeMounts`                               | Optionally specify extra list of additional volumeMounts for the agent container(s)                                         | `[]`                        |
| `agent.sidecars`                                        | Add additional sidecar containers to the agent pod(s)                                                                       | `{}`                        |
| `agent.initContainers`                                  | Add additional init containers to the agent pod(s)                                                                          | `{}`                        |
| `agent.serviceAccount.create`                           | Specifies whether a ServiceAccount should be created                                                                        | `true`                      |
| `agent.serviceAccount.annotations`                      | Additional Service Account annotations (evaluated as a template)                                                            | `{}`                        |
| `agent.serviceAccount.automountServiceAccountToken`     | Automount service account token for the server service account                                                              | `true`                      |


### Controller Parameters

| Name                                                              | Description                                                                                                                 | Value                       |
| ----------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| `agentController.image.registry`                                  | agent controller image registry                                                                                             | `gcr.io/devel-309501/cf-fn` |
| `agentController.image.repository`                                | agent controller image repository                                                                                           | `aperture-controller`       |
| `agentController.image.tag`                                       | agent controller image tag (immutable tags are recommended)                                                                 | `latest`                    |
| `agentController.image.pullPolicy`                                | agent controller image pull policy                                                                                          | `IfNotPresent`              |
| `agentController.image.pullSecrets`                               | agent controller image pull secrets                                                                                         | `[]`                        |
| `agentController.serverPort`                                      | The controller's server port                                                                                                | `80`                        |
| `agentController.log.prettyConsole`                               | Additional log writer: pretty console (stdout) logging (not recommended for prod environments)                              | `false`                     |
| `agentController.log.nonBlocking`                                 | Use non-blocking log writer (can lose logs at high throughput)                                                              | `true`                      |
| `agentController.log.level`                                       | Log level. Keywords allowed - ["debug", "info", "warn", "fatal", "panic", "trace"]. Defaults to info                        | `info`                      |
| `agentController.log.file`                                        | Output file for logs. Keywords allowed - ["stderr", "stderr", "default"]. "default" maps to /var/log/aperture/<service>.log | `stderr`                    |
| `agentController.livenessProbe.enabled`                           | Enable livenessProbe on agent controller containers                                                                         | `true`                      |
| `agentController.livenessProbe.initialDelaySeconds`               | Initial delay seconds for livenessProbe                                                                                     | `15`                        |
| `agentController.livenessProbe.periodSeconds`                     | Period seconds for livenessProbe                                                                                            | `15`                        |
| `agentController.livenessProbe.timeoutSeconds`                    | Timeout seconds for livenessProbe                                                                                           | `5`                         |
| `agentController.livenessProbe.failureThreshold`                  | Failure threshold for livenessProbe                                                                                         | `6`                         |
| `agentController.livenessProbe.successThreshold`                  | Success threshold for livenessProbe                                                                                         | `1`                         |
| `agentController.readinessProbe.enabled`                          | Enable readinessProbe on agent controller containers                                                                        | `true`                      |
| `agentController.readinessProbe.initialDelaySeconds`              | Initial delay seconds for readinessProbe                                                                                    | `15`                        |
| `agentController.readinessProbe.periodSeconds`                    | Period seconds for readinessProbe                                                                                           | `15`                        |
| `agentController.readinessProbe.timeoutSeconds`                   | Timeout seconds for readinessProbe                                                                                          | `5`                         |
| `agentController.readinessProbe.failureThreshold`                 | Failure threshold for readinessProbe                                                                                        | `6`                         |
| `agentController.readinessProbe.successThreshold`                 | Success threshold for readinessProbe                                                                                        | `1`                         |
| `agentController.customLivenessProbe`                             | Custom livenessProbe that overrides the default one                                                                         | `{}`                        |
| `agentController.customReadinessProbe`                            | Custom readinessProbe that overrides the default one                                                                        | `{}`                        |
| `agentController.resources.limits`                                | The resources limits for the agentController containers                                                                     | `{}`                        |
| `agentController.resources.requests`                              | The requested resources for the agentController containers                                                                  | `{}`                        |
| `agentController.podSecurityContext.enabled`                      | Enables agentController pods' Security Context                                                                              | `false`                     |
| `agentController.podSecurityContext.fsGroup`                      | Set agentController pod's Security Context fsGroup                                                                          | `1001`                      |
| `agentController.containerSecurityContext.enabled`                | Enables agentController containers' Security Context                                                                        | `false`                     |
| `agentController.containerSecurityContext.runAsUser`              | Set agentController containers' Security Context runAsUser                                                                  | `1001`                      |
| `agentController.containerSecurityContext.runAsNonRoot`           | Set agentController containers' Security Context runAsNonRoot                                                               | `false`                     |
| `agentController.containerSecurityContext.readOnlyRootFilesystem` | Set agent controller containers' Security Context runAsNonRoot                                                              | `false`                     |
| `agentController.command`                                         | Override default container command (useful when using custom images)                                                        | `[]`                        |
| `agentController.args`                                            | Override default container args (useful when using custom images)                                                           | `[]`                        |
| `agentController.hostAliases`                                     | agent controller pods host aliases                                                                                          | `[]`                        |
| `agentController.podLabels`                                       | Extra labels for agent controller pods                                                                                      | `{}`                        |
| `agentController.podAnnotations`                                  | Annotations for agent controller pods                                                                                       | `{}`                        |
| `agentController.affinity`                                        | Affinity for agent controller pods assignment                                                                               | `{}`                        |
| `agentController.nodeSelector`                                    | Node labels for agent controller pods assignment                                                                            | `{}`                        |
| `agentController.tolerations`                                     | Tolerations for agent controller pods assignment                                                                            | `[]`                        |
| `agentController.terminationGracePeriodSeconds`                   | Seconds Redmine pod needs to terminate gracefully                                                                           | `""`                        |
| `agentController.lifecycleHooks`                                  | for the agent controller container(s) to automate configuration before or after startup                                     | `{}`                        |
| `agentController.extraEnvVars`                                    | Array with extra environment variables to add to agent controller nodes                                                     | `[]`                        |
| `agentController.extraEnvVarsCM`                                  | Name of existing ConfigMap containing extra env vars for agentController nodes                                              | `""`                        |
| `agentController.extraEnvVarsSecret`                              | Name of existing Secret containing extra env vars for agentController nodes                                                 | `""`                        |
| `agentController.extraVolumes`                                    | Optionally specify extra list of additional volumes for the agentController pod(s)                                          | `[]`                        |
| `agentController.extraVolumeMounts`                               | Optionally specify extra list of additional volumeMounts for the agentController container(s)                               | `[]`                        |
| `agentController.sidecars`                                        | Add additional sidecar containers to the agentController pod(s)                                                             | `{}`                        |
| `agentController.initContainers`                                  | Add additional init containers to the agent controller pod(s)                                                               | `{}`                        |
| `agentController.serviceAccount.create`                           | Specifies whether a ServiceAccount should be created                                                                        | `true`                      |
| `agentController.serviceAccount.annotations`                      | Additional Service Account annotations (evaluated as a template)                                                            | `{}`                        |
| `agentController.serviceAccount.automountServiceAccountToken`     | Automount service account token for the server service account                                                              | `true`                      |


### Traffic Exposure Parameters

| Name                                  | Description                                                | Value |
| ------------------------------------- | ---------------------------------------------------------- | ----- |
| `service.agent.annotations`           | Additional custom annotations for agent service            | `{}`  |
| `service.agentController.annotations` | Additional custom annotations for agent controller service | `{}`  |


### Other Parameters




### etcd

| Name                    | Description                                                                                                              | Value    |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------ | -------- |
| `etcd.enabled`          | Whether to deploy a small etcd cluster as part of this chart                                                             | `true`   |
| `etcd.endpoints`        | List of Etcd server endpoints. Example, ["https://etcd:2379"]. This must not be empty when etcd.enabled is set to false. | `[]`     |
| `etcd.leaseTtl`         | Lease time-to-live.                                                                                                      | `60s`    |
| `etcd.auth.rbac.create` | specifies whether to create the RBAC resources for Etcd                                                                  | `false`  |
| `etcd.auth.token.type`  | specifies the type of token to use                                                                                       | `simple` |


### prometheus

| Name                 | Description                                                                                                     | Value  |
| -------------------- | --------------------------------------------------------------------------------------------------------------- | ------ |
| `prometheus.enabled` | specifies whether to deploy embedded prometheus                                                                 | `true` |
| `prometheus.address` | specifies the address of the Prometheus server. This must not be empty when prometheus.enabled is set to false. | `""`   |


### istio

| Name                                 | Description                                         | Value                   |
| ------------------------------------ | --------------------------------------------------- | ----------------------- |
| `istio.enabled`                      | specifies whether to deploy embedded istio          | `true`                  |
| `istio.global.tag`                   | specifies the image tag for all the istio resources | `1.11.8`                |
| `istio.envoyFilter.install`          | specifies whether to deploy EnvoyFilter             | `true`                  |
| `istio.envoyFilter.name`             | specifies name for the EnvoyFilter                  | `aperture-envoy-filter` |
| `istio.envoyFilter.authzGrpcTimeout` | specifies timeout for the AuthZ gRPC connection     | `0.01s`                 |
| `istio.envoyFilter.maxRequestBytes`  | specifies allowed maximum request bytes             | `8192`                  |

