package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/template"
)

// addLeadingZero prepends numbers smaller than 10 with a leading zero
func addLeadingZero(assignment uint32) string {
	if assignment < 10 {
		return fmt.Sprintf("0%d", assignment)
	}
	return fmt.Sprintf("%d", assignment)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate new assignments from the template",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			// update configuration status
			i, err := strconv.Atoi(args[0])
			if err == nil && !noIncrement {
				configuration.Status.Assignment = uint32(i)
			}
		}

		due := promptDueDate()

		spec := configuration.Spec

		bindings := &template.TemplateBinding{
			ClassPath: "../assignments",
			Course:    spec.Course,
			Group:     spec.Group,
			Sheet:     addLeadingZero(configuration.Status.Assignment),
			Due:       due,
			Members:   spec.Members,
		}

		sheetSource, err := template.GenerateAssignmentTemplate(bindings)

		if err != nil {
			logger.Errorf("Failed to template TeX source code, %v", err)
			os.Exit(1)
		}

		assignmentDirectory := fmt.Sprintf("assignment-%s", bindings.Sheet)

		err = os.Mkdir(assignmentDirectory, os.ModeDir)

		if err != nil {
			logger.Errorf("Failed to create assignment directory, %v", err)
			os.Exit(1)
		}

		file := filepath.Join(assignmentDirectory, "assignment.tex")

		err = os.WriteFile(file, []byte(sheetSource), 0644)
		if err != nil {
			logger.Errorf("Failed to write to file, %v", err)
			os.Exit(1)
		}

		logger.Infof("Generated assignment at %s", file)

		defer config.WriteConfigMap(configuration)
	},
}

func promptDueDate() string {
	fmt.Print("⏱️  When is the assignment due? (e.g.,'April 20, 2021): ")
	due := ""
	fmt.Scanln("%s", due)
	return due
}

var noIncrement bool
var force bool

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().BoolVar(&noIncrement, "no-increment", false, "Skip incrementing assignment number in configuration")
	generateCmd.Flags().BoolVarP(&force, "force", "f", false, "Overrides any existing assignment source files")

	cfg, err := config.ReadConfigMap()

	if err != nil {
		logger.Fatalf("Failed to read config file, %v", err)
	}

	configuration = cfg
}
