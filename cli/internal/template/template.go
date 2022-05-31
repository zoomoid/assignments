package template

import (
	"bytes"
	"text/template"

	_ "embed"

	"github.com/Masterminds/sprig/v3"
	config "github.com/zoomoid/assignments/v1/internal/config"
)

//go:embed sheet.tmpl
var sheetTemplate string

type TemplateBinding struct {
	ClassPath string
	Course    string
	Group     string
	Sheet     string
	Due       string
	Members   []config.GroupMember
}

func GenerateAssignmentTemplate(bindings *TemplateBinding) (string, error) {
	tmpl := template.Must(template.New("assignment").Funcs(sprig.TxtFuncMap()).Parse(sheetTemplate))

	var output bytes.Buffer

	err := tmpl.Execute(&output, bindings)

	if err != nil {
		return "", err
	}
	return output.String(), nil
}
