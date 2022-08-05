package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lithammer/dedent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/context"
	"github.com/zoomoid/assignments/v1/internal/runner"
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
		log.Fatal().
			Err(err).
			Msg("Failed to read config file")
	}
	defer ctx.Write()

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build an assignment from source",
		Long:  buildLongDescription,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runs := []runner.RunnerOptions{}
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

			if ctx.Configuration.Spec.BuildOptions.Cleanup != nil &&
				ctx.Configuration.Spec.BuildOptions.Cleanup.Command != nil &&
				ctx.Configuration.Spec.BuildOptions.Cleanup.Glob != nil {
				return errors.New("found ambiguous cleanup mode, only use either glob or command")
			}

			if data.all {
				directories, err := filepath.Glob(filepath.Join(ctx.Root, "assignment-*"))
				if err != nil {
					return fmt.Errorf("failed to glob directories in %s, %v", ctx.Root, err)
				}
				for _, dir := range directories {
					filename := "assignment.tex"
					runs = append(runs, runner.RunnerOptions{
						TargetDirectory:   filepath.Base(dir),
						Filename:          filename,
						Quiet:             data.quiet,
						OverrideArtifacts: data.force,
					})
				}
			} else {
				targetDirectory := fmt.Sprintf("assignment-%s", util.AddLeadingZero(assignmentNo))
				filename := "assignment.tex"
				if data.file != "" {
					targetDirectory, filename, err = targetDirectoryFromFlag(data.file)
					if err != nil {
						return err
					}
				}

				runs = []runner.RunnerOptions{{
					TargetDirectory:   targetDirectory,
					Filename:          filename,
					Quiet:             data.quiet,
					OverrideArtifacts: data.force,
				}}
			}

			startTime := time.Now()
			for _, run := range runs {
				runner, err := runner.New(ctx, &run)
				if err != nil {
					log.Error().Msgf("failed to initialize runner for %s", run.Filename)
					return err
				}
				err = runner.Build().Run()
				if err != nil {
					log.Error().Err(err).Msgf("run failed for %s", run.Filename)
					log.Warn().Msgf("Leaving working directory %s dirty, might require manual cleanup", run.TargetDirectory)
					return err
				}

				if !data.keep {
					err = runner.Clean().Run()
					if err != nil {
						log.Error().Err(err).Msgf("failed to clean up for %s", run.Filename)
						log.Warn().Msgf("Leaving working directory %s dirty, might require manual cleanup", run.TargetDirectory)
						return err
					}
				}
			}
			log.Debug().
				Dur("duration", time.Since(startTime)).
				Int("jobCount", len(runs)).
				Msg("Finished all build jobs successfully")
			return nil
		},
	}

	addBuildFlags(buildCmd.Flags(), data)

	return buildCmd
}

// targetDirectoryFromFlag uses the --file flag to determine the targetDirectory and filename to build.
// It returns a triplet containing (targetDirectory, filename, error), where error is not nil iff and only if
// an error occured during file system interaction, i.e. expanding the path to absolute, and probing the path.
func targetDirectoryFromFlag(file string) (string, string, error) {
	var targetDirectory string
	var filename string
	// augment path to be absolute, e.g. "-f ." is expanded to "CWD/assignment.tex"
	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", "", err
	}

	// override targetDirectory and filename conditionally, depending on whether the argument is a file or a directory
	fi, err := os.Stat(absPath)
	if err != nil {
		return "", "", err
	}
	if fi.IsDir() {
		// passed in a directory, defaulting filename to assignment.tex and the absolute path created from the flag
		// NOTE that e.g. runner.SetRoot() WILL NOT change the target directory when used this way, as the path
		// pased down to the runner is absolute.
		targetDirectory = absPath
		filename = "assignment.tex"
	} else {
		// file is a regular file, use filepath.Dir() for target directory and filepath.Base() to make filename relative/local again
		targetDirectory = filepath.Dir(absPath)
		filename = filepath.Base(absPath)
	}
	return targetDirectory, filename, nil
}

func addBuildFlags(flags *pflag.FlagSet, data *buildData) {
	flags.BoolVar(&data.force, options.Force, false, "Override any existing assignments with the same name")
	flags.BoolVarP(&data.all, options.All, options.AllShort, false, "Build all assignments in assignment-*/")
	flags.BoolVar(&data.keep, options.Keep, false, "Skip latexmk -C cleaning up all files in the source directory")
	flags.BoolVar(&data.quiet, options.Quiet, false, "Suppress output from latexmk subprocesses")
	flags.StringVarP(&data.file, options.File, options.FileShort, "", "Specify a file to build, will override any derived behaviour from the repository's configmap")
}
