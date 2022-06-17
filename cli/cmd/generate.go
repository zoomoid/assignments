package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
	"github.com/zoomoid/assignments/v1/internal/template"
	"github.com/zoomoid/assignments/v1/internal/util"

	"github.com/lithammer/dedent"
)

var (
	generateLongDescription = dedent.Dedent(`
		The command generates a new assignment, either given by number
		as an argument to the command, or otherwise from the local 
		configuration file, which keeps track of the upstream assignment.

		Generating (or templating) a new assignment requires a due date.
		As this is usually given, you can either use the --due flag, or
		wait for the CLI to prompt you. If however the due date is *not*
		provided by the assignment, just pressing ENTER during the prompt
		will leave it empty and thus not printed in the assignment's
		header.
		
		You can make the command skip incrementing the status counter in
		the local configuration file by passing the --no-increment flag.
		
		If there already exists an assignment in the target directory,
		the command will exit with an error. If however you pass the 
		--force flag, any files in the target directory will be overriden.
		Be careful!
		
		The default template for new assignments looks like this:
		
		\documentclass{csassignments}
		{{- range $_, $input := .Includes -}}
		\input{ {{- $input -}} }
		{{ end }}
		\course{ {{- .Course -}} }
		\group{ {{- .Group | default "" -}} }
		\sheet{ {{- .Sheet | default "" -}} }
		\due{ {{- .Due | default "" -}} }
		{{- range $_, $member := .Members }}
		{{- $firstname := ($member.Name | splitList " " | initial | join " ") | default "" -}}
		{{- $lastname := ($member.Name | splitList " " | last) | default "" -}} 
		\member{ {{- $firstname -}} }{ {{- $lastname -}} }{ {{- $member.ID -}} }
		{{ end }}
		\begin{document}
		\maketitle
		\gradingtable
		
		% Start the assignment here
		
		\end{document}
		
		You can provide your own template from the configuration file, by 
		setting .spec.template to a Golang template. You can use any Sprig 
		template function in your custom template.
		
		The command creates a new directory from the current assignment number
	`)
)

type generateData struct {
	noIncrement bool
	force       bool
	due         string
}

func newGenerateData() *generateData {
	return &generateData{
		noIncrement: false,
		force:       false,
		due:         "",
	}
}

func NewGenerateCommand(ctx *context.AppContext, data *generateData) *cobra.Command {
	if data == nil {
		data = newGenerateData()
	}

	cfg, err := config.ReadConfigMap()

	if err != nil {
		ctx.Logger.Fatalf("Failed to read config file, %v", err)
	}

	ctx.Configuration = cfg

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate new assignments from the template",
		Long:  generateLongDescription,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				// update configuration status
				i, err := strconv.Atoi(args[0])
				if err == nil && !data.noIncrement {
					cfg.Status.Assignment = uint32(i)
				}
			}

			due := promptDueDate()

			spec := cfg.Spec

			var classPath string

			var tpl string

			if spec.Template != "" {
				tpl = dedent.Dedent(spec.Template)
			}

			bindings := &template.TemplateBinding{
				ClassPath: classPath,
				Course:    spec.Course,
				Group:     spec.Group,
				Sheet:     util.AddLeadingZero(cfg.Status.Assignment),
				Due:       due,
				Members:   spec.Members,
				Includes:  spec.Includes,
			}

			sheetSource, err := template.GenerateAssignmentTemplate(&tpl, bindings)

			if err != nil {
				ctx.Logger.Errorf("Failed to template TeX source code, %v", err)
				return err
			}

			assignmentDirectory := fmt.Sprintf("assignment-%s", bindings.Sheet)

			err = os.Mkdir(assignmentDirectory, os.ModeDir)

			if err != nil {
				ctx.Logger.Errorf("Failed to create assignment directory, %v", err)
				return err
			}

			file := filepath.Join(assignmentDirectory, "assignment.tex")

			err = os.WriteFile(file, []byte(sheetSource), 0644)
			if err != nil {
				ctx.Logger.Errorf("Failed to write to file, %v", err)
				return err
			}

			ctx.Logger.Infof("Generated assignment at %s", file)

			defer config.WriteConfigMap(cfg)
			return nil
		},
	}

	addGenerateFlags(generateCmd.Flags(), data)

	return generateCmd
}

func addGenerateFlags(flags *pflag.FlagSet, data *generateData) {
	flags.BoolVar(&data.noIncrement, options.NoIncrement, false, "Skip incrementing assignment number in configuration")
	flags.BoolVarP(&data.force, options.Force, options.ForceShort, false, "Overrides any existing assignment source files")
	flags.StringVar(&data.due, options.Due, "", "Due date of the assignment to generate. If not provided, you'll be prompted for a due date")
}

func promptDueDate() string {
	fmt.Print("⏱️  When is the assignment due? (e.g.,'April 20, 2021): ")
	due := ""
	fmt.Scanln("%s", due)
	return due
}
