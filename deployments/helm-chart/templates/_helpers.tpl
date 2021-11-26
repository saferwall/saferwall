{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart.  This is suffixed with -saferwall, which means subtract 9 from longest 63 available */}}
*/}}
{{- define "saferwall.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 54 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "saferwall.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "saferwall.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels.
*/}}
{{- define "saferwall.labels" -}}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ include "saferwall.chart" . }}
{{- end -}}


{{/*
Create the name of the service account to use
*/}}
{{- define "saferwall.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "saferwall.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}


{{/*
Create the name of the hostnames
*/}}
{{- define "saferwall.ui-hostname" -}}
{{ .Values.global.hostname }}
{{- end -}}
{{- define "saferwall.webapis-hostname" -}}
{{ printf "api.%s" .Values.global.hostname }}
{{- end -}}
{{- define "saferwall.website-hostname" -}}
{{ printf "about.%s" .Values.global.hostname }}
{{- end -}}


{{/*
Create the docker private registry server token.
*/}}
{{- define "saferwall.private-registry.token" -}}
{{- if .Values.global.privateRegistryServer.enabled -}}
{{ printf "{\"auths\":{\"https://index.docker.io/v1/\":{\"auth\":\"%s\"}}}" .Values.global.privateRegistryServer.token }}
{{- end -}}
{{- end -}}