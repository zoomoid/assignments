package runner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/util"
)

type builder struct {
	*RunnerContext
}

type Builder interface {
	Runner
}

// MakeCommand implements the Runner spec in terms of transforming a given recipe into a
// slice of exec.Cmd, or using the default recipe
func (b *builder) MakeCommand() ([]*exec.Cmd, error) {
	recipe := b.configuration.Spec.BuildOptions.BuildRecipe

	if len(*recipe) == 0 {
		// use the default latexmk recipe
		recipe = &config.Recipe{
			{
				Command: DefaultBuildProgram,
				Args:    DefaultBuildArgs,
			},
		}
	}

	cmds, err := commandsFromRecipe(recipe, b.TargetDirectory(), b.Filename(), b.Quiet())
	return cmds, err
}

// Run implements the Runner specification for running a set of commands in terms of building
func (b *builder) Run() error {
	startTime := time.Now()
	log.Debug().Msgf("[runner/build] Started building %s", filepath.Join(b.TargetDirectory(), b.filename))

	cmds, err := b.MakeCommand()
	if err != nil {
		return err
	}

	b.Commands = cmds

	for i, cmd := range b.Commands {
		if cmd == nil {
			return fmt.Errorf("command %d is nil", i)
		}
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	log.Debug().Msgf("[runner/build] Finished building %s in %v", filepath.Join(b.TargetDirectory(), b.filename), time.Since(startTime))

	exportTime := time.Now()
	log.Debug().Msgf("[runner/export] Starting export of %s", filepath.Join(b.TargetDirectory(), b.filename))
	dest, err := b.exportArtifacts()
	if err != nil {
		return err
	}
	log.Debug().Msgf("[runner/export] Finished exporting from %s to %s in %v", filepath.Join(b.targetDirectory, b.filename), dest, time.Since(exportTime))
	return nil
}

// exportArtifacts copies the PDF from compilation to another directory for exporting artifacts collectively
func (b *builder) exportArtifacts() (string, error) {
	err := b.makeArtifactsDirectory()
	if err != nil {
		return "", err
	}

	d := b.ArtifactsDirectory()

	err = os.MkdirAll(d, 0777) // returns nil if d already exists
	if err != nil {
		return "", err
	}

	ai, err := b.assignmentNumber()
	if err != nil {
		return "", fmt.Errorf("failed to extract assignment number from target directory, got %s, %w", b.TargetDirectory(), err)
	}

	artifactsPdf := strings.Replace(b.Filename(), ".tex", ".pdf", 1)

	srcPath := filepath.Join(b.TargetDirectory(), artifactsPdf)
	destPath := filepath.Join(d, fmt.Sprintf("assignment-%s.pdf", ai))

	if _, err := os.Stat(destPath); !b.overrideArtifacts && err == nil {
		// file exists and the user did not specify --force flag,
		// exit before truncating the destination
		return destPath, errors.New("not overwriting existing file, add --force")
	}

	srcWriter, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}

	destWriter, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(destWriter, srcWriter)
	if err != nil {
		return "", err
	}

	log.Debug().Msgf("[runner/export] Copied artifact PDF from %s to %s", srcPath, destPath)

	pdfPath, err := filepath.Abs(destPath)
	if err != nil {
		// fallback to relative version
		pdfPath = destPath
	}

	return pdfPath, nil
}

// makeArtifactsDirectory ensures that the directory to copy artifact files to exists so the file
// descriptors can safely be created
func (b *builder) makeArtifactsDirectory() error {
	d := b.ArtifactsDirectory()
	if _, err := os.Stat(d); errors.Is(err, os.ErrNotExist) {
		// artifacts directory does not exist yet, try to create it
		err = os.MkdirAll(d, 0777)
		if err != nil {
			return err
		}
		log.Debug().Msgf("[runner/export] Created artifacts directory %s", d)
		return nil
	}
	log.Debug().Msgf("[runner/export] Artifacts directory %s already exists, skipping creation...", d)
	return nil
}

// assignmentNumber returns the assignment's number as a leading-zero string.
// Returns error when util.AssignmentNumberFromFilename returns an error
func (b *builder) assignmentNumber() (string, error) {
	s := filepath.Base(b.TargetDirectory())
	return util.AssignmentNumberFromRegex(util.AssignmentDirectoryPattern, s)
}
