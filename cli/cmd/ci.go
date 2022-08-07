package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/lithammer/dedent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/ci"
	"github.com/zoomoid/assignments/v1/internal/context"
)

var (
	ciReleaseLongDescription = dedent.Dedent(`
		The command is meant for usage inside CI pipelines to create release objects 
		for Gitlab and exports several environment variables in a file that are 
		required for the job running the release-cli.
	`)
)

type SCMProvider string

const (
	GitlabSCM SCMProvider = "gitlab"
	GithubSCM SCMProvider = "github"
)

func NewCiCommand(ctx *context.AppContext) *cobra.Command {
	err := ctx.Read()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to read config file")
	}
	defer ctx.Write()

	cmd := &cobra.Command{
		Use:   "ci",
		Short: "Export environment variables and context to a CI job",
	}

	cmd.AddCommand(NewCiReleaseCommand(ctx))
	cmd.AddCommand(NewCiBootstrapCommand(ctx))

	return cmd
}

func NewCiReleaseCommand(ctx *context.AppContext) *cobra.Command {
	file := ""

	cmd := &cobra.Command{
		Use:   "release",
		Short: "Creates an ENV file to source variables for the selected SCM provider",
		Long:  ciReleaseLongDescription,
		Args:  cobra.ExactValidArgs(1),
		ValidArgs: []string{
			string(GithubSCM),
			string(GitlabSCM),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			out, isStdout := ci.OpenOrFallbackToStdout(file)
			if !isStdout {
				defer out.Close()
			}

			artifactsDirectory := filepath.Join(ctx.Root, "dist")
			archiveTemplate := ""
			data := make(map[string]interface{})
			if b := ctx.Configuration.Spec.BundleOptions; b != nil {
				archiveTemplate = b.Template
			}
			if d := ctx.Configuration.Spec.BundleOptions.Data; d != nil {
				data = d
			}

			// with validators configured this can be assumed to be the only argument
			t := SCMProvider(args[0])
			if t == GithubSCM {
				o, err := ci.TemplateGithubActionsEnvFile(artifactsDirectory, archiveTemplate, data)
				if err != nil {
					return err
				}
				out.Write(o.Bytes())
				return nil
			}

			if t == GitlabSCM {
				o, err := ci.TemplateGitlabCIEnvFile(artifactsDirectory, archiveTemplate, data)
				if err != nil {
					return err
				}
				out.Write(o.Bytes())
				return nil
			}

			return fmt.Errorf("%s is not a supported SCM provider", args[0])
		},
	}

	cmd.PersistentFlags().StringVarP(&file, options.File, options.FileShort, "", "Write the template directly to a file instead of Stdout")

	return cmd
}

func NewCiBootstrapCommand(ctx *context.AppContext) *cobra.Command {
	file := ""

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Create a template CI workflow file for either Github or Gitlab",
		Args:  cobra.ExactValidArgs(1),
		ValidArgs: []string{
			string(GithubSCM),
			string(GitlabSCM),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			out, isStdout := ci.OpenOrFallbackToStdout(file)
			if !isStdout {
				defer out.Close()
			}

			// with validators configured this can be assumed to be the only argument
			t := SCMProvider(args[0])
			if t == GithubSCM {
				out.WriteString(ci.GithubActionTemplate)
				return nil
			}

			if t == GitlabSCM {
				out.WriteString(ci.GitlabCITemplate)
				return nil
			}
			return fmt.Errorf("%s is not a supported SCM provider", args[0])
		},
	}

	cmd.PersistentFlags().StringVarP(&file, options.File, options.FileShort, "", "Write the template directly to a file instead of Stdout")

	return cmd
}
