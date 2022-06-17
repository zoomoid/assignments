package cmd

import (
	"errors"
	"strconv"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/assignments/v1/cmd/options"
	"github.com/zoomoid/assignments/v1/internal/bundle"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
)

var (
	bundleLongDescription = dedent.Dedent(`
		Bundling compiles all files relevant for an assignment into an archive
		format. The backend defaults to zip, but can be set to tarball by
		passing the --tar flag. 
	`)
)

type bundleData struct {
	all   bool
	quiet bool
	force bool
	tar   bool
}

func newBundleData() *bundleData {
	return &bundleData{
		all:   false,
		quiet: false,
		force: false,
		tar:   false,
	}
}

func NewBundleCommand(ctx *context.AppContext, data *bundleData) *cobra.Command {

	if data == nil {
		data = newBundleData()
	}

	cfg, err := config.ReadConfigMap()

	if err != nil {
		ctx.Logger.Fatalf("Failed to read config file, %v", err)
	}

	ctx.Configuration = cfg

	bundleCommand := &cobra.Command{
		Use:   "bundle",
		Short: "Bundles an assignment with all additional files inside the assignment's directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			assignmentNo := cfg.Status.Assignment
			if len(args) != 0 {
				i, err := strconv.Atoi(args[0])
				if err == nil {
					assignmentNo = uint32(i)
				}
			}

			if data.all && len(args) > 0 {
				return errors.New("cannot use --all flag with specific assignment")
			}

			backend := bundle.BundlerBackendZip
			if data.tar {
				backend = bundle.BundlerBackendTar
			}

			var template *string = nil
			if ctx.Configuration.Spec.BuildOptions != nil && ctx.Configuration.Spec.BundleOptions.Template == "" {
				template = &ctx.Configuration.Spec.BundleOptions.Template
			}

			templateBindings := ctx.Configuration.Spec.BundleOptions.Data
			templateBindings["_id"] = assignmentNo

			if !data.all {
				opts := &bundle.BundlerOptions{
					Backend:  backend,
					Template: template,
					Quiet:    data.quiet,
					Data:     templateBindings,
				}

				bundle.New(ctx, opts)
			} else {

			}

			return nil
		},
	}

	return bundleCommand
}

func addBundleFlags(flags *pflag.FlagSet, data *bundleData) {
	flags.BoolVarP(&data.all, options.All, options.AllShort, false, "Bundle all assignments")
	flags.BoolVarP(&data.quiet, options.Quiet, options.QuietShort, false, "Suppress output from bundling")
	flags.BoolVarP(&data.force, options.Force, options.ForceShort, false, "Override any existing archives with the same name")
	flags.BoolVar(&data.tar, options.Tar, false, "Use tar as a backend for archive bundling")
}
