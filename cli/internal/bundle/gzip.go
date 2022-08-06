package bundle

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zoomoid/assignments/v1/internal/context"
)

type GzipBundlerOptions struct {
	// ArchiveName created in Bundler.New()
	ArchiveName string
	// SourceDirectory is the directory where the additional files and
	// the assignment's source originate
	SourceDirectory string
	// Artifacts directory is the directory where the PDF originates from
	ArtifactsDirectory string
}

type GzipBundler struct {
	// context.AppContext is a reference to a cloned AppContext
	*context.AppContext
	// embedded struct of the options passed to NewGzipBundler
	*GzipBundlerOptions
	// backend is the selected backend, here, always BundlerBackendTarGzip, for aligning file endings
	backend BundlerBackend
	// tarWriter is a pointer to the writer context for tar files. It is configured to pipe its
	// output to the gzip writer for compression
	tarWriter *tar.Writer
	// gzipWriter is a pointer to the writer context for gzip encoding. Its output is set to
	// the file descriptor in GzipBundler.archive
	gzipWriter *gzip.Writer
	// archive is the file descriptor opened to write the gzipped tarball to. Calling Bundler.Close()
	// will close it, as well as the writers linked to it
	archive *os.File
	// sourceDirectory contains the directory from which to include additional files from.
	// It is the base for the Include patterns in Configuration.Spec.BundleOptions.
	sourceDirectory string
	// files contains paths to all auxilliary files to be added to an archive
	files []additionalFile
}

// Compile-time check for GzipBundler implmenting the Bundler interface
var _ Bundler = &GzipBundler{}

// NewGzipBundler makes a new gzip+tar bundler that encodes a tarball using gzip.
// It returns a GzipBundler that implements the Bundler interface
//
// If the archive file descriptor cannot be created, NewGzipBundler returns an error.
func NewGzipBundler(ctx *context.AppContext, files []additionalFile, options *GzipBundlerOptions) (*GzipBundler, error) {
	archive, err := os.Create(filepath.Join(options.ArtifactsDirectory, options.ArchiveName))
	if err != nil {
		return nil, err
	}

	gzipWriter := gzip.NewWriter(archive)
	tarWriter := tar.NewWriter(gzipWriter)

	bundler := &GzipBundler{
		backend:            BundlerBackendTar,
		AppContext:         ctx,
		GzipBundlerOptions: options,
		tarWriter:          tarWriter,
		gzipWriter:         gzipWriter,
		files:              files,
		sourceDirectory:    options.SourceDirectory,
	}

	return bundler, nil

}

// Close closes the remaining file descriptors for the tar archive
func (b *GzipBundler) Close() error {
	defer b.archive.Close()
	defer b.tarWriter.Close()
	return b.gzipWriter.Close()
}

// Type returns the static type of the bundler
func (b *GzipBundler) Type() BundlerBackend {
	return BundlerBackendTarGzip
}

// AddAssignmentToArchive implements writing the assignment's PDF to the tar archive
func (b *GzipBundler) AddAssignment() error {
	if b.tarWriter == nil {
		return errors.New("writer not created yet")
	}
	assignmentPdfName := fmt.Sprintf("%s.pdf", filepath.Base(b.SourceDirectory))
	src, err := os.Open(filepath.Join(b.ArtifactsDirectory, assignmentPdfName))
	if err != nil {
		return err
	}

	info, err := src.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	err = b.tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(b.tarWriter, src)
	if err != nil {
		return err
	}
	return nil
}

func (b *GzipBundler) AddAuxilliaryFiles() error {
	for _, f := range b.files {
		if err := b.addAuxilliaryFile(f); err != nil {
			return err
		}
	}
	return nil
}

// addAuxilliaryFile opens a file descriptor for the file and
// writes it to the tar archive file
func (b *GzipBundler) addAuxilliaryFile(file additionalFile) error {
	if b.tarWriter == nil {
		return errors.New("writer not created yet")
	}
	fd, err := os.Open(filepath.Join(b.sourceDirectory, file.rootPath))
	if err != nil {
		return err
	}

	info, err := fd.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// remove root-level path prefix
	header.Name = file.archivePath

	err = b.tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	// write to the tarball, its output is automatically piped to the gzip writer
	_, err = io.Copy(b.tarWriter, fd)
	if err != nil {
		return err
	}

	return nil
}
