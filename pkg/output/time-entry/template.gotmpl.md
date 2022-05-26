ID: `{{ .ID }}`  
Billable: `{{ if .Billable }}yes{{ else }}no{{ end }}`  
Locked: `{{ if .IsLocked }}yes{{ else }}no{{ end }}`  
Project: {{ if eq .ProjectID "" -}}
  No Project
{{- else -}}
  {{ .Project.Name }} (`{{ .Project.ID }}`)
{{- end }}  
{{ with .Task -}}
Task: {{ .Name }} (`{{ .ID }}`)  
{{ end -}}
Interval: `{{ formatDateTime .TimeInterval.Start }}` until `{{ with .TimeInterval.End -}}
  {{ formatDateTime . }}
{{- else -}}
  now
{{- end }}`  
Description:
> {{ .Description }}
{{- with .Tags }}

Tags:
{{ range . }}
 * {{ .Name }} (`{{ .ID }}`)
{{- end }}
{{- end }}
{{- if not .Last }}

---
{{ end -}}
