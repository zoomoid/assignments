package bundle

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zoomoid/assignments/v1/internal/context"
)

type ZipBundlerOptions struct {
	// ArchiveName created in Bundler.New()
	ArchiveName string
	// SourceDirectory is the directory where the additional files and
	// the assignment's source originate
	SourceDirectory string
	// Artifacts directory is the directory where the PDF originates from
	ArtifactsDirectory string
}

type ZipBundler struct {
	// context.AppContext is a reference to a cloned AppContext
	*context.AppContext
	// embedded struct of the options passed to NewGzipBundler
	*ZipBundlerOptions
	// backend is the selected backend, here, always BundlerBackendTarGzip, for aligning file endings
	backend BundlerBackend
	// writer is a pointer to the writer context for tar files. It is configured to pipe its
	// output to the gzip writer for compression
	writer *zip.Writer
	// archive is the file descriptor opened to write the gzipped tarball to. Calling Bundler.Close()
	// will close it, as well as the writers linked to it
	archive *os.File
	// sourceDirectory contains the directory from which to include additional files from.
	// It is the base for the Include patterns in Configuration.Spec.BundleOptions.
	sourceDirectory string
	// files contains paths to all auxilliary files to be added to an archive
	files []string
}

var _ Bundler = &ZipBundler{}

// NewZipBundler makes a new zip bundler that uses the archive/zip module.
// It returns a ZipBundler that implements the Bundler interface.
//
// If the archive file descriptor cannot be created, NewZipBundler returns an error.
func NewZipBundler(ctx *context.AppContext, files []string, options *ZipBundlerOptions) (*ZipBundler, error) {
	archive, err := os.Create(filepath.Join(options.ArtifactsDirectory, options.ArchiveName))
	if err != nil {
		return nil, err
	}

	zw := zip.NewWriter(archive)

	bundler := &ZipBundler{
		backend:           BundlerBackendZip,
		AppContext:        ctx,
		ZipBundlerOptions: options,
		writer:            zw,
		files:             files,
		sourceDirectory:   options.SourceDirectory,
	}

	return bundler, nil
}

// Close closes the remaining file descriptors for the zip archive
func (b *ZipBundler) Close() error {
	defer b.archive.Close()
	return b.writer.Close()
}

// Type returns the static type of the bundler
func (b *ZipBundler) Type() BundlerBackend {
	return BundlerBackendZip
}

// AddAssignmentToArchive implements writing the assignment's PDF to the zip archive
func (b *ZipBundler) AddAssignment() error {
	if b.writer == nil {
		return errors.New("writer not created yet")
	}
	assignmentPdfName := fmt.Sprintf("%s.pdf", filepath.Base(b.SourceDirectory))
	src, err := os.Open(filepath.Join(b.ArtifactsDirectory, assignmentPdfName))
	if err != nil {
		return err
	}
	dst, err := b.writer.Create(assignmentPdfName)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return err
}

// AddAuxilliaryFiles adds all additional files added to the bundler
// instance.
//
// If a file failes to be added to the archive, AddAuxilliaryFiles will
// exit prematurely, without attempting to add any of the subsequent files,
// and return the error.
func (b *ZipBundler) AddAuxilliaryFiles() error {
	for _, f := range b.files {
		if err := b.addAuxilliaryFile(f); err != nil {
			return err
		}
	}
	return nil
}

// addAuxilliaryFile opens a file descriptor for the file and
// writes it to the zip archive file
func (b *ZipBundler) addAuxilliaryFile(filename string) error {
	if b.writer == nil {
		return errors.New("writer not created yet")
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	// remove root-level prefix from filename to preserve assignment directory structure
	f, err := b.writer.Create(filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	return nil
}
