{{/*
memcached Service
*/}}
{{- define "mimir.memcached.service" -}}
{{ with (index $.ctx.Values $.component) }}
{{- if .enabled -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ include "mimir.resourceName" (dict "ctx" $.ctx "component" $.component) }}-headless
  labels:
    {{- include "mimir.labels" (dict "ctx" $.ctx "component" $.component) | nindent 4 }}
    prometheus.io/service-monitor: "false"
    {{- with .service.labels }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
  annotations:
    {{- toYaml .service.annotations | nindent 4 }}
spec:
  type: ClusterIP
  clusterIP: None
  ports:
    - name: memcached-client
      port: {{ .port }}
      targetPort: {{ .port }}
    {{ if .metrics.enabled }}
    - name: exporter-http-metrics
      port: {{ .metrics.containerPort }}
      targetPort: {{ .metrics.containerPort }}
    {{ end }}
  selector:
    {{- include "mimir.selectorLabels" (dict "ctx" $.ctx "component" $.component) | nindent 4 }}
{{- end -}}
{{- end -}}
{{- end -}}
