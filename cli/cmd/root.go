package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/context"
	"go.uber.org/zap"
)

var (
	rootLongDescription = dedent.Dedent(`
		Manangement CLI for doing university assignments with and around
		the "csassignments" LaTeX class.
	`)
)

func Execute() {

	var verbose bool
	flag.BoolVar(&verbose, options.VerboseShort, false, "Prints debug logs")
	flag.Parse()

	// program's working directory
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to determine current working directory")
	}

	// if we cannot find a configuration file in here, traverse the file tree upwards
	// until either the root or we find a config file
	cfgPath, err := findConfigmapFile(pwd)
	if err != nil {
		log.Fatal("Failed to find configmap in working directory or above. Is the directory initialized?")
	}

	rootCmd := NewRootCommand(&rootOptions{
		root:    cfgPath,
		cwd:     pwd,
		verbose: verbose,
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
	logger *zap.SugaredLogger
}

func NewRootCommand(opts *rootOptions) *cobra.Command {

	logger, _ := makeLogger(opts.verbose)

	data := &rootData{
		rootOptions: opts,
		logger:      logger,
	}

	rootCmd := &cobra.Command{
		Use:   "assignment",
		Short: "assignments CLI for conveniently templating, building, and bundling course assignment",
		Long:  rootLongDescription,

		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("assignments requires a subcommand to run")
		},
	}

	ctx := &context.AppContext{
		Logger:        data.logger,
		Cwd:           data.cwd,
		Root:          data.root,
		Configuration: nil,
	}

	rootCmd.AddCommand(NewBootstrapCommand(ctx, nil))
	rootCmd.AddCommand(NewGenerateCommand(ctx, nil))
	rootCmd.AddCommand(NewBuildCommand(ctx, nil))
	rootCmd.AddCommand(NewBundleCommand(ctx, nil))

	return rootCmd
}

func findConfigmapFile(start string) (string, error) {
	p := start
	for {
		_, err := os.Stat(p)
		if err != nil {
			if filepath.Base(filepath.ToSlash(p)) == "/" {
				break
			}
			// move up one level and try again
			p = filepath.Clean(fmt.Sprintf("%s/../", p))
			continue
		}
		// found configmap file at p
		return p, nil
	}
	return "", errors.New("failed to find configmap in working directory or above")
}

func makeLogger(verbose bool) (*zap.SugaredLogger, error) {
	var l *zap.Logger
	var err error
	if verbose {
		l, err = zap.NewDevelopment()

	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}
	defer l.Sync()
	return l.Sugar(), nil
}
