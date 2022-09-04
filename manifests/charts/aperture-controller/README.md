# Aperture Controller

## Parameters

### Global Parameters

| Name                      | Description                                     | Value |
| ------------------------- | ----------------------------------------------- | ----- |
| `global.imageRegistry`    | Global Docker image registry                    | `nil` |
| `global.imagePullSecrets` | Global Docker registry secret names as an array | `[]`  |


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

| Name                                                         | Description                                                                                                            | Value                 |
| ------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------- | --------------------- |
| `operator.image.registry`                                    | Operator image registry                                                                                                | `docker.io/fluxninja` |
| `operator.image.repository`                                  | Operator image repository                                                                                              | `aperture-operator`   |
| `operator.image.tag`                                         | Operator image tag (immutable tags are recommended)                                                                    | `latest`              |
| `operator.image.pullPolicy`                                  | Operator image pull policy                                                                                             | `Never`               |
| `operator.image.pullSecrets`                                 | Operator image pull secrets                                                                                            | `[]`                  |
| `operator.replicaCount`                                      | Number of replicas for Operator deployment                                                                             | `1`                   |
| `operator.podAnnotations`                                    | Pod annotations                                                                                                        | `{}`                  |
| `operator.podLabels`                                         | Additional pod labels                                                                                                  | `{}`                  |
| `operator.updateStrategy.type`                               | Set up update strategy for Aperture Operator installation.                                                             | `Recreate`            |
| `operator.priorityClassName`                                 | pods' priorityClassName                                                                                                | `""`                  |
| `operator.topologySpreadConstraints`                         | Topology Spread Constraints for pod assignment                                                                         | `[]`                  |
| `operator.schedulerName`                                     | Alternative scheduler                                                                                                  | `""`                  |
| `operator.hostAliases`                                       | Add deployment host aliases                                                                                            | `[]`                  |
| `operator.nodeSelector`                                      | Node labels for pod assignment                                                                                         | `{}`                  |
| `operator.podAffinityPreset`                                 | Pod affinity preset                                                                                                    | `""`                  |
| `operator.podAntiAffinityPreset`                             | Pod anti-affinity preset                                                                                               | `soft`                |
| `operator.nodeAffinityPreset.type`                           | Set nodeAffinity preset type                                                                                           | `""`                  |
| `operator.nodeAffinityPreset.key`                            | Set nodeAffinity preset key                                                                                            | `""`                  |
| `operator.nodeAffinityPreset.values`                         | Set nodeAffinity preset values                                                                                         | `[]`                  |
| `operator.affinity`                                          | Affinity for controller pod assignment                                                                                 | `{}`                  |
| `operator.tolerations`                                       | Tolerations for controller pod assignment                                                                              | `[]`                  |
| `operator.podSecurityContext.enabled`                        | Enable pods security context                                                                                           | `true`                |
| `operator.podSecurityContext.runAsUser`                      | User ID for the pods                                                                                                   | `1001`                |
| `operator.podSecurityContext.runAsGroup`                     | User ID for the pods                                                                                                   | `1001`                |
| `operator.podSecurityContext.runAsNonRoot`                   | Aperture Operator must run as nonRoot                                                                                  | `true`                |
| `operator.podSecurityContext.fsGroup`                        | Group ID for the pods                                                                                                  | `1001`                |
| `operator.podSecurityContext.supplementalGroups`             | Which group IDs containers add                                                                                         | `[]`                  |
| `operator.containerSecurityContext.enabled`                  | Enable container security context                                                                                      | `true`                |
| `operator.containerSecurityContext.runAsUser`                | User ID for the operator container                                                                                     | `1001`                |
| `operator.containerSecurityContext.runAsGroup`               | User ID for the operator container                                                                                     | `1001`                |
| `operator.containerSecurityContext.runAsNonRoot`             | Force the container to be run as non-root                                                                              | `true`                |
| `operator.containerSecurityContext.privileged`               | Decide if the container runs privileged.                                                                               | `false`               |
| `operator.containerSecurityContext.readOnlyRootFilesystem`   | ReadOnlyRootFilesystem fot the operator container                                                                      | `false`               |
| `operator.containerSecurityContext.allowPrivilegeEscalation` | Allow Privilege Escalation for the operator container                                                                  | `false`               |
| `operator.terminationGracePeriodSeconds`                     | In seconds, time the given to the pod needs to terminate gracefully                                                    | `10`                  |
| `operator.initContainers`                                    | Add additional init containers to the etcd pods                                                                        | `[]`                  |
| `operator.command`                                           | Default container command (useful when using custom images)                                                            | `[]`                  |
| `operator.args`                                              | Default container args (useful when using custom images)                                                               | `[]`                  |
| `operator.lifecycleHooks`                                    | for the aperture-operator container to automate configuration before or after startup                                  | `{}`                  |
| `operator.extraEnvVars`                                      | Array with extra environment variables to add to RabbitMQ Cluster Operator nodes                                       | `[]`                  |
| `operator.extraEnvVarsCM`                                    | Name of existing ConfigMap containing extra env vars for RabbitMQ Cluster Operator nodes                               | `""`                  |
| `operator.extraEnvVarsSecret`                                | Name of existing Secret containing extra env vars for RabbitMQ Cluster Operator nodes                                  | `""`                  |
| `operator.resources`                                         | Container resource requests and limits                                                                                 | `{}`                  |
| `operator.livenessProbe.enabled`                             | Enable livenessProbe                                                                                                   | `true`                |
| `operator.livenessProbe.initialDelaySeconds`                 | Initial delay seconds for livenessProbe                                                                                | `10`                  |
| `operator.livenessProbe.periodSeconds`                       | Period seconds for livenessProbe                                                                                       | `10`                  |
| `operator.livenessProbe.timeoutSeconds`                      | Timeout seconds for livenessProbe                                                                                      | `1`                   |
| `operator.livenessProbe.failureThreshold`                    | Failure threshold for livenessProbe                                                                                    | `3`                   |
| `operator.livenessProbe.successThreshold`                    | Success threshold for livenessProbe                                                                                    | `1`                   |
| `operator.readinessProbe.enabled`                            | Enable readinessProbe                                                                                                  | `true`                |
| `operator.readinessProbe.initialDelaySeconds`                | Initial delay seconds for readinessProbe                                                                               | `10`                  |
| `operator.readinessProbe.periodSeconds`                      | Period seconds for readinessProbe                                                                                      | `10`                  |
| `operator.readinessProbe.timeoutSeconds`                     | Timeout seconds for readinessProbe                                                                                     | `1`                   |
| `operator.readinessProbe.failureThreshold`                   | Failure threshold for readinessProbe                                                                                   | `3`                   |
| `operator.readinessProbe.successThreshold`                   | Success threshold for readinessProbe                                                                                   | `1`                   |
| `operator.startupProbe.enabled`                              | Enable startupProbe                                                                                                    | `true`                |
| `operator.startupProbe.initialDelaySeconds`                  | Initial delay seconds for startupProbe                                                                                 | `10`                  |
| `operator.startupProbe.periodSeconds`                        | Period seconds for startupProbe                                                                                        | `10`                  |
| `operator.startupProbe.timeoutSeconds`                       | Timeout seconds for startupProbe                                                                                       | `1`                   |
| `operator.startupProbe.failureThreshold`                     | Failure threshold for startupProbe                                                                                     | `3`                   |
| `operator.startupProbe.successThreshold`                     | Success threshold for startupProbe                                                                                     | `1`                   |
| `operator.customLivenessProbe`                               | Override default liveness probe                                                                                        | `{}`                  |
| `operator.customReadinessProbe`                              | Override default readiness probe                                                                                       | `{}`                  |
| `operator.customStartupProbe`                                | Override default startup probe                                                                                         | `{}`                  |
| `operator.extraVolumes`                                      | Optionally specify extra list of additional volumes                                                                    | `[]`                  |
| `operator.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts                                                               | `[]`                  |
| `operator.rbac.create`                                       | Create specifies whether to install and use RBAC rules                                                                 | `true`                |
| `operator.serviceAccount.create`                             | Specifies whether a service account should be created                                                                  | `true`                |
| `operator.serviceAccount.name`                               | The name of the service account to use. If not set and create is true, a name is generated using the fullname template | `""`                  |
| `operator.serviceAccount.annotations`                        | Add annotations                                                                                                        | `{}`                  |
| `operator.serviceAccount.automountServiceAccountToken`       | Automount API credentials for a service account.                                                                       | `true`                |


### Controller Custom Resource Parameters

| Name                                                         | Description                                                                                                                                                       | Value    |
| ------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `controller.create`                                          | Specifies whether a CR for Controller should be created                                                                                                           | `true`   |
| `controller.fluxninjaPlugin.enabled`                         | Boolean flag for enabling FluxNinja cloud connection from Controller                                                                                              | `false`  |
| `controller.fluxninjaPlugin.endpoint`                        | FluxNinja cloud instance endpoint                                                                                                                                 | `nil`    |
| `controller.fluxninjaPlugin.heartbeatsInterval`              | specifies how often to send heartbeats to the cloud. Defaults to '30s'.                                                                                           | `nil`    |
| `controller.fluxninjaPlugin.tls.insecure`                    | specifies whether to communicate with FluxNinja cloud over TLS or in plain text. Defaults to false.                                                               | `nil`    |
| `controller.fluxninjaPlugin.tls.insecureSkipVerify`          | specifies whether to verify FluxNinja cloud certificate. Defaults to false.                                                                                       | `nil`    |
| `controller.fluxninjaPlugin.tls.caFile`                      | specifies an alternative CA certificates bundle to use to validate FluxNinja cloud certificate                                                                    | `nil`    |
| `controller.fluxninjaPlugin.apiKeySecret.create`             | Whether to create Kubernetes Secret with provided Controller API Key.                                                                                             | `true`   |
| `controller.fluxninjaPlugin.apiKeySecret.secretKeyRef.name`  | specifies a name of the Secret for Controller API Key to be used. This defaults to {{ .Release.Name }}-controller-apikey                                          | `nil`    |
| `controller.fluxninjaPlugin.apiKeySecret.secretKeyRef.key`   | specifies which key from the Secret for Controller API Key to use                                                                                                 | `apiKey` |
| `controller.fluxninjaPlugin.apiKeySecret.value`              | API Key to use when creating a new Controller API Key Secret                                                                                                      | `nil`    |
| `controller.image.registry`                                  | Controller image registry. Defaults to 'docker.io/fluxninja'.                                                                                                     | `nil`    |
| `controller.image.repository`                                | Controller image repository. Defaults to 'aperture-controller'.                                                                                                   | `nil`    |
| `controller.image.tag`                                       | Controller image tag (immutable tags are recommended). Defaults to 'latest'.                                                                                      | `nil`    |
| `controller.image.pullPolicy`                                | Controller image pull policy. Defaults to 'IfNotPresent'.                                                                                                         | `nil`    |
| `controller.image.pullSecrets`                               | Controller image pull secrets                                                                                                                                     | `[]`     |
| `controller.serverPort`                                      | The Controller's server port. Defaults to 80.                                                                                                                     | `nil`    |
| `controller.hostAliases`                                     | Add deployment host aliases for Controller deployment                                                                                                             | `[]`     |
| `controller.log.prettyConsole`                               | Additional log writer: pretty console (stdout) logging (not recommended for prod environments). Defaults to false.                                                | `nil`    |
| `controller.log.nonBlocking`                                 | Use non-blocking log writer (can lose logs at high throughput). Defaults to True.                                                                                 | `nil`    |
| `controller.log.level`                                       | Log level. Keywords allowed - ["debug", "info", "warn", "fatal", "panic", "trace"]. Defaults to 'info'.                                                           | `nil`    |
| `controller.log.file`                                        | Output file for logs. Keywords allowed - ["stderr", "stderr", "default"]. Defaults to 'stderr'.                                                                   | `nil`    |
| `controller.serviceAccount.create`                           | Specifies whether a ServiceAccount should be created                                                                                                              | `true`   |
| `controller.serviceAccount.annotations`                      | Additional Service Account annotations (evaluated as a template)                                                                                                  | `{}`     |
| `controller.serviceAccount.automountServiceAccountToken`     | Automount service account token for the server service account. Defaults to true                                                                                  | `nil`    |
| `controller.livenessProbe.enabled`                           | Enable livenessProbe on Controller containers                                                                                                                     | `true`   |
| `controller.livenessProbe.initialDelaySeconds`               | Initial delay seconds for livenessProbe. Defaults to 15.                                                                                                          | `nil`    |
| `controller.livenessProbe.periodSeconds`                     | Period seconds for livenessProbe. Defaults to 15.                                                                                                                 | `nil`    |
| `controller.livenessProbe.timeoutSeconds`                    | Timeout seconds for livenessProbe. Defaults to 5.                                                                                                                 | `nil`    |
| `controller.livenessProbe.failureThreshold`                  | Failure threshold for livenessProbe. Defaults to 6.                                                                                                               | `nil`    |
| `controller.livenessProbe.successThreshold`                  | Success threshold for livenessProbe. Defaults to 1.                                                                                                               | `nil`    |
| `controller.readinessProbe.enabled`                          | Enable readinessProbe on Controller containers                                                                                                                    | `true`   |
| `controller.readinessProbe.initialDelaySeconds`              | Initial delay seconds for readinessProbe. Defaults to 15.                                                                                                         | `nil`    |
| `controller.readinessProbe.periodSeconds`                    | Period seconds for readinessProbe. Defaults to 15.                                                                                                                | `nil`    |
| `controller.readinessProbe.timeoutSeconds`                   | Timeout seconds for readinessProbe. Defaults to 5.                                                                                                                | `nil`    |
| `controller.readinessProbe.failureThreshold`                 | Failure threshold for readinessProbe. Defaults to 6.                                                                                                              | `nil`    |
| `controller.readinessProbe.successThreshold`                 | Success threshold for readinessProbe. Defaults to 1.                                                                                                              | `nil`    |
| `controller.customLivenessProbe`                             | Custom livenessProbe that overrides the default one                                                                                                               | `{}`     |
| `controller.customReadinessProbe`                            | Custom readinessProbe that overrides the default one                                                                                                              | `{}`     |
| `controller.resources.limits`                                | The resources limits for the Controller containers                                                                                                                | `{}`     |
| `controller.resources.requests`                              | The requested resources for the Controller containers                                                                                                             | `{}`     |
| `controller.podSecurityContext.enabled`                      | Enabled Controller pods' Security Context                                                                                                                         | `false`  |
| `controller.podSecurityContext.fsGroup`                      | Set Controller pod's Security Context fsGroup. Defaults to 1001.                                                                                                  | `nil`    |
| `controller.containerSecurityContext.enabled`                | Enabled Controller containers' Security Context. Defaults to false.                                                                                               | `false`  |
| `controller.containerSecurityContext.runAsUser`              | Set Controller containers' Security Context runAsUser. Defaults to 1001.                                                                                          | `nil`    |
| `controller.containerSecurityContext.runAsNonRoot`           | Set Controller containers' Security Context runAsNonRoot. Defaults to false.                                                                                      | `nil`    |
| `controller.containerSecurityContext.readOnlyRootFilesystem` | Set Controller containers' Security Context runAsNonRoot. Defaults to false.                                                                                      | `nil`    |
| `controller.command`                                         | Override default container command (useful when using custom images)                                                                                              | `[]`     |
| `controller.args`                                            | Override default container args (useful when using custom images)                                                                                                 | `[]`     |
| `controller.podLabels`                                       | Extra labels for Controller pods                                                                                                                                  | `{}`     |
| `controller.podAnnotations`                                  | Annotations for Controller pods                                                                                                                                   | `{}`     |
| `controller.affinity`                                        | Affinity for Controller pods assignment                                                                                                                           | `{}`     |
| `controller.nodeSelector`                                    | Node labels for Controller pods assignment                                                                                                                        | `{}`     |
| `controller.tolerations`                                     | Tolerations for Controller pods assignment                                                                                                                        | `[]`     |
| `controller.terminationGracePeriodSeconds`                   | configures how long kubelet gives Controller chart to terminate cleanly                                                                                           | `nil`    |
| `controller.lifecycleHooks`                                  | for the Controller container(s) to automate configuration before or after startup                                                                                 | `{}`     |
| `controller.extraEnvVars`                                    | Array with extra environment variables to add to Controller nodes                                                                                                 | `[]`     |
| `controller.extraEnvVarsCM`                                  | Name of existing ConfigMap containing extra env vars for Controller nodes                                                                                         | `""`     |
| `controller.extraEnvVarsSecret`                              | Name of existing Secret containing extra env vars for Controller nodes                                                                                            | `""`     |
| `controller.extraVolumes`                                    | Optionally specify extra list of additional volumes for the Controller pod(s)                                                                                     | `[]`     |
| `controller.extraVolumeMounts`                               | Optionally specify extra list of additional volumeMounts for the Controller container(s)                                                                          | `[]`     |
| `controller.sidecars`                                        | Add additional sidecar containers to the Controller pod(s)                                                                                                        | `[]`     |
| `controller.initContainers`                                  | Add additional init containers to the Controller pod(s)                                                                                                           | `[]`     |
| `controller.service.annotations`                             | Additional custom annotations for Controller service                                                                                                              | `{}`     |
| `controller.sidecar.enabled`                                 | Enables sidecar mode for the Controller                                                                                                                           | `false`  |
| `controller.sidecar.enableNamespacesByDefault`               | List of namespaces in which sidecar injection will be enabled when Sidecar mode is enabled.                                                                       | `[]`     |
| `controller.etcd.endpoints`                                  | List of Etcd server endpoints. Example, ["https://etcd:2379"]. This must not be empty when etcd is not craeted by this chart (i.e. etcd.enabled is set to false). | `[]`     |
| `controller.etcd.leaseTtl`                                   | Lease time-to-live.                                                                                                                                               | `60s`    |
| `controller.prometheus.address`                              | specifies the address of the Prometheus server. This must not be empty when prometheus is not created by this chart (i.e. prometheus.enabled is set to false).    | `nil`    |
| `controller.otelConfig.grpcAddr`                             | GRPC listener addr for OTEL Collector. Defaults to ":4317"                                                                                                        | `nil`    |
| `controller.otelConfig.httpAddr`                             | HTTP listener addr for OTEL Collector. Defaults to ":4318"                                                                                                        | `nil`    |
| `controller.otelConfig.batchPrerollup.sendBatchSize`         | SendBatchSize is the size of a batch which after hit, will trigger it to be sent. Defaults to "10000".                                                            | `nil`    |
| `controller.otelConfig.batchPrerollup.timeout`               | Timeout sets the time after which a batch will be sent regardless of size. Defaults to "1s".                                                                      | `nil`    |
| `controller.otelConfig.batchPostrollup.sendBatchSize`        | SendBatchSize is the size of a batch which after hit, will trigger it to be sent. Defaults to "10000".                                                            | `nil`    |
| `controller.otelConfig.batchPostrollup.timeout`              | Timeout sets the time after which a batch will be sent regardless of size. Defaults to "1s".                                                                      | `nil`    |


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


