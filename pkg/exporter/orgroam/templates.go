/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package orgroam

type OrgTemplate struct {
	TitleTemplate string
	EntryTemplate string
}

const commonOrgTitleTpl = `:PROPERTIES:
:ID:       {{ .UUID }}
:END:
#+title: {{ .Title }}
#+filetags: :{{ .Author }}:

`

var OrgTemplates = []OrgTemplate{
	{
		TitleTemplate: commonOrgTitleTpl,
		EntryTemplate: `
* {{ .Data }}
:PROPERTIES:
:ID:       {{ .UUID }}
:TYPE:     {{ .Type }}
{{- if ne .Location.Chapter "" }}
:CHAPTER:  {{ .Location.Chapter }}
{{- else }}
:CHAPTER:  {{ .Section }}
{{- end }}
{{- if gt .Location.Page 0 }}
:PAGE:     {{ .Location.Page }}
{{- end }}
{{- if gt .Location.Location 0 }}
:LOCATION: {{ .Location.Location }}
{{- end }}
:END:
`,
	},
	{
		TitleTemplate: commonOrgTitleTpl,
		EntryTemplate: `
* {{ .Data }}
:PROPERTIES:
:ID:       {{ .UUID }}
:END:
{{ .Type }} @
{{- if ne .Location.Chapter "" }}
Chapter: {{ .Location.Chapter }}
{{- else }}
Chapter: {{ .Section }}
{{- end }}
{{ .Location }}
`,
	},
}
