
{{- define "helm.release.name" -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
