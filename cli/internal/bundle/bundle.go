package bundle

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
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

// BundlerBackend is a specific string type for
type BundlerBackend string

const (
	// BundlerBackendZip selects archive/zip and the corresponding bundler as backend
	BundlerBackendZip BundlerBackend = "zip"
	// BunlderBackendTar selects archive/tar and the corresponding bundler as backend
	BundlerBackendTar BundlerBackend = "tar"
	// BundlerBackendTarGzip combines archive/tar and encoding/gzip in the corresponding bundler as backend
	BundlerBackendTarGzip BundlerBackend = "tar.gz"
)

var defaultArchiveNameTemplate string = "assignment-{{._id}}.{{.format}}"

// Bundler interface all backends should implement
type Bundler interface {
	// AddAssignment adds the main assignment PDF at the root of the archive
	AddAssignment() error
	// AddAuxilliaryFile adds all files from directories defined in the bundling spec to the archive
	AddAuxilliaryFile(filename string) error
	// Close finishes the archive creation by closing all remaining writers
	Close() error
}

// BundlerOptions captures shared bundler options
type BundlerOptions struct {
	// Backend selects the bundling backend to use
	Backend BundlerBackend
	// Template is a string pointer to the template to be used for the filename
	Template *string
	// Data contains any data bindings required for the template execution
	// Users should ensure that the data matches the template, otherwise template
	// execution will fail ungracefully
	Data map[string]interface{}
	// Target is the basename of the assignment pdf to bundle.
	// Used to derive paths to other relevant files
	Target string
	// Includes are the directories to additionally be included in the archive
	// These are defined in the configuration file and should be relative to
	// each assignment's root
	Includes []string
}

type BundlerContext struct {
	context.AppContext
	BundlerOptions
	// files contain all the files additionally to be included
	files []string
	// sourceDirectory is the directory used for any defined additional files
	sourceDirectory string
	// artifactsDirectory is the ./dist/ directory from which the PDF originates
	artifactsDirectory string
	// base is a string of the form assignment-<no> required for structural assumptions
	// about the directory structure
	base string
}

// New makes a new bundling context from the context and the options passed as parameters
func New(ctx *context.AppContext, options *BundlerOptions) (*BundlerContext, error) {
	artifactsDirectory := filepath.Join(ctx.Cwd, "dist")
	sourceDirectory := strings.ReplaceAll(options.Target, ".pdf", "")
	base := filepath.Join(ctx.Cwd, sourceDirectory)

	additionalFiles := make([]string, 0)
	for _, f := range options.Includes {
		if strings.Contains(f, "*") {
			// f is a glob pattern
			files, err := filepath.Glob(filepath.Join(sourceDirectory, f))
			if err != nil {
				return nil, fmt.Errorf("malformed glob pattern in %s", f)
			}
			// append all globbed files to the list of additional files
			additionalFiles = append(additionalFiles, files...)
		} else {
			if !strings.HasSuffix(f, string(filepath.Separator)) {
				// f likely is a file, append it
				additionalFiles = append(additionalFiles, filepath.Join(sourceDirectory, f))
			} // otherwise f is a directory, which, if it does not contain a glob pattern, will be ignored
		}
	}

	bundler := &BundlerContext{
		BundlerOptions:     *options,
		AppContext:         *ctx,
		files:              additionalFiles,
		artifactsDirectory: artifactsDirectory,
		sourceDirectory:    sourceDirectory,
		base:               base, // this is always the same as sourceDirectory, maybe save on this field?
	}

	return bundler, nil
}

// Make runs the bundling action by picking a bundle implementor from the selected backend
// Returns the archive's filename when successful, otherwise an error and the empty string
func (b *BundlerContext) Make() (string, error) {
	archiveName, err := GenerateArchiveName(b.Template, b.Data, b.Backend)
	if err != nil {
		return "", err
	}

	if b.Backend == BundlerBackendTar {
		bundler, err := NewTarBundler(&b.AppContext, &TarBundlerOptions{
			ArchiveName:        archiveName,
			AssignmentBase:     b.base,
			ArtifactsDirectory: b.artifactsDirectory,
		})
		if err != nil {
			return "", err
		}

		if err := bundler.AddAssignment(); err != nil {
			return "", err
		}
		for _, file := range b.files {
			if err := bundler.AddAuxilliaryFile(file); err != nil {
				return "", err
			}
		}
		if err := bundler.Close(); err != nil {
			return "", err
		}

		return archiveName, nil
	}
	if b.Backend == BundlerBackendTarGzip {
		return archiveName, nil
	}
	if b.Backend == BundlerBackendZip {
		bundler, err := NewZipBundler(&b.AppContext, &ZipBundlerOptions{
			ArtifactsDirectory: b.artifactsDirectory,
			ArchiveName:        archiveName,
			AssignmentBase:     b.base,
		})
		if err != nil {
			return "", err
		}

		if err := bundler.AddAssignment(); err != nil {
			return "", err
		}
		for _, file := range b.files {
			if err := bundler.AddAuxilliaryFile(file); err != nil {
				return "", err
			}
		}

		if err := bundler.Close(); err != nil {
			return "", err
		}

		return archiveName, nil
	}
	return "", fmt.Errorf("bundling backend %s not supported", b.Backend)
}

// GenerateArchiveName executes the template with the data given in the config file
// Returns the archive's filename when successfully executed the template, otherwise
// returns the occurred error and
func GenerateArchiveName(tpl *string, data map[string]interface{}, backend BundlerBackend) (string, error) {
	if tpl == nil {
		tpl = &defaultArchiveNameTemplate
	}

	data = deriveOrOverrideFormat(data, backend)

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
func deriveOrOverrideFormat(data map[string]interface{}, backend BundlerBackend) map[string]interface{} {
	if _, ok := data["format"]; !ok {
		// no format override in map, derive from chosen backend
		switch backend {
		case BundlerBackendTar:
		case BundlerBackendZip:
		case BundlerBackendTarGzip:
			data["format"] = backend
		default:
			data["format"] = ""
		}
		return data
	}
	return data
}
