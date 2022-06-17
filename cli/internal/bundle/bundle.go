package bundle

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
)

// BundlingContext contains all fields passed down to the template renderer for the filename of a bundle
type BundlingContext struct {
	// ID is the assignment id, e.g. "01"
	ID string `json:"id"`
	// Members contains group member information, possibly included in the filename
	Members []config.GroupMember `json:"members"`
	Format  BundlerBackend
}

type BundlerBackend string

const (
	BundlerBackendZip BundlerBackend = "zip"
	BundlerBackendTar BundlerBackend = "tar"
)

var defaultArchiveNameTemplate string = "assignment-{{.id}}.{{.format}}"

type BundlerOptions struct {
	Backend  BundlerBackend
	Template *string
	Quiet    bool
	Data     map[string]interface{}
}

type BundlerResponse struct {
	Filename string
}

type BundlerContext struct {
	context.AppContext
	BundlerOptions
}

func New(ctx *context.AppContext, options *BundlerOptions) (*BundlerContext, error) {
	bundler := &BundlerContext{
		BundlerOptions: *options,
		AppContext:     *ctx,
	}

	return bundler, nil
}

func (b *BundlerContext) Make() (*BundlerResponse, error) {

	archiveName, err := b.generateArchiveName()

	if err != nil {
		return nil, err
	}

	return &BundlerResponse{
		Filename: archiveName,
	}, nil
}

// generateArchiveName executes the template with the data given in the config file
func (b *BundlerContext) generateArchiveName() (string, error) {
	tpl := b.Template
	if tpl == nil {
		tpl = &defaultArchiveNameTemplate
	}

	data := b.deriveOrOverrideFormat(b.AppContext.Configuration.Spec.BundleOptions.Data)

	tmpl := template.Must(template.New("bundleName").Funcs(sprig.TxtFuncMap()).Parse(*tpl))
	var output bytes.Buffer

	err := tmpl.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// deriveOrOverrideFormat takes an arbitrary map of data bindings, looks for a predefined field for
// the archive's format and otherwise derives the format from the chosen bundler backend
func (b *BundlerContext) deriveOrOverrideFormat(data map[string]interface{}) map[string]interface{} {
	if _, ok := data["format"]; !ok {
		// no format override in map, derive from chosen backend
		if b.BundlerOptions.Backend == BundlerBackendTar {
			data["format"] = "tar.gz"
			return data
		}
		data["format"] = "zip"
		return data
	}
	return data
}
