{{/*
Expand the name of the chart.
*/}}
{{- define "waldritter-website.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "waldritter-website.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "waldritter-website.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "waldritter-website.labels" -}}
helm.sh/chart: {{ include "waldritter-website.chart" . }}
{{ include "waldritter-website.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "waldritter-website.selectorLabels" -}}
app.kubernetes.io/name: {{ include "waldritter-website.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "waldritter-website.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "waldritter-website.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the image for a component
*/}}
{{- define "waldritter-website.image" -}}
{{- $imageRoot := .imageRoot -}}
{{- $global := .global -}}
{{- $registry := $global.imageRegistry | default $imageRoot.repository -}}
{{- printf "%s:%s" $registry ($imageRoot.tag | default "latest") -}}
{{- end }}

{{/*
Content API labels
*/}}
{{- define "waldritter-website.contentApi.labels" -}}
{{ include "waldritter-website.labels" . }}
app.kubernetes.io/component: content-api
{{- end }}

{{/*
Rails API labels
*/}}
{{- define "waldritter-website.railsApi.labels" -}}
{{ include "waldritter-website.labels" . }}
app.kubernetes.io/component: rails-api
{{- end }}

{{/*
Website UI labels
*/}}
{{- define "waldritter-website.websiteUi.labels" -}}
{{ include "waldritter-website.labels" . }}
app.kubernetes.io/component: website-ui
{{- end }}

{{/*
Admin UI labels
*/}}
{{- define "waldritter-website.adminUi.labels" -}}
{{ include "waldritter-website.labels" . }}
app.kubernetes.io/component: admin-ui
{{- end }}