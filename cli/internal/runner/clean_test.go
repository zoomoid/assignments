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
