package bundle

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
	"github.com/zoomoid/assignments/v1/internal/context"
	"github.com/zoomoid/assignments/v1/internal/util"
)

// BundlerBackend is a specific string type for picking backends, and consequently file endings
type BundlerBackend string

const (
	// BundlerBackendZip selects archive/zip and the corresponding bundler as backend
	BundlerBackendZip BundlerBackend = "zip"
	// BunlderBackendTar selects archive/tar and the corresponding bundler as backend
	BundlerBackendTar BundlerBackend = "tar"
	// BundlerBackendTarGzip combines archive/tar and encoding/gzip in the corresponding bundler as backend
	BundlerBackendTarGzip BundlerBackend = "tar.gz"
)

var (
	// ErrAchiveExists is a static error that indicates an archive already existing without explitly truncating it
	ErrArchiveExists error = errors.New("archive already exists")
	// default archive template. This contains the fields _id and format, which are automatically aliased
	// into the map that contains the data bindings, thus are ALWAYS available
	defaultArchiveNameTemplate string = "assignment-{{._id}}.{{.format}}"
)

// Bundler interface all backends should implement
type Bundler interface {
	// AddAssignment adds the main assignment PDF at the root of the archive
	AddAssignment() error
	// AddAuxilliaryFiles adds all files from directories defined in the bundling spec to the archive
	AddAuxilliaryFiles() error
	// Close finishes the archive creation by closing all remaining writers
	Close() error
	// Type returns the BundlerBackend type of this particular instance
	Type() BundlerBackend
}

// BundlerOptions captures shared bundler options
type BundlerOptions struct {
	// Backend selects the bundling backend to use
	Backend BundlerBackend
	// Template is a string template to be used for the filename
	Template string
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
	// Force indicates truncating any existing archives with the same name and creating it from scratch
	Force bool
}

type BundlerContext struct {
	context.AppContext
	// BundlerOptions are all fields passed into the New constructor for a bundler
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
	// archiveName is the name of the archive created when executing the template
	archiveName string
}

// New makes a new bundling context from the context and the options passed as parameters
func New(ctx *context.AppContext, options *BundlerOptions) (*BundlerContext, error) {
	artifactsDirectory := filepath.Join(ctx.Root, "dist")
	sourceDirectory := strings.ReplaceAll(options.Target, ".pdf", "")
	base := filepath.Join(ctx.Root, sourceDirectory)

	additionalFiles, err := additionalFiles(base, options.Includes)
	if err != nil {
		return nil, err
	}

	data := options.Data
	if data == nil {
		data = make(map[string]interface{})
	}
	id, err := util.AssignmentNumberFromFilename(filepath.Base(options.Target))
	if err != nil {
		return nil, err
	}
	data["_id"] = id
	if _, ok := data["format"]; !ok {
		data["format"] = format(options.Backend)
	}

	archiveName, err := makeArchiveName(options.Template, data, options.Backend)
	if err != nil {
		return nil, err
	}

	bundlerCtx := ctx.Clone()

	bundler := &BundlerContext{
		BundlerOptions:     *options,
		AppContext:         *bundlerCtx,
		files:              additionalFiles,
		artifactsDirectory: artifactsDirectory,
		sourceDirectory:    sourceDirectory,
		base:               base, // this is always the same as sourceDirectory, maybe save on this field?
		archiveName:        archiveName,
	}

	if !options.Force && bundler.ArchiveExists() {
		return bundler, ErrArchiveExists
	}

	return bundler, nil
}

// Bundle runs the bundling action by picking a bundle implementor from the selected backend
// Returns the archive's filename when successful, otherwise an error and the empty string
func (b *BundlerContext) Bundle() error {
	bundler, err := b.makeBundler()
	if err != nil {
		return err
	}

	if err := bundler.AddAssignment(); err != nil {
		return err
	}
	if err := bundler.AddAuxilliaryFiles(); err != nil {
		return err
	}
	if err := bundler.Close(); err != nil {
		return err
	}

	return nil
}

// ArchiveName returns the context's archive name field to other modules
func (b *BundlerContext) ArchiveName() string {
	return b.archiveName
}

// ArchiveExists checks if the archive that is created by the bundler already exists.
// Returns true if `os.Stat` is successful. Returns false otherwise, *even if the
// error returned by `os.Stat` is not os.ErrNotExist*.
func (b *BundlerContext) ArchiveExists() bool {
	_, err := os.Stat(filepath.Join(b.artifactsDirectory, b.archiveName))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Warn().Msgf("failed to stat %s with error other than ErrNotExist, %v", filepath.Join(b.artifactsDirectory, b.archiveName), err)
		}
		return false
	}
	return true
}

// makeBundler internally differentiates between bundler implementations chosen
// by the backend, and returns an instance of that bundler.
//
// If a backend is selected that isn't explicitly supported, makeBundler returns
// an error containing the name of the chosen backend.
func (b *BundlerContext) makeBundler() (Bundler, error) {
	switch b.Backend {
	case BundlerBackendTar:
		bundler, err := NewTarBundler(&b.AppContext, b.files, &TarBundlerOptions{
			ArchiveName:        b.archiveName,
			SourceDirectory:    b.base,
			ArtifactsDirectory: b.artifactsDirectory,
		})
		return bundler, err
	case BundlerBackendTarGzip:
		bundler, err := NewGzipBundler(&b.AppContext, b.files, &GzipBundlerOptions{
			ArchiveName:        b.archiveName,
			SourceDirectory:    b.base,
			ArtifactsDirectory: b.artifactsDirectory,
		})
		return bundler, err
	case BundlerBackendZip:
		bundler, err := NewZipBundler(&b.AppContext, b.files, &ZipBundlerOptions{
			ArchiveName:        b.archiveName,
			ArtifactsDirectory: b.artifactsDirectory,
			SourceDirectory:    b.base,
		})
		return bundler, err
	}
	return nil, fmt.Errorf("backend %s is not supported", b.Backend)
}

// additionalFiles takes a slice of paths or glob patterns and a source directory,
// and executes the glob pattern in that source directory. It returns a slice of
// all files that matched the glob pattern and all files directly matched. If an
// element in includes is a directory *without* a glob pattern, any children of that
// directory are ignored. You will have to use a glob pattern to include all children
// of a directory.
//
// If executing the glob pattern fails, additionalFiles returns nil and an error containing
// the pattern that failed to glob.
func additionalFiles(sourceDirectory string, includes []string) ([]string, error) {
	additionalFiles := make([]string, 0)
	for _, f := range includes {
		if _, err := os.Stat(f); errors.Is(err, fs.ErrNotExist) {
			p := filepath.Join(sourceDirectory, f)
			matches, err := filepath.Glob(p)
			if err == nil && len(matches) == 0 {
				// return nil, fmt.Errorf("the path %q does not exist", pattern)
				return []string{}, nil
			}
			if err == filepath.ErrBadPattern {
				return nil, fmt.Errorf("patterns %q is not valid: %w", f, err)
			}
			additionalFiles = append(additionalFiles, matches...)
			continue
		}
		additionalFiles = append(additionalFiles, f)
	}
	return additionalFiles, nil
}

// makeArchiveName executes the template with the data given in the config file
// Returns the archive's filename when successfully executed the template, otherwise
// returns the occurred error and an empty string
func makeArchiveName(tpl string, data map[string]interface{}, backend BundlerBackend) (string, error) {
	if tpl == "" {
		tpl = defaultArchiveNameTemplate
	}

	tmpl := template.Must(template.New("bundleName").Funcs(sprig.TxtFuncMap()).Parse(tpl))
	var output bytes.Buffer

	err := tmpl.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// format takes an arbitrary map of data bindings, looks for a predefined field for
// the archive's format and otherwise derives the format from the chosen bundler backend
func format(backend BundlerBackend) BundlerBackend {
	// no format override in map, derive from chosen backend
	switch backend {
	case BundlerBackendTar,
		BundlerBackendZip,
		BundlerBackendTarGzip:
		return backend
	default:
		return ""
	}
}
