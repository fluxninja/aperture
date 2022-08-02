{{/*
Return a remote image path based on `.Values` (passed as root) and `.` (any `.image` from `.Values` passed as parameter)
*/}}
{{- define "image-path" -}}
{{- if .image.repository -}}
{{- .image.repository -}}:{{ .image.tag }}
{{- else -}}
{{ .root.registry }}/{{ .image.name }}:{{ .image.tag }}
{{- end -}}
{{- end -}}

{{/*
Return the proper agent image name
*/}}
{{- define "aperture-agent.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.agent.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper agent image name
*/}}
{{- define "aperture-agent-controller.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.agentController.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "aperture-agent.imagePullSecrets" -}}
{{- include "common.images.pullSecrets" (dict "images" (list .Values.agent.image .Values.agentController.image) "global" .Values.global) -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "aperture-agent.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (printf "%s" (include "common.names.fullname" .)) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the agent apiKey secret to use
*/}}
{{- define "aperture-agent.apiKeySecretRefName" -}}
{{- if .Values.fluxninjaPlugin.apiKeySecret.agent.secretKeyRef.name -}}
    {{- .Values.fluxninjaPlugin.apiKeySecret.agent.secretKeyRef.name -}}
{{- else -}}
    {{ include "common.names.fullname" . }}-agent-apikey
{{- end -}}
{{- end -}}

{{/*
Create the name of the controller apiKey secret to use
*/}}
{{- define "aperture-controller.apiKeySecretRefName" -}}
{{- if .Values.fluxninjaPlugin.apiKeySecret.controller.secretKeyRef.name -}}
    {{- .Values.fluxninjaPlugin.apiKeySecret.controller.secretKeyRef.name -}}
{{- else -}}
    {{ include "common.names.fullname" . }}-controller-apikey
{{- end -}}
{{- end -}}

{{/*
Compile all warnings into a single message.
*/}}
{{- define "aperture-agent.validateValues" -}}
{{- $messages := list -}}
{{- $messages := append $messages "placeholder" -}}
{{- $messages := without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{-   printf "\nVALUES VALIDATION:\n%s" $message -}}
{{- end -}}
{{- end -}}

{{/*
Create the endpoint of the etcd for Aperture Agent and Controller
{{ include "aperture.etcd.endpoints" ( dict "etcd" .Values.path.to.the.etcd "context" $.context $) }}
*/}}
{{- define "aperture.etcd.endpoints" -}}
{{- $endpoints := list -}}
{{- if .etcd.enabled -}}
    {{- $endpoints = append $endpoints (printf "http://%s-etcd:2379" .context.Release.Name) -}}
{{- else -}}
    {{ $endpoints = without .etcd.endpoints "" }}
    {{- if empty $endpoints -}}
        {{- fail "Value for etcd endpoints of Agent or Controller cannot be empty when .Values.etcd.enbled is set to false." -}}
    {{- end -}}
{{- end -}}
{{ print $endpoints }}
{{- end -}}

{{/*
Create the address of the Prometheus for Aperture Agent and Controller
{{ include "aperture.prometheus.address" ( dict "prometheus" .Values.path.to.the.prometheus "context" $.context $) }}
*/}}
{{- define "aperture.prometheus.address" -}}
{{- if .prometheus.enabled -}}
    {{- printf "http://%s-prometheus-server:80" .context.Release.Name -}}
{{- else -}}
    {{- if .prometheus.address -}}
        {{ print .prometheus.address }}
    {{- else -}}
        {{- fail "Value for prometheus address of Agent or Controller cannot be empty when .Values.prometheus.enbled is set to false." -}}
    {{- end -}}
{{- end -}}
{{- end -}}
