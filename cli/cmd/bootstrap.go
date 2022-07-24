package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
)

type bootstrapData struct {
	course   string
	group    string
	members  []string
	includes []string
	cfg      *config.Configuration
}

func newBootstrapData() *bootstrapData {
	return &bootstrapData{
		course:  "",
		group:   "",
		members: []string{},
		cfg:     &config.Configuration{},
	}
}

func NewBootstrapCommand(ctx *context.AppContext, data *bootstrapData) *cobra.Command {
	if data == nil {
		data = newBootstrapData()
	}

	bootstrapCommand := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstraps the repository with a configuration file ready to go",
		Long:  "",
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
					ctx.Logger.Warn("group member name is illformed", "raw", m)
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

			ctx.Configuration = data.cfg
			return ctx.Write()
		},
	}

	addBootstrapFlags(bootstrapCommand.Flags(), data)

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
}

func promptCourseName() string {
	fmt.Print("❓ Please enter the course's name: ")
	c := ""
	fmt.Scanln("%s", &c)
	return c
}

func promptGroupName() string {
	fmt.Print("❓ Please enter a group name (or leave empty to use nothing): ")
	g := ""
	fmt.Scanln("%s", &g)
	return g
}

func promptMembers() []string {
	m := make([]string, 0)
	for {
		fmt.Print("❓ Please enter a group member's name followed by their student ID (e.g. Max Mustermann;123456), or press 'q' to move on: ")
		rawMember := ""
		fmt.Scanln("%s", &rawMember)
		if rawMember == "q" {
			break
		}

		m = append(m, strings.TrimSpace(rawMember))
	}
	return m
}
