package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/lithammer/dedent"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
)

var (
	rootLongDescription = dedent.Dedent(`
		Manangement CLI for doing university assignments with and around
		the "csassignments" LaTeX class.
	`)
)

func Execute() {
	// program's working directory
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to determine current working directory")
	}

	// if we cannot find a configuration file in here, traverse the file tree upwards
	// until either the root or we find a config file
	cfgPath, err := config.Find(pwd)
	if err != nil {
		log.Fatal("Failed to find configmap in working directory or above. Is the directory initialized?")
	}

	rootCmd := NewRootCommand(&rootOptions{
		root:    cfgPath,
		cwd:     pwd,
		verbose: false,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
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

	rootCmd := &cobra.Command{
		Use:              "assignmentctl",
		Short:            "assignments CLI for conveniently templating, building, and bundling course assignment",
		Long:             rootLongDescription,
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if data.verbose {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("assignments requires a subcommand to run")
		},
		SilenceUsage: true,
	}

	ctx := &context.AppContext{
		Cwd:           data.cwd,
		Root:          data.root,
		Configuration: nil,
	}

	rootCmd.Flags().BoolVarP(&data.verbose, options.Verbose, options.VerboseShort, false, "Sets logging verbosity level to high")

	rootCmd.AddCommand(NewBootstrapCommand(ctx, nil))
	rootCmd.AddCommand(NewGenerateCommand(ctx, nil))
	rootCmd.AddCommand(NewBuildCommand(ctx, nil))
	rootCmd.AddCommand(NewBundleCommand(ctx, nil))

	return rootCmd
}
