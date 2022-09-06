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
{{ include "controller.fluxninjaPlugin.endpoint" ( dict "controller" .Values.path.to.the.controller $) }}
*/}}
{{- define "controller.fluxninjaPlugin.endpoint" -}}
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
{{- if .controller.secrets.fluxninjaPlugin.create -}}
    {{- if .controller.secrets.fluxninjaPlugin.value -}}
        {{ print .controller.secrets.fluxninjaPlugin.value }}
    {{- else -}}
        {{- fail "Value of API Key for Controller cannot be empty when .Values.controller.secrets.fluxninjaPlugin.create is set to true." -}}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}
