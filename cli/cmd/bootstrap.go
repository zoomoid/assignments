package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/bundle"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
	"github.com/zoomoid/assignments/v1/internal/runner"
	"github.com/zoomoid/assignments/v1/internal/template"
)

var (
	bootstrapLongDescription = dedent.Dedent(`
	Running bootstrap will create a local configuration file for the current
	working directory containing minimal information about the course,
	passed in either as flags or during interactive prompts.

	The configuration file can be customized further afterwards, e.g., by
	adding different build recipes and bundling options. For this, see
	documentation.
	`)

	instructionsPreamble = dedent.Dedent(`
	Completed initialization of directory with assignmentctl.

	Configuration is stored in .assignments.yaml. You can start customizing
	it to your liking. 

	You can start a new assignment by running 
	
	  assignmentctl generate

	`)

	instructionsGit = dedent.Dedent(`
	Also created a git repository for your assignments to be stored in.
	I've already created a commit containing your configuration file and
	setup your branch accordingly.

	Now is probably a good time to add a remote:
	
	  git remote add origin $GIT_REPOSITORY

	`)

	instructionsAppendix = dedent.Dedent(`
	Best of luck to you and keep it up!
	`)
)

type bootstrapData struct {
	course   string
	group    string
	members  []string
	includes []string
	cfg      *config.Configuration
	git      bool
	full     bool
}

func newBootstrapData() *bootstrapData {
	return &bootstrapData{
		course:  "",
		group:   "",
		members: []string{},
		full:    false,
		git:     false,
		cfg:     config.Minimal(),
	}
}

func NewBootstrapCommand(ctx *context.AppContext, data *bootstrapData) *cobra.Command {
	if data == nil {
		data = newBootstrapData()
	}

	bootstrapCommand := &cobra.Command{
		Use:     "bootstrap",
		Short:   "Bootstraps the repository with a configuration file ready to go",
		Long:    bootstrapLongDescription,
		Aliases: []string{"init", "initialize"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if data.course == "" {
				data.course = promptCourseName()
			}
			if data.group == "" {
				data.group = promptGroupName()
			}
			if data.course == "" {
				data.course = promptCourseName()
			}

			if len(data.members) == 0 {
				data.members = promptMembers()
			}

			data.cfg.Spec.Course = data.course
			data.cfg.Spec.Group = data.group

			members := []config.GroupMember{}

			for _, m := range data.members {
				m, err := transformMember(m)
				if err != nil {
					log.Warn().Interface("raw", m).Msg("group member name is illformed, skipping...")
					break
				}
				members = append(members, *m)
			}

			data.cfg.Spec.Members = members

			includes := []config.Include{}

			for _, include := range data.includes {
				includes = append(includes, config.Include{
					Path: "../" + include,
				})
			}

			data.cfg.Spec.Includes = includes

			if data.full {
				data.cfg = augmentDefaults(data.cfg)
			}

			ctx.Configuration = data.cfg
			err := ctx.Write()
			if err != nil {
				return err
			}

			if data.git {
				err := createGitRepository(ctx.Verbose)
				if err != nil {
					return err
				}
				// Log instructions with default print to prevent prefixing of timestamp etc. from zerolog

			}

			fmt.Print(instructionsPreamble)
			if data.git {
				fmt.Print(instructionsGit)
			}
			fmt.Print(instructionsAppendix)

			return nil
		},
	}

	addBootstrapFlags(bootstrapCommand.PersistentFlags(), data)

	return bootstrapCommand
}

func transformMember(member string) (*config.GroupMember, error) {
	s := strings.Split(member, ";")

	if len(s) == 2 {
		return &config.GroupMember{
			Name: s[0],
			ID:   s[1],
		}, nil
	}

	if len(s) == 1 {
		return &config.GroupMember{
			Name: s[0],
			ID:   "",
		}, nil
	}

	return nil, errors.New("group member name is illformed")

}

func addBootstrapFlags(flags *pflag.FlagSet, data *bootstrapData) {
	flags.StringVar(&data.course, options.CourseName, "", "Course name")
	flags.StringVar(&data.group, options.GroupName, "", "Group name")
	flags.StringSliceVar(&data.members, options.Members, []string{}, "Group members, as comma-separated <Name>;<ID> tuples")
	flags.StringSliceVar(&data.includes, options.Includes, []string{}, "Custom TeX includes for the template. Paths are relative to the REPOSITORY root, not the actual assignment source file")
	flags.BoolVar(&data.git, options.Git, false, "Create a git repository in the current directory and commit the configuration file immediately")
	flags.BoolVar(&data.full, options.Full, false, "Include all defaults in configuration file")
}

func promptCourseName() string {
	fmt.Print("❓ Please enter the course's name: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal().Err(err)
	}
	input = strings.TrimSuffix(input, "\n")
	return input
}

func promptGroupName() string {
	fmt.Print("❓ Please enter a group name (or leave empty to use nothing): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal().Err(err)
	}
	input = strings.TrimSuffix(input, "\n")
	return input
}

func promptMembers() []string {
	m := make([]string, 0)
	for {
		fmt.Print("❓ Please enter a group member's name followed by their student ID (e.g. Max Mustermann;123456), or press 'q' to move on: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal().Err(err)
		}
		input = strings.TrimSuffix(input, "\n")
		if input == "q" {
			break
		}

		m = append(m, strings.TrimSpace(input))
	}
	return m
}

func createGitRepository(verbose bool) error {
	var stdout *os.File
	if verbose {
		stdout = os.Stdout
	}

	initCmd := exec.Command("git", "init", ".")
	initCmd.Stdout = stdout
	if err := initCmd.Run(); err != nil {
		return err
	}

	addCmd := exec.Command("git", "add", ".assignments.yaml")
	addCmd.Stdout = stdout
	if err := addCmd.Run(); err != nil {
		return err
	}

	commitCmd := exec.Command("git", "commit", "-s", "-m", "chore(*): initialize repository with assignmentctl")
	commitCmd.Stdout = stdout
	if err := commitCmd.Run(); err != nil {
		return err
	}

	branchCmd := exec.Command("git", "branch", "-M", "main")
	branchCmd.Stdout = stdout
	if err := branchCmd.Run(); err != nil {
		return err
	}

	return nil
}

func augmentDefaults(cfg *config.Configuration) *config.Configuration {
	return &config.Configuration{
		Spec: &config.ConfigurationSpec{
			Course:   cfg.Spec.Course,
			Group:    cfg.Spec.Group,
			Members:  cfg.Spec.Members,
			Includes: cfg.Spec.Includes,
			Template: template.DefaultSheetTemplate,
			GenerateOptions: &config.GenerateOptions{
				Create: []string{},
			},
			BuildOptions: &config.BuildOptions{
				BuildRecipe: &config.Recipe{
					{
						Command: runner.DefaultBuildProgram,
						Args:    runner.DefaultBuildArgs,
					},
				},
				Cleanup: runner.DefaultCleaner,
			},
			BundleOptions: &config.BundleOptions{
				Template: bundle.DefaultArchiveNameTemplate,
				Data:     make(map[string]interface{}),
				Include:  []string{},
			},
		},
		Status: &config.ConfigurationStatus{
			Assignment: 1,
		},
	}
}
