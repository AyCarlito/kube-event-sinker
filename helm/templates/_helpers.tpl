{{/*
Include a nodeSelector
*/}}
{{- define "kube-event-sinker.nodeSelector" -}}
  {{- if .Values.nodeSelector }}
nodeSelector: {{ toYaml .Values.nodeSelector | nindent 2 }}
  {{- end }}
{{- end -}}

{{/*
Include affinity
*/}}
{{- define "kube-event-sinker.affinity" -}}
  {{- if .Values.affinity }}
affinity: {{ toYaml .Values.affinity | nindent 2 -}}
  {{- end }}
{{- end -}}

{{/*
Include tolerations
*/}}
{{- define "kube-event-sinker.tolerations" -}}
  {{- if .Values.tolerations }}
tolerations:
{{ toYaml .Values.tolerations }}
  {{- end }}
{{- end -}}

{{/*
Include placement instructions: affinity, tolerations and nodeSelector
*/}}
{{- define "kube-event-sinker.placement" -}}
{{ include "kube-event-sinker.affinity" . }}
{{ include "kube-event-sinker.tolerations" . }}
{{ include "kube-event-sinker.nodeSelector" . }}
{{- end -}}


{{/*
Include resources
*/}}
{{- define "kube-event-sinker.resources" -}}
  {{- if .Values.resources }}
resources:
{{ toYaml .Values.resources }}
  {{- end }}
{{- end -}}

{{/*
Include image pull secrets.
*/}}
{{- define "kube-event-sinker.imagePullSecrets" -}}
  {{- if .Values.image.pullSecrets }}
imagePullSecrets:
{{ toYaml .Values.image.pullSecrets }}
  {{- end }}
{{- end -}}
