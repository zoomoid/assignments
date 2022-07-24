package template

import (
	"bytes"
	"strings"
	"text/template"

	_ "embed"

	"github.com/Masterminds/sprig/v3"
	"github.com/lithammer/dedent"
	config "github.com/zoomoid/assignments/v1/internal/config"
)

var (
	defaultSheetTemplate = strings.TrimPrefix(dedent.Dedent(`
		\documentclass{csassignments}
		{{- range $_, $input := .Includes -}}
		\input{ {{- $input -}} }
		{{ end }}
		\course{ {{- .Course -}} }
		\group{ {{- .Group | default "" -}} }
		\sheet{ {{- .Sheet | default "" -}} }
		\due{ {{- .Due | default "" -}} }
		{{- range $_, $member := .Members }}
		{{- $firstname := ($member.Name | splitList " " | initial | join " ") | default "" }}
		{{- $lastname := ($member.Name | splitList " " | last) | default "" }}
		\member{ {{- $firstname -}} }{ {{- $lastname -}} }{ {{- $member.ID -}} }
		{{- end }}
		
		\begin{document}
		\maketitle
		\gradingtable
		
		% Start the assignment here
		
		\end{document}
	`), "\n")
)

type TemplateBinding struct {
	ClassPath string
	Course    string
	Group     string
	Sheet     string
	Due       string
	Members   []config.GroupMember
	Includes  []config.Include
}

func GenerateAssignmentTemplate(tpl *string, bindings *TemplateBinding) (*bytes.Buffer, error) {
	if tpl == nil || *tpl == "" {
		tpl = &defaultSheetTemplate
	}
	tmpl := template.Must(template.New("assignment").Funcs(sprig.TxtFuncMap()).Parse(*tpl))

	var output bytes.Buffer

	err := tmpl.Execute(&output, bindings)

	if err != nil {
		return nil, err
	}
	return &output, nil
}
