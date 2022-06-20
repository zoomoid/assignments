package latexmk

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zoomoid/assignments/v1/internal/context"
)

// RunnerOptions struct to carry configuration for the latexmk runs
type RunnerOptions struct {
	// Directory to run latexmk from
	TargetDirectory string
	// TeX source file to compile, defaults to "assignment.tex"
	Filename string
	// Number of latexmk runs to make when running RunBuildMultiple
	Runs int
	// ArtifactsDirectory to copy to, defaults to "$PWD/dist"
	ArtifactsDirectory string
	// Quiet makes the latexmk run capture the output inside a buffer instead of piping to stdout
	Quiet bool
}

type RunnerContext struct {
	context.AppContext
	RunnerOptions
}

var (
	defaultProgram        = "latexmk"
	defaultLatexmkOptions = []string{"-pdf", "-interaction=nonstopmode", "-file-line-error", "-shell-escape"}
)

// New creates a new runner context from the given parameters and applies sensible defaults
func New(context *context.AppContext, options *RunnerOptions) (*RunnerContext, error) {
	runner := &RunnerContext{
		AppContext: *context,
	}

	if options.Runs <= 0 {
		runner.Runs = 3
	} else {
		runner.Runs = options.Runs
	}

	if options.TargetDirectory == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, errors.New("failed to determine current working directory for runner context creation")
		}
		runner.TargetDirectory = pwd
	} else {
		runner.TargetDirectory = options.TargetDirectory
	}

	if options.Filename == "" {
		runner.Filename = options.Filename
	} else {
		runner.Filename = "assignment.tex"
	}

	return runner, nil
}

// NewMust creates a new runner context or exits with error if creation fails
func NewMust(context *context.AppContext, options *RunnerOptions) *RunnerContext {
	r, err := New(context, options)
	if err != nil {
		context.Logger.Fatalf("Failed to create runner context, %v", err)
	}
	return r
}

// RunClean runs latexmk with cleaning options in the target directory to clean up after building
func (r *RunnerContext) RunClean() error {
	out := &bytes.Buffer{}

	r.Logger.Infof("Cleaning up using latexmk", "pwd", r.Cwd)

	cmd := exec.Command(defaultProgram, "-C")

	if r.Quiet {
		cmd.Stdout = out
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Dir = r.TargetDirectory

	if err := cmd.Run(); err != nil {
		r.Logger.Debugf("Cleaning up failed with %v", err)
		return err
	}
	return nil
}

// combineBuffers takes a slice of byte buffers and combines those elements into a single byte buffer
func combineBuffers(buffers *[]bytes.Buffer) (*bytes.Buffer, error) {
	o := &bytes.Buffer{}
	for _, b := range *buffers {
		_, err := o.Write(b.Bytes())
		if err != nil {
			return nil, err
		}
	}
	return o, nil
}

// RunBuildOnce runs latexmk once and returns the output from the run
func (r *RunnerContext) RunBuildOnce() (*bytes.Buffer, error) {
	r.Logger.Info("[1/1] Running compiler", "pwd", r.TargetDirectory)
	output, err := r.runBuild()
	return output, err
}

// RunBuildMultiple runs latexmk a specified number of times and returns the combined output of all runs
func (r *RunnerContext) RunBuildMultiple() (*bytes.Buffer, error) {
	outs := []bytes.Buffer{}

	var err error
	for i := 1; i <= r.Runs; i++ {
		r.Logger.Infof("[%d/%d] Running compiler", i, r.Runs, "pwd", r.TargetDirectory)
		b, err := r.runBuild()
		if err != nil {
			// break on first error occurence
			break
		}

		outs = append(outs, *b)
	}

	combinedBuffer, e := combineBuffers(&outs)
	if e != nil {
		r.Logger.Debugf("Failed to concatenate output buffers, %v", e)
		return nil, e
	}

	return combinedBuffer, err
}

// runBuild parses latexmk arguments, runs the command, and captures the output
func (r *RunnerContext) runBuild() (*bytes.Buffer, error) {
	out := &bytes.Buffer{}

	recipe := r.Configuration.Spec.BuildOptions.Recipe

	program := defaultProgram
	args := defaultLatexmkOptions
	if recipe != nil && recipe.Command != "" {
		program = recipe.Command
	}
	if recipe != nil && len(recipe.Args) > 0 {
		args = recipe.Args
	}

	args = append(args, "-f", r.Filename)

	cmd := exec.Command(program, args...)

	if r.Quiet {
		cmd.Stdout = out
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Dir = r.TargetDirectory

	if err := cmd.Run(); err != nil {
		r.Logger.Debugf("Failed to build with %v", err)
		return nil, err
	}

	return out, nil
}

// ExportArtifacts copies the PDF from compilation to another directory for exporting artifacts collectively
func (r *RunnerContext) ExportArtifacts() error {
	if r.ArtifactsDirectory == "" {
		r.ArtifactsDirectory = filepath.Join(r.Cwd, "dist")
	}

	if _, err := os.Stat(r.ArtifactsDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(r.ArtifactsDirectory, os.ModeDir)
		if err != nil {
			return err
		}
		r.Logger.Info("Created artifacts directory", "directory", r.ArtifactsDirectory)
	}

	base := filepath.Base(r.TargetDirectory)

	d := filepath.Join(r.ArtifactsDirectory, base)

	err := os.MkdirAll(d, os.ModeDir)

	if err != nil {
		return err
	}

	src := filepath.Join(r.TargetDirectory, r.Filename)

	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	dest := filepath.Join(d, "assignment.pdf")

	err = ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		return err
	}

	r.Logger.Info("Copied artifact PDF", "src", src, "dest", dest)

	return nil
}
