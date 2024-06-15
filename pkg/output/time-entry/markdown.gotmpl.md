{{- $project := "" -}}
{{- if eq .Project nil }}
  {{- $project = "No Project" -}}
{{- else -}}
  {{- $project = concat "**" .Project.Name "**" -}}
  {{- if ne .Task nil -}}
    {{- $project = concat $project ": " .Task.Name  -}}
  {{- else if ne .Project.ClientName "" -}}
    {{- $project = concat $project " - " .Project.ClientName  -}}
  {{- end -}}
{{- end -}}

{{- $bil := "No" -}}
{{- if .Billable -}}{{ $bil = "Yes" }}{{- end -}}

{{- $tags := "" -}}
{{- with .Tags -}}
  {{- range $index, $element := . -}}
    {{- if ne $index 0 }}{{ $tags = concat $tags ", " }}{{ end -}}
    {{- $tags = concat $tags $element.Name -}}
  {{- end -}}
{{- else -}}
  {{- $tags = "No Tags" -}}
{{- end -}}

{{- $pad := maxLength .Description $project $tags $bil -}}

## _Time Entry_: {{ .ID }}

_Time and date_  
**{{ dsf .TimeInterval.Duration }}** | {{ if eq .TimeInterval.End nil -}}
Start Time: _{{ formatTimeWS .TimeInterval.Start }}_ 🗓 Today
{{- else -}}
{{ formatTimeWS .TimeInterval.Start }} - {{ formatTimeWS .TimeInterval.End }} 🗓
{{- .TimeInterval.Start.Format " 01/02/2006" }}
{{- end }}

|               | {{ pad "" $pad }} |
|---------------|-{{ repeatString "-" $pad }}-|
| _Description_ | {{ pad .Description $pad }} |
| _Project_     | {{ pad $project $pad }} |
| _Tags_        | {{ pad $tags $pad }} |
| _Billable_    | {{ pad $bil $pad }} |
