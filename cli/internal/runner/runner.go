package runner

import (
	"bufio"
	"bytes"
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
	// ArtifactsDirectory to copy to, defaults to "$PWD/dist"
	ArtifactsDirectory string
	// Quiet makes the latexmk run capture the output inside a buffer instead of piping to stdout
	Quiet bool
	// OverrideArtifacts makes the builder override any existing artifacts
	OverrideArtifacts bool
}

type RunnerContext struct {
	*context.AppContext
	options            *RunnerOptions
	filename           string
	quiet              bool
	overrideArtifacts  bool
	targetDirectory    string
	artifactsDirectory string
	Commands           []*exec.Cmd
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
		AppContext: context,
		options:    options,
	}

	if options.TargetDirectory == "" {
		// when TargetDirectory is not specified, use the current working dir as target
		runner.targetDirectory = context.Cwd
	} else {
		if filepath.IsAbs(options.TargetDirectory) {
			runner.targetDirectory = options.TargetDirectory
		} else {
			runner.targetDirectory = filepath.Join(context.Root, options.TargetDirectory)
		}
	}

	if options.ArtifactsDirectory == "" {
		runner.artifactsDirectory = filepath.Join(context.Root, "dist")
	}

	if options.Filename == "" {
		runner.filename = options.Filename
	} else {
		runner.filename = "assignment.tex"
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

		if b.quiet {
			sink := bufio.NewWriter(out)
			destCmd.Stdout = sink
		} else {
			destCmd.Stdout = os.Stdout
		}
		cmds = append(cmds, destCmd)
	}
	return &RunnerContext{
		targetDirectory:    b.TargetDirectory(),
		filename:           b.Filename(),
		artifactsDirectory: b.ArtifactsDirectory(),
		quiet:              b.Quiet(),
		overrideArtifacts:  b.OverrideArtifacts(),
		Commands:           cmds,
		AppContext: &context.AppContext{
			Cwd:           b.Cwd,
			Root:          b.Root,
			Configuration: b.Configuration.Clone(),
			Logger:        b.Logger,
		},
	}
}

func (r *RunnerContext) Build() *builder {
	return &builder{RunnerContext: r}
}

func (r *RunnerContext) Clean() *cleaner {
	return &cleaner{RunnerContext: r}
}

func (r *RunnerContext) ArtifactsDirectory() string {
	return r.artifactsDirectory
}

func (r *RunnerContext) TargetDirectory() string {
	return r.targetDirectory
}

func (r *RunnerContext) Filename() string {
	return r.filename
}

func (r *RunnerContext) Quiet() bool {
	return r.quiet
}

func (r *RunnerContext) OverrideArtifacts() bool {
	return r.overrideArtifacts
}

func (r *RunnerContext) SetTargetDirectory(targetDirectory string) {
	if filepath.IsAbs(targetDirectory) {
		r.targetDirectory = targetDirectory
	} else {
		r.targetDirectory = filepath.Join(r.Root, r.options.TargetDirectory)
	}
}

func (r *RunnerContext) SetRoot(newRoot string) {
	r.Root = newRoot

	// update root-dependent fields for runner
	if r.options.TargetDirectory != "" {
		if !filepath.IsAbs(r.options.TargetDirectory) {
			r.targetDirectory = filepath.Join(r.Root, r.options.TargetDirectory)
		}
	}
	if r.options.ArtifactsDirectory == "" {
		r.artifactsDirectory = filepath.Join(r.Root, "dist")
	}
}
