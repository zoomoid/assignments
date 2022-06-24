package latexmk

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"

	"github.com/zoomoid/assignments/v1/internal/context"
)

// RunnerOptions struct to carry configuration for the latexmk runs
type RunnerOptions struct {
	// Directory to run latexmk from
	TargetDirectory string
	// TeX source file to compile, defaults to "assignment.tex"
	Filename string
	// ArtifactsDirectory to copy to, defaults to "$PWD/dist"
	ArtifactsDirectory string
	// Quiet makes the latexmk run capture the output inside a buffer instead of piping to stdout
	Quiet bool
	// OverrideArtifacts makes the builder override any existing artifacts
	OverrideArtifacts bool
}

type RunnerContext struct {
	context.AppContext
	RunnerOptions
	Commands []*exec.Cmd
}

// Job runners should implement this interface, i.e.,
// the clean job and the build job implement a shared interface
type Runner interface {
	Run() error
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

	if options.TargetDirectory == "" {
		// when TargetDirectory is not specified, use the current working dir as target
		runner.TargetDirectory = context.Cwd
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

func (b *RunnerContext) Clone() *RunnerContext {
	cmds := []*exec.Cmd{}
	for _, srcCmd := range b.Commands {

		destCmd := exec.Command(
			srcCmd.Args[0],
			srcCmd.Args[1:]...,
		)
		destCmd.Dir = srcCmd.Dir
		out := &bytes.Buffer{}

		if b.Quiet {
			sink := bufio.NewWriter(out)
			destCmd.Stdout = sink
		} else {
			destCmd.Stdout = os.Stdout
		}
		cmds = append(cmds, destCmd)
	}
	return &RunnerContext{
		RunnerOptions: RunnerOptions{
			TargetDirectory:    b.TargetDirectory,
			Filename:           b.Filename,
			ArtifactsDirectory: b.ArtifactsDirectory,
			Quiet:              b.Quiet,
			OverrideArtifacts:  b.OverrideArtifacts,
		},
		Commands: cmds,
		AppContext: context.AppContext{
			Cwd:           b.Cwd,
			Root:          b.Root,
			Configuration: b.Configuration.Clone(),
		},
	}
}

func (r *RunnerContext) Build() *builder {
	return &builder{RunnerContext: r}
}

func (r *RunnerContext) Clean() *cleaner {
	return &cleaner{RunnerContext: r}
}
