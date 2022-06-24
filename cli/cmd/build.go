package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
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
	keep  bool
	quiet bool
	file  string
}

func newBuildData() *buildData {
	return &buildData{
		force: false,
		all:   false,
		keep:  false,
		quiet: false,
		file:  "",
	}
}

func NewBuildCommand(ctx *context.AppContext, data *buildData) *cobra.Command {
	if data == nil {
		data = newBuildData()
	}

	err := ctx.Read()
	if err != nil {
		ctx.Logger.Fatalf("Failed to read config file", err)
	}

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build an assignment from source",
		Long:  buildLongDescription,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runs := []latexmk.RunnerOptions{}
			assignmentNo := ctx.Configuration.Status.Assignment

			if len(args) != 0 {
				i, err := strconv.Atoi(args[0])
				if err == nil {
					assignmentNo = uint32(i)
				}
			}

			if data.all && len(args) > 0 {
				return errors.New("cannot use --all flag with specific assignment")
			}

			if data.file != "" && len(args) > 0 {
				return errors.New("cannot use -f flag with specific assignment")
			}

			artifactsDirectory := filepath.Join(ctx.Root, "dist")

			if data.all {
				directories, err := filepath.Glob(filepath.Join(ctx.Root, "assignment-*"))
				if err != nil {
					return fmt.Errorf("failed to glob directories in %s, %v", ctx.Root, err)
				}
				for _, dir := range directories {
					filename := filepath.Join(dir, "assignment.tex")
					runs = append(runs, latexmk.RunnerOptions{
						TargetDirectory:    dir,
						Filename:           filename,
						ArtifactsDirectory: artifactsDirectory,
						Quiet:              data.quiet,
					})
				}
			} else {
				targetDirectory := filepath.Join(ctx.Root, fmt.Sprintf("assignment-%s", util.AddLeadingZero(assignmentNo)))
				filename := filepath.Join(targetDirectory, "assignment.tex")
				if data.file != "" {
					absPath, err := filepath.Abs(data.file)
					if err != nil {
						return err
					}

					// override targetDirectory and filename conditionally, depending on whether the argument is a file or a directory
					fi, err := os.Stat(absPath)
					if err != nil {
						return err
					}

					if fi.IsDir() {
						// append default filename
						targetDirectory = absPath
						filename = filepath.Join(absPath, "assignment.tex")
					} else {
						// file is a regular file
						filename = absPath
						targetDirectory = filepath.Dir(absPath)
					}
				}

				runs = []latexmk.RunnerOptions{{
					TargetDirectory:    targetDirectory,
					Filename:           filename,
					ArtifactsDirectory: artifactsDirectory,
					Quiet:              data.quiet,
				}}
			}

			startTime := time.Now()
			for _, run := range runs {
				runner, err := latexmk.New(ctx, &run)
				if err != nil {
					ctx.Logger.Errorf("failed to initialize runner for assignment %d in %s ", assignmentNo, run.Filename)
					return err
				}
				err = runner.Build().Run()
				if err != nil {
					ctx.Logger.Errorf("run failed for assignment %d in %s, %v", assignmentNo, run.Filename, err)
					ctx.Logger.Warnf("Leaving working directory %s dirty, might require manual cleanup", run.TargetDirectory)
					return err
				}

				if !data.keep {
					err = runner.Clean().Run()
					if err != nil {
						ctx.Logger.Errorf("failed to clean up for assignment %d in %s, %v", assignmentNo, run.Filename, err)
						ctx.Logger.Warnf("Leaving working directory %s dirty, might require manual cleanup", run.TargetDirectory)
						return err
					}
				}
			}
			ctx.Logger.Debug("Finished all build jobs successfully", "jobCount", len(runs), "duration", time.Since(startTime))
			return nil
		},
	}

	addBuildFlags(buildCmd.Flags(), data)

	return buildCmd

}

func addBuildFlags(flags *pflag.FlagSet, data *buildData) {
	flags.BoolVarP(&data.force, options.Force, options.ForceShort, false, "Override any existing assignments with the same name")
	flags.BoolVarP(&data.all, options.All, options.AllShort, false, "Build all assignments in assignment-*/")
	flags.BoolVar(&data.keep, options.Keep, false, "Skip latexmk -C cleaning up all files in the source directory")
	flags.BoolVar(&data.quiet, options.Quiet, false, "Suppress output from latexmk subprocesses")
	flags.StringVarP(&data.file, options.File, options.FileShort, "", "Specify a file to build, will override any derived behaviour from the repository's configmap")
}
