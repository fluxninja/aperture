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
Create the endpoint of the etcd for Aperture Agent
{{ include "aperture.etcd.endpoints" ( dict "etcd" .Values.path.to.the.etcd "context" $.context $) }}
*/}}
{{- define "aperture.etcd.endpoints" -}}
{{- $endpoints := list -}}
{{ $endpoints = without .etcd.endpoints "" }}
{{- if empty $endpoints -}}
    {{- fail "Value for etcd endpoints of Agent cannot be empty." -}}
{{- end -}}
{{ print $endpoints }}
{{- end -}}

{{/*
Create the address of the Prometheus for Aperture Agent
{{ include "aperture.prometheus.address" ( dict "prometheus" .Values.path.to.the.prometheus "context" $.context $) }}
*/}}
{{- define "aperture.prometheus.address" -}}
{{- if .prometheus.address -}}
    {{ print .prometheus.address }}
{{- else -}}
    {{- fail "Value for prometheus address of Agent cannot be empty." -}}
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
Fetch the value of the API Key secret for Aperture Agent
{{ include "aperture.apiSecret.value" ( dict "object" .Values.path.to.the.agent/controller $) }}
*/}}
{{- define "aperture.apisecret.value" -}}
{{- if .object.fluxninjaPlugin.enabled -}}
    {{- if .object.fluxninjaPlugin.apiKeySecret.create -}}
        {{- if .object.fluxninjaPlugin.apiKeySecret.value -}}
            {{ print .object.fluxninjaPlugin.apiKeySecret.value }}
        {{- else -}}
            {{- fail "Value of API Key for Agent cannot be empty when .Values.agent/controller.fluxninjaPlugin.enabled and .Values.agent/controller.fluxninjaPlugin.apiKeySecret.create is set to true." -}}
        {{- end -}}
    {{- else -}}
        {{ print "" }}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}
