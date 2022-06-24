package latexmk

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/util"
)

type builder struct {
	*RunnerContext
}

var _ Runner = &builder{}

func (b *builder) makeBuildCommands() error {
	recipe := b.Configuration.Spec.BuildOptions.Recipe

	if recipe == nil {
		// use the default latexmk recipe
		recipe = []config.Recipe{
			{
				Command: defaultProgram,
				Args:    defaultLatexmkOptions,
			},
		}
	}

	cmds := []*exec.Cmd{}
	for i, tool := range recipe {
		if tool.Command == "" {
			return fmt.Errorf("failed to make build commands, missing program in recipe step %d", i)
		}
		program := tool.Command
		args := []string{}
		if len(tool.Args) > 0 {
			args = tool.Args
		}

		args = append(args, filepath.ToSlash(b.Filename))
		cmd := exec.Command(program, args...)
		out := &bytes.Buffer{}

		if b.Quiet {
			sink := bufio.NewWriter(out)
			cmd.Stdout = sink
		} else {
			cmd.Stdout = os.Stdout
		}
		cmd.Dir = b.TargetDirectory

		cmds = append(cmds, cmd)
	}

	b.Commands = cmds
	return nil
}

func (b *builder) Run() error {
	startTime := time.Now()
	b.Logger.Debug("[runner/build] Starting builder", "pwd", b.Root, "file", b.Filename)
	err := b.makeBuildCommands()
	if err != nil {
		return err
	}
	for _, cmd := range b.Commands {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	b.Logger.Debug("[runner/build] Finished building stage", "pwd", b.Root, "file", b.Filename, "duration", time.Since(startTime))

	exportTime := time.Now()
	b.Logger.Debug("[runner/export] Starting export", "pwd", b.Root, "file", b.Filename)
	dest, err := b.exportArtifacts()
	if err != nil {
		return err
	}
	b.Logger.Debug("[runner/export] Finished exporting stage", "pwd", b.Root, "src_file", b.Filename, "dest_file", dest, "duration", time.Since(exportTime))
	return nil
}

// exportArtifacts copies the PDF from compilation to another directory for exporting artifacts collectively
func (b *builder) exportArtifacts() (string, error) {
	if b.ArtifactsDirectory == "" {
		b.ArtifactsDirectory = filepath.Join(b.Root, "dist")
	}

	err := b.makeArtifactsDirectory()
	if err != nil {
		return "", err
	}

	d := filepath.Join(b.ArtifactsDirectory)

	err = os.MkdirAll(d, os.ModeDir) // returns nil if d already exists
	if err != nil {
		return "", err
	}

	ai, err := b.assignmentNumber()
	if err != nil {
		return "", fmt.Errorf("failed to derive assignment number from target directory, got %s", b.TargetDirectory)
	}

	srcPath := filepath.Join(b.TargetDirectory, b.Filename)
	destPath := filepath.Join(d, fmt.Sprintf("assignment-%s.pdf", ai))

	if _, err := os.Stat(destPath); !b.OverrideArtifacts && err == nil {
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
	_, err = io.Copy(srcWriter, destWriter)
	if err != nil {
		return "", err
	}

	b.Logger.Debugf("[runner/export] Copied artifact PDF from %s to %s", srcPath, destPath)

	pdfPath, err := filepath.Abs(destPath)
	if err != nil {
		// fallback to relative version
		pdfPath = destPath
	}

	return pdfPath, nil
}

func (b *builder) makeArtifactsDirectory() error {
	if _, err := os.Stat(b.ArtifactsDirectory); errors.Is(err, os.ErrNotExist) {
		// artifacts directory does not exist yet, try to create it
		err = os.MkdirAll(b.ArtifactsDirectory, os.ModeDir)
		if err != nil {
			return err
		}
		b.Logger.Debug("[runner/export] Created artifacts directory", "directory", b.ArtifactsDirectory)
	}
	b.Logger.Debug("[runner/export] Artifacts directory already exists, skipping creation...", "directory", b.ArtifactsDirectory)
	return nil
}

func (b *builder) assignmentNumber() (string, error) {
	s := filepath.Base(b.TargetDirectory)
	s = strings.ReplaceAll(s, "assignment-", "")
	s = strings.ReplaceAll(s, ".pdf", "")
	i, err := strconv.Atoi(s)
	if err != nil {
		return "", err
	}
	return util.AddLeadingZero(uint32(i)), nil
}
