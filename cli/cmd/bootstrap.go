package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var bootstrapCommand = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstraps an environment with LaTeX class and .assignments.rc file",
	Run: func(cmd *cobra.Command, args []string) {
		if course == "" {
			course = promptCourseName()
		}

		if group == "" {
			group = promptGroupName()
		}

		if len(members) == 0 {
			members = promptMembers()
		}
	},
}

var course string
var group string
var members []string

func init() {
	rootCmd.AddCommand(bootstrapCommand)

	bootstrapCommand.Flags().StringVar(&course, "course", "", "Course name")
	bootstrapCommand.Flags().StringVar(&group, "group", "", "Group name")
	bootstrapCommand.Flags().StringArrayVar(&members, "members", []string{}, "Group members")
}

func promptCourseName() string {
	fmt.Print("❓ Please enter the course's name: ")
	course := ""
	fmt.Scanln("%s", course)
	return course
}

func promptGroupName() string {
	fmt.Print("❓ Please enter a group name (or leave empty to use nothing): ")
	group := ""
	fmt.Scanln("%s", group)
	return group
}

func promptMembers() []string {
	members = make([]string, 0)
	for {
		fmt.Print("❓ Please enter a group member's name followed by their student ID (e.g. Max Mustermann,123456), or press 'q' to move on: ")
		rawMember := ""
		fmt.Scanln("%s", rawMember)
		if rawMember == "q" {
			break
		}

		members = append(members, strings.TrimSpace(rawMember))
	}
	return members
}
