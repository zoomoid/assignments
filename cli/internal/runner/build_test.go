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

package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/zoomoid/assignments/v1/internal/config"
)

func TestBuildRunner(t *testing.T) {
	workingDirectory := t.TempDir()

	targetDirectory, err := makeSourceFile(workingDirectory)
	if err != nil {
		t.Error(err)
		return
	}

	ctx, err := makeAppContext(workingDirectory)
	if err != nil {
		t.Error(err)
		return
	}

	testRunner, err := New(ctx, &RunnerOptions{
		TargetDirectory:   targetDirectory,
		OverrideArtifacts: true,
		Quiet:             false,
	})
	if err != nil {
		t.Error(err)
	}

	t.Run("assignmentNumber", func(t *testing.T) {
		r := testRunner.Clone().Build()
		t.Run("TargetDirectory=assignment", func(t *testing.T) {
			r.SetTargetDirectory("assignment")
			i, err := r.assignmentNumber()
			if err != nil {
				// this is to be expected
				return
			}
			t.Error(fmt.Errorf("found assignment number where there shouldn't have, found %s, saw %s ", i, r.TargetDirectory()))
		})
		t.Run("TargetDirectory=assignment-03", func(t *testing.T) {
			r.SetTargetDirectory("assignment-03")
			s, err := r.assignmentNumber()
			if err != nil {
				t.Error(err)
				return
			}
			if s != "03" {
				t.Error(fmt.Errorf("Expected '%s', found '%s'", "03", s))
			}
		})
		t.Run("TargetDirectory=assignment-13", func(t *testing.T) {
			r.SetTargetDirectory("assignment-13")
			s, err := r.assignmentNumber()
			if err != nil {
				t.Error(err)
				return
			}
			if s != "13" {
				t.Error(fmt.Errorf("Expected '%s', found '%s'", "13", s))
			}
		})
		t.Run("TargetDirectory=assignment-3", func(t *testing.T) {
			r.SetTargetDirectory("assignment-3")
			s, err := r.assignmentNumber()
			if err != nil {
				t.Error(err)
				return
			}
			if s != "03" {
				t.Error(fmt.Errorf("Expected '%s', found '%s'", "03", s))
			}
		})
	})
	t.Run("MakeCommand", func(t *testing.T) {
		r := testRunner.Clone().Build()
		t.Run("predefined recipe", func(t *testing.T) {
			cmds, err := r.MakeCommand()
			if err != nil {
				t.Error(err)
				return
			}
			if len(cmds) != 1 {
				t.Error(fmt.Errorf("cmds is of length %d, expected %d ", len(cmds), 1))
				return
			}
			cmd := cmds[0]
			program := (*r.configuration.Spec.BuildOptions.BuildRecipe)[0].Command
			if cmd.Args[0] != program {
				t.Error(fmt.Errorf("cmd runs %s, expected %s", cmd.Args[0], program))
				return
			}
		})
		t.Run("empty recipe", func(t *testing.T) {
			r.configuration.Spec.BuildOptions.BuildRecipe = &config.Recipe{}
			cmds, err := r.MakeCommand()
			if err != nil {
				t.Error(err)
				return
			}
			// defaults to single latexmk run
			if len(cmds) != 1 {
				t.Error(fmt.Errorf("cmds is of length %d, expected %d ", len(cmds), 1))
				return
			}

			cmd := cmds[0]
			if cmd.Args[0] != DefaultBuildProgram {
				t.Error(fmt.Errorf("cmd runs %s, expected %s", cmd.Args[0], DefaultBuildProgram))
				return
			}
		})
		t.Run("pdflatex -> bibtex -> pdflatex recipe", func(t *testing.T) {
			r.configuration.Spec.BuildOptions.BuildRecipe = &config.Recipe{{
				Command: "pdflatex",
				Args:    []string{"-interaction=nonstopmode", "-file-line-error"},
			}, {
				Command: "bibtex",
			}, {
				Command: "pdflatex",
				Args:    []string{"-interaction=nonstopmode", "-file-line-error"},
			}}
			cmds, err := r.MakeCommand()
			if err != nil {
				t.Error(err)
				return
			}
			// defaults to single latexmk run
			if len(cmds) != len(*r.configuration.Spec.BuildOptions.BuildRecipe) {
				t.Error(fmt.Errorf("cmds is of length %d, expected %d ", len(cmds), len(*r.configuration.Spec.BuildOptions.BuildRecipe)))
				return
			}
			cmd := cmds[0]
			expected := (*r.configuration.Spec.BuildOptions.BuildRecipe)[0].Command
			if cmd.Args[0] != expected {
				t.Error(fmt.Errorf("cmd runs %s, expected %s", cmd.Args[0], expected))
				return
			}
		})
	})
	t.Run("makeArtifactsDirectory", func(t *testing.T) {
		r := testRunner.Clone().Build()
		t.Run("does not exist", func(t *testing.T) {
			newRoot := t.TempDir() // get a fresh temporary directory
			r.SetRoot(newRoot)
			r.SetArtifactsDirectory("dist")
			err := r.makeArtifactsDirectory()
			if err != nil {
				t.Error(err)
				return
			}
			// directory should exist now
			_, err = os.Stat(filepath.Join(newRoot, "dist"))
			if err != nil {
				t.Error(err)
				return
			}
			// Cleanup removes inner temp dir
		})
		t.Run("already exists", func(t *testing.T) {
			newRoot := t.TempDir()
			err := os.Mkdir(filepath.Join(newRoot, "dist"), 0777)
			if err != nil {
				t.Error(err)
				return
			}
			r.SetRoot(newRoot)
			// should return nil, as directory should already exist
			err = r.makeArtifactsDirectory()
			if err != nil {
				t.Error(err)
				return
			}
			// Cleanup removes inner temp dir
		})
	})
	// not actually building anything in between because that would
	// introduce a dependency on latex in testing - which is definitely
	// not what we want
	t.Run("exportArtifacts", func(t *testing.T) {
		r := testRunner.Clone().Build()
		// since we don't need an actual PDF at the src path,
		// just create a small test file named "assignment.pdf"
		newRoot := t.TempDir()
		// after switching root, create new source directory structure
		// to copy the artificial assignment.pdf to
		_, err := makeSourceFile(newRoot)
		if err != nil {
			t.Error(err)
			return
		}
		r.SetRoot(newRoot)
		pdfPath := filepath.Join(r.TargetDirectory(), "assignment.pdf")
		pdf, err := os.Create(pdfPath)
		if err != nil {
			t.Error(err)
			return
		}
		_, err = pdf.WriteString(validAssignmentTexCode)
		if err != nil {
			t.Error(err)
			return
		}
		err = pdf.Close()
		if err != nil {
			t.Error(err)
			return
		}

		dest, err := r.exportArtifacts()
		if err != nil {
			t.Error(err)
			return
		}

		sfi, _ := os.Stat(filepath.Join(r.TargetDirectory(), "assignment.pdf"))

		// check if assignment is in dist folder
		fi, err := os.Stat(dest)
		if err != nil {
			t.Error(err)
			return
		}

		if sfi.Size() != fi.Size() {
			t.Error(fmt.Errorf("src and dest files are not the same size, expected %d bytes, found %d bytes", sfi.Size(), fi.Size()))
		}
	})
}
