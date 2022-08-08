package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
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
		
		The command creates a new directory from the current assignment number,
		as well as all directories defined in the .spec.generate.create list.
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

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate new assignments from the template",
		Long:  generateLongDescription,
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
			if len(args) != 0 {
				// update configuration status
				i, err := strconv.Atoi(args[0])
				if err == nil && !data.noIncrement {
					ctx.Configuration.Status.Assignment = uint32(i)
				}
			}

			due := data.due
			if data.due == "" {
				due = promptDueDate()
			}

			spec := ctx.Configuration.Spec

			var classPath string

			var tpl string

			if spec.Template != "" {
				tpl = dedent.Dedent(spec.Template)
			}

			bindings := &template.TemplateBinding{
				ClassPath: classPath,
				Course:    spec.Course,
				Group:     spec.Group,
				Sheet:     util.AddLeadingZero(ctx.Configuration.Status.Assignment),
				Due:       due,
				Members:   spec.Members,
				Includes:  spec.Includes,
			}

			sheetSource, err := template.GenerateAssignmentTemplate(&tpl, bindings)

			if err != nil {
				return err
			}

			// create the assignment's main directory
			assignmentDirectory := fmt.Sprintf("assignment-%s", bindings.Sheet)
			if data.force {
				// when using --force to override any existing assignments, clean up before creating
				// assignment directory from scratch
				_ = os.RemoveAll(assignmentDirectory)
			}
			err = os.Mkdir(assignmentDirectory, 0777)
			if err != nil {
				return err
			}

			// create the additional directories defined in the spec
			additionalDirectories := []string{}
			if spec.GenerateOptions != nil {
				additionalDirectories = spec.GenerateOptions.Create
			}
			for _, dir := range additionalDirectories {
				err = os.Mkdir(filepath.Join(assignmentDirectory, dir), 0777)
				if err != nil {
					return err
				}
			}

			file := filepath.Join(assignmentDirectory, "assignment.tex")

			err = os.WriteFile(file, sheetSource.Bytes(), 0644)
			if err != nil {
				return err
			}

			if !data.noIncrement {
				ctx.Configuration.Status.Assignment += 1
			}

			log.Info().Msgf("Generated assignment at %s", file)

			defer ctx.Write()
			return nil
		},
	}

	addGenerateFlags(generateCmd.PersistentFlags(), data)

	return generateCmd
}

func addGenerateFlags(flags *pflag.FlagSet, data *generateData) {
	flags.BoolVar(&data.noIncrement, options.NoIncrement, false, "Skip incrementing assignment number in configuration")
	flags.BoolVarP(&data.force, options.Force, options.ForceShort, false, "Overrides any existing assignment source files")
	flags.StringVar(&data.due, options.Due, "", "Due date of the assignment to generate. If not provided, you'll be prompted for a due date")
}

func promptDueDate() string {
	fmt.Print("⏱️  When is the assignment due? (e.g.,'April 20, 2021): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal().Err(err)
	}
	input = strings.TrimSuffix(input, "\n")
	return input
}
