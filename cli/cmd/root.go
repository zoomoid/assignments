package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zoomoid/assignments/v1/internal/config"
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
		return fmt.Errorf("assignments requires a subcommand to run")
	},
}

var logger *zap.SugaredLogger
var cwd string
var configuration *config.Configuration

func Execute() {
	l, _ := zap.NewProduction()

	defer l.Sync()
	sugar := l.Sugar()

	dir, err := os.Getwd()

	if err != nil {
		logger.Fatal("Failed to determine current working directory")
	}

	cwd = dir
	logger = sugar

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}
