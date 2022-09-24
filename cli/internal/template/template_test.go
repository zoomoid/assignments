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
	"testing"

	"github.com/lithammer/dedent"
	"github.com/zoomoid/assignments/v1/internal/config"
)

func TestGenerateDefaultAssignmentTemplate(t *testing.T) {
	bindings := TemplateBinding{
		ClassPath: "../assignments",
		Course:    "Example Course",
		Group:     "Group Alpha",
		Sheet:     "01",
		Due:       "April 20th, 2021",
		Includes: []config.Include{{
			Path: "./math.tex",
		}},
		Members: []config.GroupMember{
			{
				Name: "Max Mustermann",
				ID:   "123456",
			},
			{
				Name: "Erika Mustermann",
				ID:   "AB123456",
			},
		},
	}

	o, err := GenerateAssignmentTemplate(nil, &bindings)
	t.Logf("Rendered template to %s", o)

	if err != nil {
		t.Fatalf(`GenerateAssignmentTemplate() should NOT return an error with this binding, %v`, err)
	}

	if o.Len() == 0 {
		t.Fatal(`GenerateAssignmentTemplate() should NOT return an empty string on this binding`)
	}

}

func TestGenerateCustomAssignmentTemplate(t *testing.T) {
	bindings := TemplateBinding{
		ClassPath: "../assignments",
		Course:    "Example Course",
		Group:     "Group Alpha",
		Sheet:     "01",
		Due:       "April 20th, 2021",
		Members: []config.GroupMember{
			{
				Name: "Max Mustermann",
				ID:   "123456",
			},
			{
				Name: "Erika Mustermann",
				ID:   "AB123456",
			},
		},
	}

	template := dedent.Dedent(`
		\documentclass{article}

		\title{ {{- .Course -}} }
		{{- $authors := list "" | compact -}}
		{{- range $_, $m := .Members }}
		{{- $authors = append $authors $m.Name -}}
		{{- end }}
		\author{ {{- $authors | join " \\and " -}} }
		\date{ {{- .Due -}} }

		\begin{document}
		\maketitle

		\end{document}
	`)

	o, err := GenerateAssignmentTemplate(&template, &bindings)
	t.Logf("Rendered template to %s", o)

	if err != nil {
		t.Fatalf(`GenerateAssignmentTemplate() should NOT return an error with this binding, %v`, err)
	}

	if o.Len() == 0 {
		t.Fatal(`GenerateAssignmentTemplate() should NOT return an empty string on this binding`)
	}
}
