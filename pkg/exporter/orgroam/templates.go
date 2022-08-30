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
{{- if eq .Type "NOTE" }}
-- "{{ .UserNote  }}"
{{- end }}
:PROPERTIES:
:ID:       {{ .UUID }}
:TYPE:     {{ .Type }}
:CHAPTER:  {{ .Location.Chapter }}
{{- if .Location.Page }}
:PAGE:     {{ .Location.Page }}
{{- end }}
{{- if .Location.Location }}
:LOCATION: {{ .Location.Location }}
{{- end }}
:END:
`,
	},
	{
		TitleTemplate: commonOrgTitleTpl,
		EntryTemplate: `
* {{ .Data }}
{{- if eq .Type "NOTE" }}
-- "{{ .UserNote  }}"
{{- end }}
{{- if ne .UserNote "" }}
-- "{{ .UserNote  }}"
{{- end }}
:PROPERTIES:
:ID:       {{ .UUID }}
:END:
{{ .Type }} @
Chapter: {{ .Location.Chapter }}
{{ .Location }}
`,
	},
}
