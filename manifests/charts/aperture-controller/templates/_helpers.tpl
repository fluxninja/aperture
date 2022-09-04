{{/*
Create the name of the service account to use
*/}}
{{- define "aperture-operator.serviceAccountName" -}}
{{- default (include "common.names.fullname" .) .Values.operator.serviceAccount.name }}
{{- end }}

{{/*
Compile all warnings into a single message.
*/}}
{{- define "aperture-operator.validateValues" -}}
{{- $messages := list -}}
{{- $messages := without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{-   printf "\nVALUES VALIDATION:\n%s" $message -}}
{{- end -}}
{{- end -}}

{{/*
Create the endpoint of the etcd for Aperture Controller
{{ include "aperture.etcd.endpoints" ( dict "etcd" .Values.path.to.the.etcd "context" $.context $) }}
*/}}
{{- define "aperture.etcd.endpoints" -}}
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
{{ include "aperture.prometheus.address" ( dict "prometheus" .Values.path.to.the.prometheus "context" $.context $) }}
*/}}
{{- define "aperture.prometheus.address" -}}
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
{{ include "aperture.fluxninjaPlugin.endpoint" ( dict "aperture" .Values.path.to.the.aperture $) }}
*/}}
{{- define "aperture.fluxninjaPlugin.endpoint" -}}
{{- if .aperture.fluxninjaPlugin.enabled -}}
    {{- if .aperture.fluxninjaPlugin.endpoint -}}
        {{ print .aperture.fluxninjaPlugin.endpoint }}
    {{- else -}}
        {{- fail "Value of endpoint for FluxNinja plugin cannot be empty when .Values.aperture.fluxninjaPlugin.enabled is set to true." -}}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}

{{/*
Fetch the value of the API Key secret for Aperture Controller
{{ include "aperture.apiSecret.value" ( dict "object" .Values.path.to.the.agent/controller $) }}
*/}}
{{- define "aperture.apisecret.value" -}}
{{- if .object.fluxninjaPlugin.enabled -}}
    {{- if .object.fluxninjaPlugin.apiKeySecret.create -}}
        {{- if .object.fluxninjaPlugin.apiKeySecret.value -}}
            {{ print .object.fluxninjaPlugin.apiKeySecret.value }}
        {{- else -}}
            {{- fail "Value of API Key for Controller cannot be empty when .Values.agent/controller.fluxninjaPlugin.enabled and .Values.agent/controller.fluxninjaPlugin.apiKeySecret.create is set to true." -}}
        {{- end -}}
    {{- else -}}
        {{ print "" }}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}
