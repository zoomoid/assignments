package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/lithammer/dedent"
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

		The default archive template is "assignment-{{._id}}.{{.format}}". Note 
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

	err := ctx.Read()
	if err != nil {
		ctx.Logger.Fatalf("Failed to read config file", err)
	}

	bundleCommand := &cobra.Command{
		Use:   "bundle",
		Short: "Bundles an assignment with all additional files inside the assignment's directory",
		Long:  bundleLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			assignmentNo := ctx.Configuration.Status.Assignment
			if len(args) != 0 {
				i, err := strconv.Atoi(args[0])
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

			var template *string = nil
			if ctx.Configuration.Spec.BuildOptions != nil && ctx.Configuration.Spec.BundleOptions.Template == "" {
				template = &ctx.Configuration.Spec.BundleOptions.Template
			}

			templateBindings := ctx.Configuration.Spec.BundleOptions.Data
			templateBindings["_id"] = assignmentNo

			files := []string{}

			if !data.all {
				assignment := fmt.Sprintf("assignments-%s.pdf", util.AddLeadingZero(assignmentNo))
				files = append(files, assignment)
			} else {
				assignments, err := filepath.Glob(filepath.Join(ctx.Root, "dist", "assignment-*.pdf"))
				if err != nil {
					return err
				}
				for _, assignment := range assignments {
					files = append(files, filepath.Base(assignment))
				}
			}

			for _, file := range files {
				opts := &bundle.BundlerOptions{
					Backend:  backend,
					Template: template,
					Data:     templateBindings,
					Target:   filepath.Base(file),
					Includes: ctx.Configuration.Spec.BundleOptions.Include,
				}
				bundler, err := bundle.New(ctx, opts)
				if err != nil {
					return err
				}
				fn, err := bundler.Make()
				if err != nil {
					return err
				}

				ctx.Logger.Infof("Finished bundling %s to ./dist/", fn)
			}

			return nil
		},
	}

	addBundleFlags(bundleCommand.Flags(), data)

	return bundleCommand
}

func addBundleFlags(flags *pflag.FlagSet, data *bundleData) {
	flags.BoolVarP(&data.all, options.All, options.AllShort, false, "Bundle all assignments")
	flags.BoolVarP(&data.force, options.Force, options.ForceShort, false, "Override any existing archives with the same name")
	flags.BoolVar(&data.tar, options.Tar, false, "Use tar as a backend for archive bundling")
	flags.BoolVar(&data.gzip, options.Gzip, false, "Use tar and gzip as backend for archive bundling")
}
