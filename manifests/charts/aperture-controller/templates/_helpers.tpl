{{/*
Create the name of the service account to use
*/}}
{{- define "controller-operator.serviceAccountName" -}}
{{- default ( print (include "common.names.fullname" .) "-operator" ) .Values.operator.serviceAccount.name }}
{{- end }}

{{/*
Compile all warnings into a single message.
*/}}
{{- define "controller-operator.validateValues" -}}
{{- $messages := list -}}
{{- $messages := without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{-   printf "\nVALUES VALIDATION:\n%s" $message -}}
{{- end -}}
{{- end -}}

{{/*
Get image tag for operator.
{{ include "controller-operator.image" ( dict "image" .Values.path.to.the.image "context" $.context $ ) }}
*/}}
{{- define "controller-operator.image" -}}
{{- $globalAzure := get .context.Values.global "azure" -}}
{{- if not (empty $globalAzure) -}}
    {{- $azureImage := "" -}}
    {{- if  .operator -}}
        {{- $azureImage = (printf "%s/%s@%s" .context.Values.global.azure.images.operator.registry  .context.Values.global.azure.images.operator.image .context.Values.global.azure.images.operator.digest) -}}
    {{- else -}}
        {{- $azureImage = (printf "%s/%s@%s" .context.Values.global.azure.images.controller.registry  .context.Values.global.azure.images.controller.image .context.Values.global.azure.images.controller.digest) -}}
    {{- end -}}
    {{ print $azureImage }}
{{- else -}}
    {{- $tag := get .image "tag" -}}
    {{- $newImage := .image -}}
    {{- if (not $tag) -}}
        {{- $tag = trimPrefix "v" .context.Chart.AppVersion -}}
    {{- end -}}
    {{- $_ := set $newImage "tag" $tag -}}
    {{ print (include "common.images.image" (dict "imageRoot" $newImage "global" .context.Values.global)) }}
{{- end -}}
{{- end -}}

{{/*
Check of etcd at global level.
{{ include "etcd.initContainer.image" ( dict "image" .Values.path.to.the.image "context" $.context $ ) }}
*/}}
{{- define "etcd.initContainer.image" -}}
{{- $globalAzure := get .context.Values.global "azure" -}}
{{- if not (empty $globalAzure) -}}
    {{- $etcdAzure := get .context.Values.global.azure.images "etcd" -}}
    {{- if not (empty $etcdAzure) -}}
        {{- $azureImage := (printf "%s/%s@%s" .context.Values.global.azure.images.etcd.registry  .context.Values.global.azure.images.etcd.image .context.Values.global.azure.images.etcd.digest) -}}
        {{ print $azureImage }}
    {{- end -}}
{{- else -}}
    {{- $newImage := .image -}}
    {{ print (include "common.images.image" (dict "imageRoot" $newImage "global" .context.Values.global)) }}
{{- end -}}
{{- end -}}

{{/*
Check of prometheus at global level.
{{ include "etcd.initContainer.image" ( dict "image" .Values.path.to.the.image "context" $.context $ ) }}
*/}}
{{- define "prometheus.initContainer.image" -}}
{{- $globalAzure := get .context.Values.global "azure" -}}
{{- if not (empty $globalAzure) -}}
    {{- $prometheusAzure := get .context.Values.global.azure.images "yq" -}}
    {{- if not (empty $globalAzure) -}}
        {{- $prometheusImage := (printf "%s/%s@%s" .context.Values.global.azure.images.yq.registry  .context.Values.global.azure.images.yq.image .context.Values.global.azure.images.yq.digest) -}}
        {{ print $prometheusImage }}
    {{- end -}}
{{- else -}}
    {{- $newImage := .image -}}
    {{ print (include "common.images.image" (dict "imageRoot" $newImage "global" .context.Values.global)) }}
{{- end -}}
{{- end -}}

{{/*
Create the endpoint of the etcd for Aperture Controller
{{ include "controller.etcd.endpoints" ( dict "etcd" .Values.path.to.the.etcd "context" $.context $ ) }}
*/}}
{{- define "controller.etcd.endpoints" -}}
{{- $endpoints := list -}}
{{- if .context.Values.etcd.enabled -}}
    {{- $endpoints = append $endpoints (printf "http://%s-etcd:2379" .context.Release.Name) -}}
{{- else -}}
    {{ $endpoints = without .etcd.endpoints "" }}
    {{- if empty $endpoints -}}
        {{- fail "Value for etcd endpoints of Controller cannot be empty when .Values.etcd.enbled is set to false." -}}
    {{- end -}}
{{- end -}}
{{ print $endpoints }}
{{- end -}}

{{/*
Create the address of the Prometheus for Aperture Controller
{{ include "controller.prometheus.address" ( dict "prometheus" .Values.path.to.the.prometheus "context" $.context $ ) }}
*/}}
{{- define "controller.prometheus.address" -}}
{{- if .context.Values.prometheus.enabled -}}
    {{- printf "http://%s-prometheus-server:80" .context.Release.Name -}}
{{- else -}}
    {{- if .prometheus.address -}}
        {{ print .prometheus.address }}
    {{- else -}}
        {{- fail "Value for prometheus address of Controller cannot be empty when .Values.prometheus.enbled is set to false." -}}
    {{- end -}}
{{- end -}}
{{- end -}}

{{/*
Fetch the value of the API Key secret for Aperture Controller
{{ include "controller.apiSecret.value" ( dict "controller" .Values.path.to.the.controller $ ) }}
*/}}
{{- define "controller.apiSecret.value" -}}
{{- if .controller.secrets.fluxNinjaExtension.create -}}
    {{- if .controller.secrets.fluxNinjaExtension.value -}}
        {{ print .controller.secrets.fluxNinjaExtension.value }}
    {{- else -}}
        {{- fail "Value of API Key for Controller cannot be empty when .Values.controller.secrets.fluxNinjaExtension.create is set to true." -}}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}

{{/*
Fetch the Name of the API Key secret for Aperture Controller
{{ include "controller.apiSecret.name" ( dict "controller" .Values.path.to.the.controller "context" $.context $ ) }}
*/}}
{{- define "controller.apiSecret.name" -}}
{{- if .controller.secrets.fluxNinjaExtension.secretKeyRef.name -}}
    {{ print .controller.secrets.fluxNinjaExtension.secretKeyRef.name }}
{{- else -}}
    {{ printf "%s-controller-apikey" .context.Release.Name }}
{{- end -}}
{{- end -}}

{{/*
Fetch the Key of the API Key secret for Aperture Controller
{{ include "controller.apiSecret.name" ( dict "controller" .Values.path.to.the.controller $ ) }}
*/}}
{{- define "controller.apiSecret.key" -}}
{{- if .controller.secrets.fluxNinjaExtension.secretKeyRef.key -}}
    {{ print .controller.secrets.fluxNinjaExtension.secretKeyRef.key }}
{{- else -}}
    {{ print "apiKey" }}
{{- end -}}
{{- end -}}

{{/*
Fetch the server port of the Aperture Controller
{{ include "controller.server.port" ( dict "controller" .Values.path.to.the.controller $ ) }}
*/}}
{{- define "controller.server.port" -}}
{{- if and .controller.config .controller.config.server .controller.config.server.listener .controller.config.server.listener.addr -}}
    {{ print (split ":" .controller.config.server.listener.addr)._1 }}
{{- else -}}
    {{ print "8080" }}
{{- end -}}
{{- end -}}

{{/*
Fetch the OTEL port of the Aperture Controller
{{ include "controller.otel.port" ( dict "controller" .Values.path.to.the.controller portName string defaultPort string $ ) }}
*/}}
{{- define "controller.otel.port" -}}
{{- if and .controller.config .controller.config.otel .controller.config.otel.ports (hasKey .controller.config.otel.ports .portName) -}}
    {{ print (get .controller.config.otel.ports .portName) }}
{{- else -}}
    {{ print .defaultPort }}
{{- end -}}
{{- end -}}

{{/*
Prepare the Host for configuring in the Ingress
{{ include "controller.ingress-endpoint" ( dict "component" "component_name" "context" $.context $ ) }}
*/}}
{{- define "controller.ingress-endpoint" -}}
{{- if .context.Values.ingress.domain_name -}}
    {{- printf "%s.%s" .component .context.Values.ingress.domain_name -}}
{{- else -}}
    {{- fail "Value of .Values.ingress.domain_name cannot be empty when .Values.ingress.enabled is set to true." -}}
{{- end -}}
{{- end -}}

{{/*
Add the pod labels when global azure field is set
*/}}
{{- define "controller.podlabels" -}}
{{- $globalAzure := get .context.Values.global "azure" -}}
{{- $podLabels := "" -}}
{{- if not (empty $globalAzure) -}}
    {{- $podLabels = (printf "%s: %s" "azure-extensions-usage-release-identifier" .context.Release.Name ) -}}
    {{ print $podLabels | toYaml | nindent 4 | replace "'" ""}}
{{- else -}}
  {{- if not (empty .podlabels) -}}
    {{- $podLabels = .podlabels | toYaml | nindent 4 -}}
    {{ print $podLabels }}
  {{- end -}}
{{- end -}}
{{- end -}}
