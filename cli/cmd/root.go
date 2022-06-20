package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/context"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "assignment",
	Short: "assignments CLI for conveniently templating, building, and bundling course assignment",
	Long: `💡 assignments CLI for conveniently templating, building, and bundling course assignment

	Run 'assignment bootstrap' to create a new environment {needed}
	Run 'assignment generate <number>' to template a new assignment
	Run 'assignment build <number>' to build a specfic assignment with latexmk (add '--all' to build all assignments)
	Run 'assignment compile <number>' to compile a zip file for a specific assignment (add '--all' to compile all assignments)
	Run 'assignment release' inside a Gitlab CI/CD pipeline to create files required for automatic release (see manual)

	See each command's help for description of arguments

	Copyright (C) zoomoid, 2022`,

	RunE: func(cmd *cobra.Command, args []string) error {
		var l *zap.Logger
		if verbose {
			l, _ = zap.NewDevelopment()

		} else {
			l, _ = zap.NewProduction()
		}
		defer l.Sync()
		sugar := l.Sugar()
		logger = sugar

		return fmt.Errorf("assignments requires a subcommand to run")
	},
}

var logger *zap.SugaredLogger
var cwd string
var verbose bool

func Execute() {

	rootCmd.Flags().BoolVarP(&verbose, options.Verbose, options.VerboseShort, false, "Prints debug logs")

	dir, err := os.Getwd()

	if err != nil {
		logger.Fatal("Failed to determine current working directory")
	}

	cwd = dir

	ctx := &context.AppContext{
		Logger:        logger,
		Cwd:           cwd,
		Configuration: nil,
	}

	rootCmd.AddCommand(NewBootstrapCommand(ctx, nil))
	rootCmd.AddCommand(NewGenerateCommand(ctx, nil))
	rootCmd.AddCommand(NewBuildCommand(ctx, nil))
	rootCmd.AddCommand(NewBundleCommand(ctx, nil))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
