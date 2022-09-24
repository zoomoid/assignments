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

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lithammer/dedent"
)

const configFile string = ".assignments.yaml"

func TestRead(t *testing.T) {
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
  generate:
    create:
      - figures
  includes:
    - path: ../include.tex
    - path: ../packages.tex
status:
  assignment: 1`
	os.WriteFile(configFile, []byte(config), 0644)
	defer os.Remove(configFile)
	cfgFile, err := Read("")

	if err != nil {
		t.Fatalf("Failed to read in freshly created config file, %v", err)
	}

	if cfgFile == nil {
		t.Fatal("Returned Configuration struct is nil")
	}
}

func TestWrite(t *testing.T) {
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
	err := Write(&config, "")
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
	err := Write(&config, "")
	if err != nil {
		t.Fatalf("Failed to write configmap to file, %v", err)
	}

	readConfig, err := Read("")
	if err != nil {
		t.Fatalf("Failed to read configmap from file, %v", err)
	}

	if config.Status.Assignment != readConfig.Status.Assignment {
		t.Fatalf("Failed to read config back in, differing values, expected .Status.Assignment to be %d, got %d", config.Status.Assignment, readConfig.Status.Assignment)
	}
}

func TestFind(t *testing.T) {
	t.Run("depth=0", func(t *testing.T) {
		dir := t.TempDir()
		f, err := os.Create(filepath.Join(dir, ConfigurationFileName))
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		path, err := Find(dir)
		if err != nil {
			t.Fatal(err)
		}
		if path != dir {
			t.Error(fmt.Errorf("Expected path to be %s, found %s", path, dir))
		}
	})
	t.Run("depth=1", func(t *testing.T) {
		dir := t.TempDir()
		subDir := filepath.Join(dir, "sub1")
		err := os.MkdirAll(subDir, 0777)
		if err != nil {
			t.Fatal(err)
		}
		f, err := os.Create(filepath.Join(dir, ConfigurationFileName))
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		path, err := Find(subDir)
		if err != nil {
			t.Fatal(err)
		}
		if path != dir {
			t.Error(fmt.Errorf("Expected path to be %s, found %s", dir, path))
		}
	})
	t.Run("depth=5", func(t *testing.T) {
		dir := t.TempDir()
		subDir := filepath.Join(dir, "sub1", "sub2", "sub3", "sub4", "sub5")
		err := os.MkdirAll(subDir, 0777)
		if err != nil {
			t.Fatal(err)
		}
		f, err := os.Create(filepath.Join(dir, ConfigurationFileName))
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		path, err := Find(subDir)
		if err != nil {
			t.Fatal(err)
		}
		if path != dir {
			t.Error(fmt.Errorf("Expected path to be %s, found %s", dir, path))
		}
	})
	t.Run("no configmap", func(t *testing.T) {
		dir := t.TempDir()
		subDir := filepath.Join(dir, "sub1", "sub2", "sub3", "sub4", "sub5")
		err := os.MkdirAll(subDir, 0777)
		if err != nil {
			t.Fatal(err)
		}
		path, err := Find(subDir)
		if !errors.Is(err, ErrNoConfigmap) {
			t.Fatal(err)
		}
		if path != "" {
			t.Error(fmt.Errorf("Expected path to be '%s', found '%s'", "", path))
		}
	})
}

func TestMarshal(t *testing.T) {
	config := Configuration{
		Spec: &ConfigurationSpec{
			Course: "Test Course",
			Group:  "Test Group",
		},
		Status: &ConfigurationStatus{
			Assignment: 1,
		},
	}
	out, err := Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Error("Expected byte array to not be length 0")
	}
}

func TestUnmarshal(t *testing.T) {
	marshalledConfig := dedent.Dedent(`
spec:
  course: Test Course
  group: Test Group
  build:
    recipe:
      - command: latexmk
        args:
          - "-pdf"
          - "-file-line-error"
          - "-shell-escape"
          - "-interaction=nonstopmode"
  bundle:
    template: "sheet_{{._id}}.zip"
    include:
      - "code/*"
status:
  assignment: 1
	`)
	out := Configuration{}
	err := Unmarshal([]byte(marshalledConfig), &out)
	if err != nil {
		t.Fatal(err)
	}
	if out.Spec.Course != "Test Course" {
		t.Fatal(fmt.Errorf("expected %s, found %s", "Test Course", out.Spec.Course))
	}
	if out.Spec.Group != "Test Group" {
		t.Fatal(fmt.Errorf("expected %s, found %s", "Test Group", out.Spec.Group))
	}
	if out.Status.Assignment != 1 {
		t.Fatal(fmt.Errorf("expected %d, found %d", 1, out.Status.Assignment))
	}
}
