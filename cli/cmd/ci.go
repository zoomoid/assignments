package cmd

import (
	"path/filepath"

	"github.com/lithammer/dedent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/context"
	"github.com/zoomoid/assignments/v1/internal/release"
)

var (
	ciLongDescription = dedent.Dedent(`
		The command is meant for usage inside CI pipelines to create release objects 
		for Gitlab and exports several environment variables in a file that are 
		required for the job running the release-cli.
	`)
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
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Create a release object for the specified git provider, either Github or Gitlab",
	}

	cmd.AddCommand(NewCiReleaseGitlabCommand(ctx))

	return cmd
}

func NewCiReleaseGitlabCommand(ctx *context.AppContext) *cobra.Command {
	file := ""

	cmd := &cobra.Command{
		Use:   "gitlab",
		Short: "Creates an ENV file to source variables for the Gitlab release-cli from",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, isStdout := release.OpenOrFallbackToStdout(file)
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

			o, err := release.TemplateGitlabCIEnvFile(artifactsDirectory, archiveTemplate, data)
			if err != nil {
				return err
			}
			out.Write(o.Bytes())

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&file, options.File, options.FileShort, "", "Write the template directly to a file instead of Stdout")

	return cmd
}

func NewCiBootstrapCommand(ctx *context.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Creates a template CI file for the specified git provider, either Github or Gitlab",
	}

	cmd.AddCommand(NewCiBootstrapGitlabCommand(ctx))

	return cmd
}

func NewCiBootstrapGitlabCommand(ctx *context.AppContext) *cobra.Command {
	file := ""

	cmd := &cobra.Command{
		Use:   "gitlab",
		Short: "Create a template CI file for Gitlab CI",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, isStdout := release.OpenOrFallbackToStdout(file)
			if !isStdout {
				defer out.Close()
			}
			out.WriteString(release.GitlabCITemplate)
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&file, options.File, options.FileShort, "", "Write the template directly to a file instead of Stdout")

	return cmd
}
