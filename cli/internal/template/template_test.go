package template

import (
	"testing"

	"github.com/zoomoid/assignments/v1/internal/config"
)

func TestGenerateAssignmentTemplate(t *testing.T) {
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

	o, err := GenerateAssignmentTemplate(&bindings)
	t.Logf("Rendered template to %s", o)

	if err != nil {
		t.Fatalf(`GenerateAssignmentTemplate() should NOT return an error with this binding, %v`, err)
	}

	if len(o) == 0 {
		t.Fatal(`GenerateAssignmentTemplate() should NOT return an empty string on this binding`)
	}

}
