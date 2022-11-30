{{/*
Create the name of the service account to use
*/}}
{{- define "agent-operator.serviceAccountName" -}}
{{- default (include "common.names.fullname" .) .Values.operator.serviceAccount.name }}
{{- end }}

{{/*
Compile all warnings into a single message.
*/}}
{{- define "agent-operator.validateValues" -}}
{{- $messages := list -}}
{{- $messages := without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{-   printf "\nVALUES VALIDATION:\n%s" $message -}}
{{- end -}}
{{- end -}}

{{/*
Get image tag for operator.
{{ include "agent-operator.image" ( dict "image" .Values.path.to.the.image "context" $.context $) }}
*/}}
{{- define "agent-operator.image" -}}
{{- $tag := get .image "tag" -}}
{{- $newImage := .image -}}
{{- if (not $tag) -}}
    {{- $tag = trimPrefix "v" .context.Chart.AppVersion -}}
{{- end -}}
{{- $_ := set $newImage "tag" $tag -}}
{{ print (include "common.images.image" (dict "imageRoot" $newImage "global" .context.Values.global)) }}
{{- end -}}

{{/*
Create the endpoint of the etcd for Aperture Agent
{{ include "agent.etcd.endpoints" ( dict "etcd" .Values.path.to.the.etcd "context" $.context $) }}
*/}}
{{- define "agent.etcd.endpoints" -}}
{{- $endpoints := list -}}
{{ $endpoints = without .etcd.endpoints "" }}
{{- if empty $endpoints -}}
    {{- fail "Value for etcd endpoints of Agent cannot be empty." -}}
{{- end -}}
{{ print $endpoints }}
{{- end -}}

{{/*
Create the address of the Prometheus for Aperture Agent
{{ include "agent.prometheus.address" ( dict "prometheus" .Values.path.to.the.prometheus "context" $.context $) }}
*/}}
{{- define "agent.prometheus.address" -}}
{{- if .prometheus.address -}}
    {{ print .prometheus.address }}
{{- else -}}
    {{- fail "Value for prometheus address of Agent cannot be empty." -}}
{{- end -}}
{{- end -}}

{{/*
Fetch the endpoint of the FluxNinja cloud instance
{{ include "agent.fluxNinjaPlugin.endpoint" ( dict "agent" .Values.path.to.the.agent $) }}
*/}}
{{- define "agent.fluxNinjaPlugin.endpoint" -}}
{{- if .agent.config.fluxninja_plugin.enabled -}}
    {{- if .agent.config.fluxninja_plugin.fluxninja_endpoint -}}
        {{ print .agent.config.fluxninja_plugin.fluxninja_endpoint }}
    {{- else -}}
        {{- fail "Value of fluxninja_endpoint for FluxNinja plugin cannot be empty when .Values.agent.config.fluxninja_plugin.enabled is set to true." -}}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}

{{/*
Fetch the value of the API Key secret for Aperture Agent
{{ include "agent.apiSecret.value" ( dict "agent" .Values.path.to.the.agent $) }}
*/}}
{{- define "agent.apisecret.value" -}}
{{- if .agent.secrets.fluxNinjaPlugin.create -}}
    {{- if .agent.secrets.fluxNinjaPlugin.value -}}
        {{ print .agent.secrets.fluxNinjaPlugin.value }}
    {{- else -}}
        {{- fail "Value of API Key for Agent cannot be empty when .Values.agent.secrets.fluxNinjaPlugin.create is set to true." -}}
    {{- end -}}
{{- else -}}
    {{ print "" }}
{{- end -}}
{{- end -}}
