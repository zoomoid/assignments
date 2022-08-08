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

		After compilation, the command also cleans up any intermediate files
		created by the LaTeX compiler. By default, the cleanup will be done
		directly in the file system, by using Glob patterns for a large
		set of intermediate TeX files. You can override this behaviour in 
		two ways:

		1) you can specify a different set of Glob patterns in the config at
		.spec.build.cleanup.glob.patterns. If you'd like to run the cleanup
		recursively, also set .spec.build.cleanup.glob.recursive to true.
		Note that the Glob patterns are not merged with the default ones:
		If you provide your own, these are the complete ones to cleanup

		2) you can change the execution from using Globs to running commands,
		e.g. latexmk -C: For this, set .spec.build.cleanup.command.recipe
		accordingly

		Note that .spec.build.cleanup.command and .spec.build.cleanup.glob
		are mutually exclusive. Presence of both will cause the CLI to throw
		an error.
		
		If you use the build command in a setup different to one-off runs, 
		for which you might want to keep the files for later runs again to save 
		times, you can use --keep to preserve those intermediate files.

		You can suppress the output of spawned shell commands by passing 
		--quiet, or -q.

		To adjust the build recipe for compilation, add a recipe to your
		configuration file at .spec.build.recipe. Recipes are order-preservent
		lists of commands with arguments in YAML format. A recipe consists
		of Tools, which must at least contain a .command string, and may
		include arbitrary .args as a YAML list.
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

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build an assignment from source",
		Long:  buildLongDescription,
		Args:  cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			err := ctx.Read()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read config file")
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			defer ctx.Write()
		},
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
				var err error
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

	addBuildFlags(buildCmd.PersistentFlags(), data)

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
