package config

import (
	"os"
	"testing"
)

const configFile string = ".assignments.yaml"

func TestReadConfigMap(t *testing.T) {
	config := `
spec:
  course: "Linear Algebra I"
  group: "Group Alpha"
  members:
    - id: "123456"
      name: "Max Mustermann"
    - id: "AB123456"
      name: "Erika Mustermann"
    - id: "69420"
      name: "Kim Took"
  includes:
    - path: ../include.tex
    - path: ../packages.tex
status:
  assignment: 1`
	os.WriteFile(configFile, []byte(config), 0644)
	defer os.Remove(configFile)
	cfgFile, err := ReadConfigMap("")

	if err != nil {
		t.Fatalf("Failed to read in freshly created config file, %v", err)
	}

	if cfgFile == nil {
		t.Fatal("Returned Configuration struct is nil")
	}
}

func TestWriteConfigMap(t *testing.T) {
	config := Configuration{
		Spec: &ConfigurationSpec{
			Course: "Linear Algebra I",
			Group:  "Group Alpha",
			Members: []GroupMember{{
				Name: "Max Mustermann",
				ID:   "12346",
			}, {
				Name: "Erika Mustermann",
				ID:   "AB123456",
			}, {
				Name: "Kim Took",
				ID:   "69420",
			}},
			Includes: []Include{{
				Path: "../includes.tex",
			}, {
				Path: "../packages.tex",
			}, {
				Path: "../custom-macros.tex",
			}},
		},
		Status: &ConfigurationStatus{
			Assignment: 1,
		},
	}
	defer os.Remove(configFile)
	err := WriteConfigMap(&config, "")
	if err != nil {
		t.Fatalf("Failed to write configmap to file, %v", err)
	}

}

func TestWriteThenRead(t *testing.T) {
	config := Configuration{
		Spec: &ConfigurationSpec{
			Course: "Linear Algebra I",
			Group:  "Group Alpha",
			Members: []GroupMember{{
				Name: "Max Mustermann",
				ID:   "12346",
			}, {
				Name: "Erika Mustermann",
				ID:   "AB123456",
			}, {
				Name: "Kim Took",
				ID:   "69420",
			}},
			Includes: []Include{{
				Path: "../includes.tex",
			}, {
				Path: "../packages.tex",
			}, {
				Path: "../custom-macros.tex",
			}},
		},
		Status: &ConfigurationStatus{
			Assignment: 1,
		},
	}
	defer os.Remove(configFile)
	err := WriteConfigMap(&config, "")
	if err != nil {
		t.Fatalf("Failed to write configmap to file, %v", err)
	}

	readConfig, err := ReadConfigMap("")
	if err != nil {
		t.Fatalf("Failed to read configmap from file, %v", err)
	}

	if config.Status.Assignment != readConfig.Status.Assignment {
		t.Fatalf("Failed to read config back in, differing values, expected .Status.Assignment to be %d, got %d", config.Status.Assignment, readConfig.Status.Assignment)
	}
}
