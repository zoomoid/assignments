/*
Copyright 2022 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

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
		required for the job running Gitlab's release-cli or Github's CLI, respectively.
		
		You can run this command outside of CI pipelines, note however that it is highly
		dependant on the ENV variables being available in the runner context. You will
		have to provide either $CI_COMMIT_TAG with $CI_JOB_ID and $CI_PROJECT_URL or 
		$GITHUB_REF_NAME for the command to output correct .env files.
	`)

	ciBootstrapLongDescription = dedent.Dedent(`
		Run this command to quickly template CI files for the supported SCM providers,
		namely Gitlab and Github. Afterwards, you can customize them to your liking.

		To learn more about the CI integration, see the documentation at 
		https://github.com/zoomoid/assignments/blob/main/ci/README.md
	`)
)

type SCMProvider string

const (
	GitlabSCM SCMProvider = "gitlab"
	GithubSCM SCMProvider = "github"
)

func NewCiCommand(ctx *context.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ci",
		Short: "Manage CI integration for assignmentctl",
		Long:  "The command is not meant to be run on its own",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := ctx.Read()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read config file")
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			defer ctx.Write()
		},
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
			fmt.Sprintf("%s\t%s", string(GithubSCM), "Creates a .env file to use in a release job in Github Action"),
			fmt.Sprintf("%s\t%s", string(GitlabSCM), "Creates a .env file to use in a release job in Gitlab CI"),
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
	cmd.RegisterFlagCompletionFunc(options.File, cobra.NoFileCompletions)

	return cmd
}

func NewCiBootstrapCommand(ctx *context.AppContext) *cobra.Command {
	file := ""

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Create a template CI workflow file for the selected SCM provider",
		Long:  ciBootstrapLongDescription,
		Args:  cobra.ExactValidArgs(1),
		ValidArgs: []string{
			fmt.Sprintf("%s\t%s", string(GithubSCM), "Creates a Github Actions YAML file"),
			fmt.Sprintf("%s\t%s", string(GitlabSCM), "Creates a Gitlab CI .gitlab-ci.yml file"),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			out, isStdout := ci.OpenOrFallbackToStdout(file)
			if !isStdout {
				defer out.Close()
			}

			// with validators configured this can be assumed to be the only argument
			t := SCMProvider(args[0])
			if t == GithubSCM {
				out.WriteString(strings.TrimSpace(ci.GithubActionTemplate))
				return nil
			}

			if t == GitlabSCM {
				out.WriteString(strings.TrimSpace(ci.GitlabCITemplate))
				return nil
			}
			return fmt.Errorf("%s is not a supported SCM provider", args[0])
		},
	}

	cmd.PersistentFlags().StringVarP(&file, options.File, options.FileShort, "", "Write the template directly to a file instead of Stdout")
	cmd.RegisterFlagCompletionFunc(options.File, cobra.NoFileCompletions)

	return cmd
}
