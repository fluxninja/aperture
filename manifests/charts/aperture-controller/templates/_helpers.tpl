{{/*
Create the name of the service account to use
*/}}
{{- define "controller-operator.serviceAccountName" -}}
{{- default (include "common.names.fullname" .) .Values.operator.serviceAccount.name }}
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
{{ include "controller-operator.image" ( dict "image" .Values.path.to.the.image "context" $.context $) }}
*/}}
{{- define "controller-operator.image" -}}
{{- $tag := get .image "tag" -}}
{{- $newImage := .image -}}
{{- if (not $tag) -}}
    {{- $tag = trimPrefix "v" .context.Chart.AppVersion -}}
{{- end -}}
{{- $_ := set $newImage "tag" $tag -}}
{{ print (include "common.images.image" (dict "imageRoot" $newImage "global" .context.Values.global)) }}
{{- end -}}

{{/*
Create the endpoint of the etcd for Aperture Controller
{{ include "controller.etcd.endpoints" ( dict "etcd" .Values.path.to.the.etcd "context" $.context $) }}
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
{{ include "controller.prometheus.address" ( dict "prometheus" .Values.path.to.the.prometheus "context" $.context $) }}
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
Fetch the endpoint of the FluxNinja cloud instance
{{ include "controller.fluxNinjaPlugin.endpoint" ( dict "controller" .Values.path.to.the.controller $) }}
*/}}
{{- define "controller.fluxNinjaPlugin.endpoint" -}}
{{- if .controller.config.fluxninja_plugin.enabled -}}
    {{- if .controller.config.fluxninja_plugin.endpoint -}}
        {{ print .controller.config.fluxninja_plugin.endpoint }}
    {{- else -}}
        {{- fail "Value of endpoint for FluxNinja plugin cannot be empty when .Values.controller.config.fluxninja_plugin.enabled is set to true." -}}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}

{{/*
Fetch the value of the API Key secret for Aperture Controller
{{ include "controller.apiSecret.value" ( dict "controller" .Values.path.to.the.controller $) }}
*/}}
{{- define "controller.apisecret.value" -}}
{{- if .controller.secrets.fluxNinjaPlugin.create -}}
    {{- if .controller.secrets.fluxNinjaPlugin.value -}}
        {{ print .controller.secrets.fluxNinjaPlugin.value }}
    {{- else -}}
        {{- fail "Value of API Key for Controller cannot be empty when .Values.controller.secrets.fluxNinjaPlugin.create is set to true." -}}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}

{{/*
Prepare the Host for configuring in the Ingress
{{ include "controller.ingress-endpoint" ( dict "component" "component_name" "context" $.context $) }}
*/}}
{{- define "controller.ingress-endpoint" -}}
{{- if .context.Values.ingress.domain_name -}}
    {{- printf "%s.%s" .component .context.Values.ingress.domain_name -}}
{{- else -}}
    {{- fail "Value of .Values.ingress.domain_name cannot be empty when .Values.ingress.enabled is set to true." -}}
{{- end -}}
{{- end -}}
