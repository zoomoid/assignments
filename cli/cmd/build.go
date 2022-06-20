package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
	"github.com/zoomoid/assignments/v1/internal/latexmk"
	"github.com/zoomoid/assignments/v1/internal/util"
)

var (
	buildLongDescription = dedent.Dedent(`
		The command builds a selected assignment, either from arguments,
		or from the state of the local configuration file using latexmk
		with the underlying LaTeX distro. After successful build, the
		artifact files are copied to a common output directory, commonly
		./dist/.

		To build *all* assignments found in the working directory, add the
		--all (or -a) flag.

		By default, latexmk is run 3 times. You can change this by specifying
		--runs with another number of rounds.

		After compilation, the command also cleans up any intermediate files
		created by the LaTeX compiler, using latexmk -C. If you use the build
		command in a setup different to one-off runs, for which you might
		want to keep the files for later runs again to save times, you
		can use --keep to preserve those intermediate files.

		You can suppress the output of latexmk by passing --quiet, or -q.

		To adjust the build recipe for compilation, add a recipe to your
		configuration file at .spec.building.recipe.
	`)
)

type buildData struct {
	force bool
	all   bool
	runs  uint
	keep  bool
	quiet bool
}

func newBuildData() *buildData {
	return &buildData{
		force: false,
		all:   false,
		runs:  uint(3),
		keep:  false,
		quiet: false,
	}
}

func NewBuildCommand(ctx *context.AppContext, data *buildData) *cobra.Command {
	if data == nil {
		data = newBuildData()
	}

	cfg, err := config.ReadConfigMap()

	if err != nil {
		ctx.Logger.Fatalf("Failed to read config file, %v", err)
	}

	ctx.Configuration = cfg

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build an assignment from source",
		Long:  buildLongDescription,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runs := []latexmk.RunnerOptions{}
			assignmentNo := cfg.Status.Assignment

			if len(args) != 0 {
				i, err := strconv.Atoi(args[0])
				if err == nil {
					assignmentNo = uint32(i)
				}
			}

			if data.all && len(args) > 0 {
				return errors.New("cannot use --all flag with specific assignment")
			}

			artifactsDirectory := filepath.Join(ctx.Cwd, "dist")

			if !data.all {
				targetDirectory := filepath.Join(ctx.Cwd, fmt.Sprintf("assignment-%s", util.AddLeadingZero(assignmentNo)))
				filename := filepath.Join(targetDirectory, "assignment.tex")

				runs = []latexmk.RunnerOptions{{
					TargetDirectory:    targetDirectory,
					Filename:           filename,
					Runs:               int(data.runs),
					ArtifactsDirectory: artifactsDirectory,
					Quiet:              data.quiet,
				}}
			} else {
				directories, err := filepath.Glob(filepath.Join(ctx.Cwd, "assignment-*"))
				if err != nil {
					return fmt.Errorf("failed to glob directories in %s, %v", ctx.Cwd, err)
				}
				for _, dir := range directories {
					filename := filepath.Join(dir, "assignment.tex")
					runs = append(runs, latexmk.RunnerOptions{
						TargetDirectory:    dir,
						Filename:           filename,
						Runs:               int(data.runs),
						ArtifactsDirectory: artifactsDirectory,
						Quiet:              data.quiet,
					})
				}
			}

			for i, run := range runs {
				runner, err := latexmk.New(ctx, &run)
				if err != nil {
					ctx.Logger.Errorf("failed to initialize runner for assignment %d in %s ", assignmentNo, run.Filename)
					return err
				}
				_, err = runner.RunBuildMultiple()
				if err != nil {
					ctx.Logger.Errorf("run failed for assignment %d in %s, %v", assignmentNo, run.Filename, err)
					ctx.Logger.Warnf("Leaving working directory %s dirty, might require manual cleanup", run.TargetDirectory)
					return err
				}

				if !data.keep {
					err = runner.RunClean()
					if err != nil {
						ctx.Logger.Errorf("failed to clean up for assignment %d in %s, %v", assignmentNo, run.Filename, err)
						ctx.Logger.Warnf("Leaving working directory %s dirty, might require manual cleanup", run.TargetDirectory)
						return err
					}
				}

				ctx.Logger.Infof("Finished building %s [%d/%d]", run.Filename, i+1, len(runs))
			}

			return nil
		},
	}

	addBuildFlags(buildCmd.Flags(), data)

	return buildCmd

}

func addBuildFlags(flags *pflag.FlagSet, data *buildData) {
	flags.BoolVarP(&data.force, options.Force, options.ForceShort, false, "Override any existing assignments with the same name")
	flags.BoolVarP(&data.all, options.All, options.AllShort, false, "Build all assignments in assignment-*/")
	flags.UintVarP(&data.runs, options.Runs, options.RunsShort, 3, "latexmk compiler runs")
	flags.BoolVar(&data.keep, options.Keep, false, "Skip latexmk -C cleaning up all files in the source directory")
	flags.BoolVar(&data.quiet, options.Quiet, false, "Suppress output from latexmk subprocesses")
}
