# Aperture Operator

Aperture Operator

## Introduction

Aperture Operator

## Installation

## Parameters

### Global Parameters

| Name                      | Description                                     | Value     |
| ------------------------- | ----------------------------------------------- | --------- |
| `global.imageRegistry`    | Global Docker image registry                    | `""`      |
| `global.imagePullSecrets` | Global Docker registry secret names as an array | `[]`      |
| `global.istioNamespace`   | Specifies namespace for Istio resources         | `default` |


### Common Parameters

| Name                | Description                                                          | Value           |
| ------------------- | -------------------------------------------------------------------- | --------------- |
| `kubeVersion`       | Force target Kubernetes version (using Helm capabilities if not set) | `""`            |
| `nameOverride`      | String to partially override common.names.name                       | `""`            |
| `fullnameOverride`  | String to fully override common.names.fullname                       | `""`            |
| `namespaceOverride` | String to fully override common.names.namespace                      | `""`            |
| `extraDeploy`       | Array of extra objects to deploy with the release                    | `[]`            |
| `commonLabels`      | Labels to add to all deployed objects                                | `{}`            |
| `commonAnnotations` | Annotations to add to all deployed objects                           | `{}`            |
| `clusterDomain`     | Kubernetes cluster domain name                                       | `cluster.local` |


### Operator Parameters

| Name                                                         | Description                                                                                                            | Value                       |
| ------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| `operator.image.registry`                                    | Operator image registry                                                                                                | `docker.io/aperturecontrol` |
| `operator.image.repository`                                  | Operator image repository                                                                                              | `aperture-operator`         |
| `operator.image.tag`                                         | Operator image tag (immutable tags are recommended)                                                                    | `latest`                    |
| `operator.image.pullPolicy`                                  | Operator image pull policy                                                                                             | `IfNotPresent`              |
| `operator.image.pullSecrets`                                 | Operator image pull secrets                                                                                            | `[]`                        |
| `operator.replicaCount`                                      | Number of replicas for Operator deployment                                                                             | `1`                         |
| `operator.podAnnotations`                                    | Pod annotations                                                                                                        | `{}`                        |
| `operator.podLabels`                                         | Additional pod labels                                                                                                  | `{}`                        |
| `operator.updateStrategy.type`                               | Set up update strategy for Aperture Operator installation.                                                             | `Recreate`                  |
| `operator.priorityClassName`                                 | pods' priorityClassName                                                                                                | `""`                        |
| `operator.topologySpreadConstraints`                         | Topology Spread Constraints for pod assignment                                                                         | `[]`                        |
| `operator.schedulerName`                                     | Alternative scheduler                                                                                                  | `""`                        |
| `operator.hostAliases`                                       | Add deployment host aliases                                                                                            | `[]`                        |
| `operator.nodeSelector`                                      | Node labels for pod assignment                                                                                         | `{}`                        |
| `operator.podAffinityPreset`                                 | Pod affinity preset                                                                                                    | `""`                        |
| `operator.podAntiAffinityPreset`                             | Pod anti-affinity preset                                                                                               | `soft`                      |
| `operator.nodeAffinityPreset.type`                           | Set nodeAffinity preset type                                                                                           | `""`                        |
| `operator.nodeAffinityPreset.key`                            | Set nodeAffinity preset key                                                                                            | `""`                        |
| `operator.nodeAffinityPreset.values`                         | Set nodeAffinity preset values                                                                                         | `[]`                        |
| `operator.affinity`                                          | Affinity for controller pod assignment                                                                                 | `{}`                        |
| `operator.tolerations`                                       | Tolerations for controller pod assignment                                                                              | `[]`                        |
| `operator.podSecurityContext.enabled`                        | Enable pods security context                                                                                           | `true`                      |
| `operator.podSecurityContext.runAsUser`                      | User ID for the pods                                                                                                   | `1001`                      |
| `operator.podSecurityContext.runAsGroup`                     | User ID for the pods                                                                                                   | `1001`                      |
| `operator.podSecurityContext.runAsNonRoot`                   | Aperture Operator must run as nonRoot                                                                                  | `true`                      |
| `operator.podSecurityContext.fsGroup`                        | Group ID for the pods                                                                                                  | `1001`                      |
| `operator.podSecurityContext.supplementalGroups`             | Which group IDs containers add                                                                                         | `[]`                        |
| `operator.containerSecurityContext.enabled`                  | Enable container security context                                                                                      | `true`                      |
| `operator.containerSecurityContext.runAsUser`                | User ID for the operator container                                                                                     | `1001`                      |
| `operator.containerSecurityContext.runAsGroup`               | User ID for the operator container                                                                                     | `1001`                      |
| `operator.containerSecurityContext.runAsNonRoot`             | Force the container to be run as non-root                                                                              | `true`                      |
| `operator.containerSecurityContext.privileged`               | Decide if the container runs privileged.                                                                               | `false`                     |
| `operator.containerSecurityContext.readOnlyRootFilesystem`   | ReadOnlyRootFilesystem fot the operator container                                                                      | `false`                     |
| `operator.containerSecurityContext.allowPrivilegeEscalation` | Allow Privilege Escalation for the operator container                                                                  | `false`                     |
| `operator.terminationGracePeriodSeconds`                     | In seconds, time the given to the pod needs to terminate gracefully                                                    | `10`                        |
| `operator.initContainers`                                    | Add additional init containers to the etcd pods                                                                        | `[]`                        |
| `operator.command`                                           | Default container command (useful when using custom images)                                                            | `[]`                        |
| `operator.args`                                              | Default container args (useful when using custom images)                                                               | `[]`                        |
| `operator.lifecycleHooks`                                    | for the aperture-operator container to automate configuration before or after startup                                  | `{}`                        |
| `operator.extraEnvVars`                                      | Array with extra environment variables to add to RabbitMQ Cluster Operator nodes                                       | `[]`                        |
| `operator.extraEnvVarsCM`                                    | Name of existing ConfigMap containing extra env vars for RabbitMQ Cluster Operator nodes                               | `""`                        |
| `operator.extraEnvVarsSecret`                                | Name of existing Secret containing extra env vars for RabbitMQ Cluster Operator nodes                                  | `""`                        |
| `operator.resources`                                         | Container resource requests and limits                                                                                 | `{}`                        |
| `operator.livenessProbe.enabled`                             | Enable livenessProbe                                                                                                   | `true`                      |
| `operator.livenessProbe.initialDelaySeconds`                 | Initial delay seconds for livenessProbe                                                                                | `10`                        |
| `operator.livenessProbe.periodSeconds`                       | Period seconds for livenessProbe                                                                                       | `10`                        |
| `operator.livenessProbe.timeoutSeconds`                      | Timeout seconds for livenessProbe                                                                                      | `1`                         |
| `operator.livenessProbe.failureThreshold`                    | Failure threshold for livenessProbe                                                                                    | `3`                         |
| `operator.livenessProbe.successThreshold`                    | Success threshold for livenessProbe                                                                                    | `1`                         |
| `operator.readinessProbe.enabled`                            | Enable readinessProbe                                                                                                  | `true`                      |
| `operator.readinessProbe.initialDelaySeconds`                | Initial delay seconds for readinessProbe                                                                               | `10`                        |
| `operator.readinessProbe.periodSeconds`                      | Period seconds for readinessProbe                                                                                      | `10`                        |
| `operator.readinessProbe.timeoutSeconds`                     | Timeout seconds for readinessProbe                                                                                     | `1`                         |
| `operator.readinessProbe.failureThreshold`                   | Failure threshold for readinessProbe                                                                                   | `3`                         |
| `operator.readinessProbe.successThreshold`                   | Success threshold for readinessProbe                                                                                   | `1`                         |
| `operator.startupProbe.enabled`                              | Enable startupProbe                                                                                                    | `true`                      |
| `operator.startupProbe.initialDelaySeconds`                  | Initial delay seconds for startupProbe                                                                                 | `10`                        |
| `operator.startupProbe.periodSeconds`                        | Period seconds for startupProbe                                                                                        | `10`                        |
| `operator.startupProbe.timeoutSeconds`                       | Timeout seconds for startupProbe                                                                                       | `1`                         |
| `operator.startupProbe.failureThreshold`                     | Failure threshold for startupProbe                                                                                     | `3`                         |
| `operator.startupProbe.successThreshold`                     | Success threshold for startupProbe                                                                                     | `1`                         |
| `operator.customLivenessProbe`                               | Override default liveness probe                                                                                        | `{}`                        |
| `operator.customReadinessProbe`                              | Override default readiness probe                                                                                       | `{}`                        |
| `operator.customStartupProbe`                                | Override default startup probe                                                                                         | `{}`                        |
| `operator.extraVolumes`                                      | Optionally specify extra list of additional volumes                                                                    | `[]`                        |
| `operator.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts                                                               | `[]`                        |
| `operator.rbac.create`                                       | Create specifies whether to install and use RBAC rules                                                                 | `true`                      |
| `operator.serviceAccount.create`                             | Specifies whether a service account should be created                                                                  | `true`                      |
| `operator.serviceAccount.name`                               | The name of the service account to use. If not set and create is true, a name is generated using the fullname template | `""`                        |
| `operator.serviceAccount.annotations`                        | Add annotations                                                                                                        | `{}`                        |
| `operator.serviceAccount.automountServiceAccountToken`       | Automount API credentials for a service account.                                                                       | `true`                      |


### Aperture Custom Resource Parameters

| Name                                                                  | Description                                                                                                              | Value    |
| --------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ | -------- |
| `aperture.create`                                                     | Specifies whether a CR for Aperture should be created                                                                    | `true`   |
| `aperture.fluxninjaPlugin.enabled`                                    | Boolean flag for enabling FluxNinja cloud connection from Agent and Controller                                           | `false`  |
| `aperture.fluxninjaPlugin.endpoint`                                   | FluxNinja cloud instance endpoint                                                                                        | `nil`    |
| `aperture.fluxninjaPlugin.heartbeatsInterval`                         | specifies how often to send heartbeats to the cloud. Defaults to '30s'.                                                  | `nil`    |
| `aperture.fluxninjaPlugin.tls.insecure`                               | specifies whether to communicate with FluxNinja cloud over TLS or in plain text. Defaults to false.                      | `nil`    |
| `aperture.fluxninjaPlugin.tls.insecureSkipVerify`                     | specifies whether to verify FluxNinja cloud certificate. Defaults to false.                                              | `nil`    |
| `aperture.fluxninjaPlugin.tls.caFile`                                 | specifies an alternative CA certificates bundle to use to validate FluxNinja cloud certificate                           | `nil`    |
| `aperture.fluxninjaPlugin.apiKeySecret.agent.create`                  | Whether to create Kubernetes Secret with provided Agent API Key.                                                         | `true`   |
| `aperture.fluxninjaPlugin.apiKeySecret.agent.secretKeyRef.name`       | specifies a name of the Secret for Agent API Key to be used. This defaults to {{ .Release.Name }}-agent-apikey           | `nil`    |
| `aperture.fluxninjaPlugin.apiKeySecret.agent.secretKeyRef.key`        | specifies which key from the Secret for Agent API Key to use                                                             | `apiKey` |
| `aperture.fluxninjaPlugin.apiKeySecret.agent.value`                   | API Key to use when creating a new Agent API Key Secret                                                                  | `nil`    |
| `aperture.fluxninjaPlugin.apiKeySecret.controller.create`             | Whether to create Kubernetes Secret with provided Controller API Key                                                     | `true`   |
| `aperture.fluxninjaPlugin.apiKeySecret.controller.secretKeyRef.name`  | specifies a name of the Secret for Controller API Key to be used. This defaults to {{ .Release.Name }}-controller-apikey | `nil`    |
| `aperture.fluxninjaPlugin.apiKeySecret.controller.secretKeyRef.key`   | specifies which key from the Secret for Controller API Key to use                                                        | `apiKey` |
| `aperture.fluxninjaPlugin.apiKeySecret.controller.value`              | API Key to use when creating a new Controller API Key Secret                                                             | `nil`    |
| `aperture.agent.image.registry`                                       | Agent image registry. Defaults to 'docker.io/aperturecontrol'.                                                           | `nil`    |
| `aperture.agent.image.repository`                                     | Agent image repository. Defaults to 'aperture-agent'.                                                                    | `nil`    |
| `aperture.agent.image.tag`                                            | Agent image tag (immutable tags are recommended). Defaults to 'latest'.                                                  | `nil`    |
| `aperture.agent.image.pullPolicy`                                     | Agent image pull policy. Defaults to 'IfNotPresent'.                                                                     | `nil`    |
| `aperture.agent.image.pullSecrets`                                    | Agent image pull secrets                                                                                                 | `[]`     |
| `aperture.agent.serverPort`                                           | The Agent's server port. Defaults to 80.                                                                                 | `nil`    |
| `aperture.agent.log.prettyConsole`                                    | Additional log writer: pretty console (stdout) logging (not recommended for prod environments). Defaults to false.       | `nil`    |
| `aperture.agent.log.nonBlocking`                                      | Use non-blocking log writer (can lose logs at high throughput). Defaults to True.                                        | `nil`    |
| `aperture.agent.log.level`                                            | Log level. Keywords allowed - ["debug", "info", "warn", "fatal", "panic", "trace"]. Defaults to 'info'.                  | `nil`    |
| `aperture.agent.log.file`                                             | Output file for logs. Keywords allowed - ["stderr", "stderr", "default"]. Defaults to 'stderr'.                          | `nil`    |
| `aperture.agent.agentGroup`                                           | Agent Group name. Defaults to 'default' Agent group                                                                      | `nil`    |
| `aperture.agent.serviceAccount.create`                                | Specifies whether a ServiceAccount should be created                                                                     | `true`   |
| `aperture.agent.serviceAccount.annotations`                           | Additional Service Account annotations (evaluated as a template)                                                         | `{}`     |
| `aperture.agent.serviceAccount.automountServiceAccountToken`          | Automount service account token for the server service account. Defaults to true                                         | `nil`    |
| `aperture.agent.livenessProbe.enabled`                                | Enable livenessProbe on Agent containers                                                                                 | `true`   |
| `aperture.agent.livenessProbe.initialDelaySeconds`                    | Initial delay seconds for livenessProbe. Defaults to 15.                                                                 | `nil`    |
| `aperture.agent.livenessProbe.periodSeconds`                          | Period seconds for livenessProbe. Defaults to 15.                                                                        | `nil`    |
| `aperture.agent.livenessProbe.timeoutSeconds`                         | Timeout seconds for livenessProbe. Defaults to 5.                                                                        | `nil`    |
| `aperture.agent.livenessProbe.failureThreshold`                       | Failure threshold for livenessProbe. Defaults to 6.                                                                      | `nil`    |
| `aperture.agent.livenessProbe.successThreshold`                       | Success threshold for livenessProbe. Defaults to 1.                                                                      | `nil`    |
| `aperture.agent.readinessProbe.enabled`                               | Enable readinessProbe on Agent containers                                                                                | `true`   |
| `aperture.agent.readinessProbe.initialDelaySeconds`                   | Initial delay seconds for readinessProbe. Defaults to 15.                                                                | `nil`    |
| `aperture.agent.readinessProbe.periodSeconds`                         | Period seconds for readinessProbe. Defaults to 15.                                                                       | `nil`    |
| `aperture.agent.readinessProbe.timeoutSeconds`                        | Timeout seconds for readinessProbe. Defaults to 5.                                                                       | `nil`    |
| `aperture.agent.readinessProbe.failureThreshold`                      | Failure threshold for readinessProbe. Defaults to 6.                                                                     | `nil`    |
| `aperture.agent.readinessProbe.successThreshold`                      | Success threshold for readinessProbe. Defaults to 1.                                                                     | `nil`    |
| `aperture.agent.customLivenessProbe`                                  | Custom livenessProbe that overrides the default one                                                                      | `{}`     |
| `aperture.agent.customReadinessProbe`                                 | Custom readinessProbe that overrides the default one                                                                     | `{}`     |
| `aperture.agent.resources.limits`                                     | The resources limits for the Agent containers                                                                            | `{}`     |
| `aperture.agent.resources.requests`                                   | The requested resources for the Agent containers                                                                         | `{}`     |
| `aperture.agent.podSecurityContext.enabled`                           | Enabled Agent pods' Security Context                                                                                     | `false`  |
| `aperture.agent.podSecurityContext.fsGroup`                           | Set Agent pod's Security Context fsGroup. Defaults to 1001.                                                              | `nil`    |
| `aperture.agent.containerSecurityContext.enabled`                     | Enabled Agent containers' Security Context. Defaults to false.                                                           | `false`  |
| `aperture.agent.containerSecurityContext.runAsUser`                   | Set Agent containers' Security Context runAsUser. Defaults to 1001.                                                      | `nil`    |
| `aperture.agent.containerSecurityContext.runAsNonRoot`                | Set Agent containers' Security Context runAsNonRoot. Defaults to false.                                                  | `nil`    |
| `aperture.agent.containerSecurityContext.readOnlyRootFilesystem`      | Set Agent containers' Security Context runAsNonRoot. Defaults to false.                                                  | `nil`    |
| `aperture.agent.command`                                              | Override default container command (useful when using custom images)                                                     | `[]`     |
| `aperture.agent.args`                                                 | Override default container args (useful when using custom images)                                                        | `[]`     |
| `aperture.agent.podLabels`                                            | Extra labels for Agent pods                                                                                              | `{}`     |
| `aperture.agent.podAnnotations`                                       | Annotations for Agent pods                                                                                               | `{}`     |
| `aperture.agent.affinity`                                             | Affinity for Agent pods assignment                                                                                       | `{}`     |
| `aperture.agent.nodeSelector`                                         | Node labels for Agent pods assignment                                                                                    | `{}`     |
| `aperture.agent.tolerations`                                          | Tolerations for Agent pods assignment                                                                                    | `[]`     |
| `aperture.agent.terminationGracePeriodSeconds`                        | configures how long kubelet gives Agent chart to terminate cleanly                                                       | `nil`    |
| `aperture.agent.lifecycleHooks`                                       | for the Agent container(s) to automate configuration before or after startup                                             | `{}`     |
| `aperture.agent.extraEnvVars`                                         | Array with extra environment variables to add to Agent nodes                                                             | `[]`     |
| `aperture.agent.extraEnvVarsCM`                                       | Name of existing ConfigMap containing extra env vars for Agent nodes                                                     | `""`     |
| `aperture.agent.extraEnvVarsSecret`                                   | Name of existing Secret containing extra env vars for Agent nodes                                                        | `""`     |
| `aperture.agent.extraVolumes`                                         | Optionally specify extra list of additional volumes for the Agent pod(s)                                                 | `[]`     |
| `aperture.agent.extraVolumeMounts`                                    | Optionally specify extra list of additional volumeMounts for the Agent container(s)                                      | `[]`     |
| `aperture.agent.sidecars`                                             | Add additional sidecar containers to the Agent pod(s)                                                                    | `[]`     |
| `aperture.agent.initContainers`                                       | Add additional init containers to the Agent pod(s)                                                                       | `[]`     |
| `aperture.controller.image.registry`                                  | Controller image registry. Defaults to 'docker.io/aperturecontrol'.                                                      | `nil`    |
| `aperture.controller.image.repository`                                | Controller image repository. Defaults to 'aperture-controller'.                                                          | `nil`    |
| `aperture.controller.image.tag`                                       | Controller image tag (immutable tags are recommended). Defaults to 'latest'.                                             | `nil`    |
| `aperture.controller.image.pullPolicy`                                | Controller image pull policy. Defaults to 'IfNotPresent'.                                                                | `nil`    |
| `aperture.controller.image.pullSecrets`                               | Controller image pull secrets                                                                                            | `[]`     |
| `aperture.controller.serverPort`                                      | The Controller's server port. Defaults to 80.                                                                            | `nil`    |
| `aperture.controller.log.prettyConsole`                               | Additional log writer: pretty console (stdout) logging (not recommended for prod environments). Defaults to false.       | `nil`    |
| `aperture.controller.log.nonBlocking`                                 | Use non-blocking log writer (can lose logs at high throughput). Defaults to true.                                        | `nil`    |
| `aperture.controller.log.level`                                       | Log level. Keywords allowed - ["debug", "info", "warn", "fatal", "panic", "trace"]. Defaults to 'info'.                  | `nil`    |
| `aperture.controller.log.file`                                        | Output file for logs. Keywords allowed - ["stderr", "stderr", "default"]. Defaults to /stderr.                           | `nil`    |
| `aperture.controller.serviceAccount.create`                           | Specifies whether a ServiceAccount should be created                                                                     | `true`   |
| `aperture.controller.serviceAccount.annotations`                      | Additional Service Account annotations (evaluated as a template)                                                         | `{}`     |
| `aperture.controller.serviceAccount.automountServiceAccountToken`     | Automount service account token for the server service account. Defaults to true.                                        | `nil`    |
| `aperture.controller.livenessProbe.enabled`                           | Enable livenessProbe on Controller containers                                                                            | `true`   |
| `aperture.controller.livenessProbe.initialDelaySeconds`               | Initial delay seconds for livenessProbe. Defaults to 15.                                                                 | `nil`    |
| `aperture.controller.livenessProbe.periodSeconds`                     | Period seconds for livenessProbe. Defaults to 15.                                                                        | `nil`    |
| `aperture.controller.livenessProbe.timeoutSeconds`                    | Timeout seconds for livenessProbe. Defaults to 5.                                                                        | `nil`    |
| `aperture.controller.livenessProbe.failureThreshold`                  | Failure threshold for livenessProbe. Defaults to 6.                                                                      | `nil`    |
| `aperture.controller.livenessProbe.successThreshold`                  | Success threshold for livenessProbe. Defaults to 1.                                                                      | `nil`    |
| `aperture.controller.readinessProbe.enabled`                          | Enable readinessProbe on Controller containers                                                                           | `true`   |
| `aperture.controller.readinessProbe.initialDelaySeconds`              | Initial delay seconds for readinessProbe. Defaults to 15.                                                                | `nil`    |
| `aperture.controller.readinessProbe.periodSeconds`                    | Period seconds for readinessProbe. Defaults to 15.                                                                       | `nil`    |
| `aperture.controller.readinessProbe.timeoutSeconds`                   | Timeout seconds for readinessProbe. Defaults to 5.                                                                       | `nil`    |
| `aperture.controller.readinessProbe.failureThreshold`                 | Failure threshold for readinessProbe. Defaults to 6.                                                                     | `nil`    |
| `aperture.controller.readinessProbe.successThreshold`                 | Success threshold for readinessProbe. Defaults to 1.                                                                     | `nil`    |
| `aperture.controller.customLivenessProbe`                             | Custom livenessProbe that overrides the default one                                                                      | `{}`     |
| `aperture.controller.customReadinessProbe`                            | Custom readinessProbe that overrides the default one                                                                     | `{}`     |
| `aperture.controller.resources.limits`                                | The resources limits for the Controller containers                                                                       | `{}`     |
| `aperture.controller.resources.requests`                              | The requested resources for the Controller containers                                                                    | `{}`     |
| `aperture.controller.podSecurityContext.enabled`                      | Enabled Controller pods' Security Context                                                                                | `false`  |
| `aperture.controller.podSecurityContext.fsGroup`                      | Set Controller pod's Security Context fsGroup. Defaults to 1001.                                                         | `1001`   |
| `aperture.controller.containerSecurityContext.enabled`                | Enabled Controller containers' Security Context                                                                          | `false`  |
| `aperture.controller.containerSecurityContext.runAsUser`              | Set Controller containers' Security Context runAsUser. Defaults to 1001.                                                 | `1001`   |
| `aperture.controller.containerSecurityContext.runAsNonRoot`           | Set Controller containers' Security Context runAsNonRoot. Defaults to false.                                             | `false`  |
| `aperture.controller.containerSecurityContext.readOnlyRootFilesystem` | Set Controller containers' Security Context runAsNonRoot. Defaults to false.                                             | `false`  |
| `aperture.controller.command`                                         | Override default container command (useful when using custom images)                                                     | `[]`     |
| `aperture.controller.args`                                            | Override default container args (useful when using custom images)                                                        | `[]`     |
| `aperture.controller.hostAliases`                                     | Controller pods host aliases                                                                                             | `[]`     |
| `aperture.controller.podLabels`                                       | Extra labels for Controller pods                                                                                         | `{}`     |
| `aperture.controller.podAnnotations`                                  | Annotations for Controller pods                                                                                          | `{}`     |
| `aperture.controller.affinity`                                        | Affinity for Controller pods assignment                                                                                  | `{}`     |
| `aperture.controller.nodeSelector`                                    | Node labels for Controller pods assignment                                                                               | `{}`     |
| `aperture.controller.tolerations`                                     | Tolerations for Controller pods assignment                                                                               | `[]`     |
| `aperture.controller.terminationGracePeriodSeconds`                   | configures how long kubelet gives Controller chart to terminate cleanly                                                  | `nil`    |
| `aperture.controller.lifecycleHooks`                                  | for the Controller container(s) to automate configuration before or after startup                                        | `{}`     |
| `aperture.controller.extraEnvVars`                                    | Array with extra environment variables to add to Controller nodes                                                        | `[]`     |
| `aperture.controller.extraEnvVarsCM`                                  | Name of existing ConfigMap containing extra env vars for Controller nodes                                                | `""`     |
| `aperture.controller.extraEnvVarsSecret`                              | Name of existing Secret containing extra env vars for Controller nodes                                                   | `""`     |
| `aperture.controller.extraVolumes`                                    | Optionally specify extra list of additional volumes for the Controller pod(s)                                            | `[]`     |
| `aperture.controller.extraVolumeMounts`                               | Optionally specify extra list of additional volumeMounts for the Controller container(s)                                 | `[]`     |
| `aperture.controller.sidecars`                                        | Add additional sidecar containers to the Controller pod(s)                                                               | `[]`     |
| `aperture.controller.initContainers`                                  | Add additional init containers to the Controller pod(s)                                                                  | `[]`     |
| `aperture.service.agent.annotations`                                  | Additional custom annotations for Agent service                                                                          | `{}`     |
| `aperture.service.controller.annotations`                             | Additional custom annotations for Controller service                                                                     | `{}`     |
| `aperture.sidecar.enabled`                                            | Enables sidecar mode for the Agent                                                                                       | `false`  |
| `aperture.sidecar.enableNamespacesByDefault`                          | List of namespaces in which sidecar injection will be enabled when Sidecar mode is enabled.                              | `[]`     |
| `aperture.etcd.endpoints`                                             | List of Etcd server endpoints. Example, ["https://etcd:2379"]. This must not be empty when etcd.enabled is set to false. | `[]`     |
| `aperture.etcd.leaseTtl`                                              | Lease time-to-live.                                                                                                      | `60s`    |
| `aperture.prometheus.address`                                         | specifies the address of the Prometheus server. This must not be empty when prometheus.enabled is set to false.          | `nil`    |


### etcd

| Name                    | Description                                                  | Value    |
| ----------------------- | ------------------------------------------------------------ | -------- |
| `etcd.enabled`          | Whether to deploy a small etcd cluster as part of this chart | `true`   |
| `etcd.auth.rbac.create` | specifies whether to create the RBAC resources for Etcd      | `false`  |
| `etcd.auth.token.type`  | specifies the type of token to use                           | `simple` |


### prometheus

| Name                 | Description                                     | Value  |
| -------------------- | ----------------------------------------------- | ------ |
| `prometheus.enabled` | specifies whether to deploy embedded prometheus | `true` |


### istio

| Name                                 | Description                                         | Value                   |
| ------------------------------------ | --------------------------------------------------- | ----------------------- |
| `istio.enabled`                      | specifies whether to deploy embedded istio          | `true`                  |
| `istio.global.tag`                   | specifies the image tag for all the istio resources | `1.11.8`                |
| `istio.envoyFilter.install`          | specifies whether to deploy EnvoyFilter             | `true`                  |
| `istio.envoyFilter.name`             | specifies name for the EnvoyFilter                  | `aperture-envoy-filter` |
| `istio.envoyFilter.authzGrpcTimeout` | specifies timeout for the AuthZ gRPC connection     | `0.01s`                 |
| `istio.envoyFilter.maxRequestBytes`  | specifies allowed maximum request bytes             | `8192`                  |

