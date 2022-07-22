/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package parser

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
:END:
{{ .Type }} @
{{- if eq .Location.Chapter "" }}
Chapter: {{ .Section }}
{{ else }}
Section: {{ .Section }}
{{ end -}}
{{ .Location }}
`,
	},
	{
		TitleTemplate: commonOrgTitleTpl,
		EntryTemplate: `
* {{ .Data }}
:PROPERTIES:
:ID:       {{ .UUID }}
:TYPE:     {{ .Type }}
{{- if eq .Location.Chapter "" }}
:CHAPTER:  {{ .Section }}
{{ else }}
:SECTION:  {{ .Section }}
{{ end -}}
:LOCATION: {{ .Location }}
:END:
`,
	},
}
