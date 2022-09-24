/*
Copyright 2022 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package template

import (
	"bytes"
	"strings"
	"text/template"

	_ "embed"

	"github.com/Masterminds/sprig/v3"
	"github.com/lithammer/dedent"
	"github.com/zoomoid/assignments/v1/internal/config"
)

var (
	DefaultSheetTemplate = strings.TrimPrefix(dedent.Dedent(`
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
		\member[{{- $member.ID -}}]{ {{- $member.Name -}} }
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
		tpl = &DefaultSheetTemplate
	}
	tmpl := template.Must(template.New("assignment").Funcs(sprig.TxtFuncMap()).Parse(*tpl))

	var output bytes.Buffer

	err := tmpl.Execute(&output, bindings)

	if err != nil {
		return nil, err
	}
	return &output, nil
}
