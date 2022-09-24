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
	"testing"
)

func TestCleanRunner(t *testing.T) {
	workingDirectory := t.TempDir()

	targetDirectory, err := makeSourceFile(workingDirectory)
	if err != nil {
		t.Error(err)
	}

	ctx, err := makeAppContext(workingDirectory)
	if err != nil {
		t.Error(err)
	}

	testRunner, err := New(ctx, &RunnerOptions{
		TargetDirectory:   targetDirectory,
		OverrideArtifacts: true,
		Quiet:             false,
	})
	if err != nil {
		t.Error(err)
	}
	t.Run("MakeCommand", func(t *testing.T) {
		r := testRunner.Clone().Clean()
		t.Run("latexmk -C", func(t *testing.T) {
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
			if len(cmd.Args) != 2 {
				t.Error(fmt.Errorf("expected command to be of length %d, found %d", 2, len(cmd.Args)))
			}
			program := cmd.Args[0]
			if program != DefaultBuildProgram {
				t.Error(fmt.Errorf("expected program to be %s, found %s", DefaultBuildProgram, program))
			}
			arg := cmd.Args[1]
			if arg != "-C" {
				t.Error(fmt.Errorf("expected arguments to be singleton %s, found %s", "-C", arg))
			}
		})
	})
}
