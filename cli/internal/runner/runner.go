package runner

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
	"go.uber.org/zap"
)

// RunnerOptions struct to carry configuration for the latexmk runs
type RunnerOptions struct {
	// Directory to run latexmk from
	TargetDirectory string
	// TeX source file to compile, defaults to "assignment.tex"
	Filename string
	// Quiet makes the latexmk run capture the output inside a buffer instead of piping to stdout
	Quiet bool
	// OverrideArtifacts makes the builder override any existing artifacts
	OverrideArtifacts bool
}

type RunnerContext struct {
	root               string
	cwd                string
	configuration      *config.Configuration
	logger             *zap.SugaredLogger
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
	MakeCommand() ([]*exec.Cmd, error)
	Run() error
}

var (
	defaultProgram        = "latexmk"
	defaultLatexmkOptions = []string{"-pdf", "-interaction=nonstopmode", "-file-line-error", "-shell-escape"}
)

// New creates a new runner context from the given parameters and applies sensible defaults
func New(context *context.AppContext, options *RunnerOptions) (*RunnerContext, error) {
	// clone the app context for the runner context to not mutate the application context's
	// state with setters on the runner
	runnerCtx := context.Clone()
	runner := &RunnerContext{
		options:       options,
		root:          runnerCtx.Root,
		cwd:           runnerCtx.Cwd,
		configuration: runnerCtx.Configuration,
		logger:        runnerCtx.Logger,
	}

	if options.TargetDirectory == "" {
		// when TargetDirectory is not specified, use the current working dir as target
		runner.targetDirectory = context.Cwd
	} else {
		runner.targetDirectory = options.TargetDirectory
	}

	runner.artifactsDirectory = "dist"

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
		targetDirectory:    b.targetDirectory,
		filename:           b.filename,
		artifactsDirectory: b.artifactsDirectory,
		quiet:              b.quiet,
		overrideArtifacts:  b.overrideArtifacts,
		Commands:           cmds,
		cwd:                b.cwd,
		root:               b.root,
		configuration:      b.configuration.Clone(),
		logger:             b.logger, // this is not actually cloned
	}
}

func (r *RunnerContext) Build() *builder {
	b := &builder{RunnerContext: r}
	return b
}

func (r *RunnerContext) Clean() *cleaner {
	c := &cleaner{RunnerContext: r}
	return c
}

func (r *RunnerContext) ArtifactsDirectory() string {
	if filepath.IsAbs(r.artifactsDirectory) {
		return r.artifactsDirectory
	} else {
		return filepath.Join(r.root, r.artifactsDirectory)
	}
}

func (r *RunnerContext) TargetDirectory() string {
	if filepath.IsAbs(r.targetDirectory) {
		return r.targetDirectory
	} else {
		return filepath.Join(r.root, r.targetDirectory)
	}
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
	if targetDirectory == "" {
		targetDirectory = r.cwd
	}
	r.targetDirectory = targetDirectory
}

func (r *RunnerContext) SetArtifactsDirectory(artifactsDirectory string) {
	if artifactsDirectory == "" {
		artifactsDirectory = "dist"
	}
	r.artifactsDirectory = artifactsDirectory
}

func (r *RunnerContext) SetRoot(newRoot string) {
	r.root = newRoot
}
