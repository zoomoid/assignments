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
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/bundle"
	"github.com/zoomoid/assignments/v1/internal/context"
	"github.com/zoomoid/assignments/v1/internal/util"
)

var (
	bundleLongDescription = dedent.Dedent(`
		Bundling compiles all files relevant for an assignment into an archive
		format. The backend defaults to zip, but can be set to tarball by
		passing the --tar flag. If you want to use tar and gzip, use --gzip.

		By default, every bundle includes at least the assignment's PDF from
		the ./dist/ directory. If you want to add further files or directories
		see the list .spec.bundle.include in your configuration file. It lets
		you specify files explicitly, or a glob pattern for multiple files,
		e.g. "code/*" or "figures/*.pdf". It is meant to complement the list 
		of directories to create when using the generate command. The bundle will
		preserve the structure of the files included, and will have the PDF
		located at the archive's root.

		You can customize how the filename for the archive is generated. For this,
		you can set .spec.bundle.template to be an arbitrary Golang text template
		(including the use of sprig text functions). Just note that this is 
		limited by what file paths are supported by your operating system, so 
		don't get too crazy. The map in .spec.bundle.data is passed down to
		the template's execution for data binding. 

		The default archive template is "assignment-{{._id}}.{{._format}}". Note 
		the _id field: this is internally augmented from the command's arguments
		or the configuration's status field (or, in case of usage of --all, all
		available assignments in the repository). "format" is derived from the 
		selected backend's common file extension, but respects overrides from
		the map at .spec.bundle.data, so you can also pick your own file extension
		without overriding the entire template.
	`)
)

type bundleData struct {
	all   bool
	force bool
	tar   bool
	gzip  bool
}

func newBundleData() *bundleData {
	return &bundleData{
		all:   false,
		force: false,
		tar:   false,
		gzip:  false,
	}
}

func NewBundleCommand(ctx *context.AppContext, data *bundleData) *cobra.Command {
	if data == nil {
		data = newBundleData()
	}

	bundleCommand := &cobra.Command{
		Use:   "bundle",
		Short: "Bundles an assignment with all additional files inside the assignment's directory",
		Long:  bundleLongDescription,
		PreRun: func(cmd *cobra.Command, args []string) {
			err := ctx.Read()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read config file")
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			defer ctx.Write()
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return getAssignmentsFromRoot(toComplete, ctx.Root), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			assignmentNo := ctx.Configuration.Status.Assignment
			if len(args) != 0 {
				assignmentArg := args[0]
				// attempt to remove prefix from arguments. This is relevant when using the autocompletion
				assignmentArg = strings.TrimPrefix(assignmentArg, "assignment-")
				// attempt parsing numerically
				i, err := strconv.Atoi(assignmentArg)
				if err == nil {
					assignmentNo = uint32(i)
				}
			}

			if data.all && len(args) > 0 {
				return errors.New("cannot use --all flag with specific assignment")
			}

			if data.gzip && !data.tar {
				return errors.New("cannot use --gzip without --tar")
			}

			backend := bundle.BundlerBackendZip
			if data.tar {
				if data.gzip {
					backend = bundle.BundlerBackendTarGzip
				} else {
					backend = bundle.BundlerBackendTar
				}
			}

			var template string
			if ctx.Configuration.Spec.BundleOptions != nil && ctx.Configuration.Spec.BundleOptions.Template != "" {
				template = ctx.Configuration.Spec.BundleOptions.Template
			}

			templateBindings := make(map[string]interface{})
			if ctx.Configuration.Spec.BundleOptions != nil && ctx.Configuration.Spec.BundleOptions.Data != nil {
				templateBindings = ctx.Configuration.Spec.BundleOptions.Data
			}
			templateBindings["_id"] = util.AddLeadingZero(assignmentNo)

			bundleRuns := []string{}

			if !data.all {
				assignment := fmt.Sprintf("assignment-%s.pdf", util.AddLeadingZero(assignmentNo))
				bundleRuns = append(bundleRuns, assignment)
			} else {
				assignments, err := filepath.Glob(filepath.Join(ctx.Root, "dist", "assignment-*.pdf"))
				if err != nil {
					return err
				}
				for _, assignment := range assignments {
					bundleRuns = append(bundleRuns, filepath.Base(assignment))
				}
			}

			includes := []string{}
			if ctx.Configuration.Spec.BundleOptions != nil && ctx.Configuration.Spec.BundleOptions.Include != nil {
				includes = ctx.Configuration.Spec.BundleOptions.Include
			}

			for _, file := range bundleRuns {
				opts := &bundle.BundlerOptions{
					Backend:  backend,
					Template: template,
					Data:     templateBindings,
					Target:   filepath.Base(file),
					Includes: includes,
					Force:    data.force,
				}
				bundler, err := bundle.New(ctx, opts)
				if err != nil {
					if errors.Is(err, bundle.ErrArchiveExists) {
						// only skip the current bundling, continue with other runs
						log.Warn().Msgf("Archive %s already exists and --force is not specified, skipping...", bundler.ArchiveName())
						break
					}
					return err
				}

				if err := bundler.Bundle(); err != nil {
					return err
				}

				archiveName := bundler.ArchiveName()

				log.Info().Msgf("Finished bundling assignment to %s in ./dist/", archiveName)
			}
			return nil
		},
	}

	addBundleFlags(bundleCommand.PersistentFlags(), data)
	addBundleFlagsCompletion(bundleCommand)

	return bundleCommand
}

func addBundleFlags(flags *pflag.FlagSet, data *bundleData) {
	flags.BoolVarP(&data.all, options.All, options.AllShort, false, "Bundle all assignments")
	flags.BoolVarP(&data.force, options.Force, options.ForceShort, false, "Override any existing archives with the same name")
	flags.BoolVar(&data.tar, options.Tar, false, "Use tar as a backend for archive bundling")
	flags.BoolVar(&data.gzip, options.Gzip, false, "Use gzip to encode the archive. Requires --tar to be specified as well")
}

func addBundleFlagsCompletion(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc(options.All, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.Force, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.Tar, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.Gzip, cobra.NoFileCompletions)
}
