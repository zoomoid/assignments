package cmd

import (
	"fmt"
	"os"

	"github.com/lithammer/dedent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/context"
)

var (
	rootLongDescription = dedent.Dedent(`
		Manangement CLI for doing university assignments with and around
		the "csassignments" LaTeX class.
	`)

	rootExample = dedent.Dedent(`
		# Initialize a directory for assignments
		assignmentctl bootstrap --course $COURSE_NAME \ 
			--group $GROUP_NAME --member "Max Mustermann;$ID" \ 
			--member "Erika Mustermann;$ID2" \
			--includes "code,feedback"

		# Generate a fresh assignment for the next number
		assignmentctl generate --due "$DUE_DATE"

		# Generate an assignment for a specific number
		assignmentctl generate 5 --due "$OTHER_DUE_DATE"

		# Build a specific assignment
		assignmentctl build 5 

		# Build a specific assignment while keeping intermediate files
		assignmentctl build 5 --keep

		# Build all assignments in the current directory, overriding existing ones
		assignmentctl build --all --force

		# Bundle a specific assignment to a zip file
		assignmentctl bundle 5 --zip 

		# Bundle all assignments to a tar.gz file
		assignmentctl bundle --all --tar --gzip

		# Create a template Gitlab CI pipeline file
		assignmentctl ci bootstrap gitlab -f .gitlab-ci.yml

		# Create a template Github Action file
		assignmentctl ci bootstrap github -f workflow.yml
	`)
)

var Version = "0.0.0-dev.0"

func Execute() {
	// program's working directory
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal().Msg("Failed to determine current working directory")
	}

	rootCmd := NewRootCommand(&rootOptions{
		root:    pwd,
		cwd:     pwd,
		verbose: false,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msgf("%v", err)
	}
}

type rootOptions struct {
	root    string
	cwd     string
	verbose bool
}

type rootData struct {
	*rootOptions
}

func NewRootCommand(opts *rootOptions) *cobra.Command {
	data := &rootData{
		rootOptions: opts,
	}

	ctx := &context.AppContext{
		Cwd:           data.cwd,
		Root:          data.root,
		Configuration: nil,
		Verbose:       data.verbose,
	}

	rootCmd := &cobra.Command{
		Use:           "assignmentctl",
		Short:         "assignments CLI for conveniently templating, building, and bundling course assignment",
		Long:          rootLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       rootExample,
		Version:       Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if data.verbose {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
			ctx.Verbose = data.verbose
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("assignments requires a subcommand to run")
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&data.verbose, options.Verbose, options.VerboseShort, false, "Sets logging verbosity level to high")

	rootCmd.AddCommand(NewBootstrapCommand(ctx, nil))
	rootCmd.AddCommand(NewGenerateCommand(ctx, nil))
	rootCmd.AddCommand(NewBuildCommand(ctx, nil))
	rootCmd.AddCommand(NewBundleCommand(ctx, nil))
	rootCmd.AddCommand(NewCiCommand(ctx))

	return rootCmd
}
