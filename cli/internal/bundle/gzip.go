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
	ArchiveName        string
	AssignmentBase     string
	ArtifactsDirectory string
}

type GzipBundler struct {
	*context.AppContext
	backend    BundlerBackend
	tarWriter  *tar.Writer
	gzipWriter *gzip.Writer
	*TarBundlerOptions
	archive         *os.File
	sourceDirectory string
}

var _ Bundler = &GzipBundler{}

func NewGzipBundler(ctx *context.AppContext, options *TarBundlerOptions) (*GzipBundler, error) {
	archive, err := os.Create(filepath.Join(options.ArtifactsDirectory, options.ArchiveName))
	if err != nil {
		return nil, err
	}

	gzipWriter := gzip.NewWriter(archive)
	tarWriter := tar.NewWriter(gzipWriter)

	bundler := &GzipBundler{
		backend:           BundlerBackendTar,
		AppContext:        ctx,
		TarBundlerOptions: options,
		tarWriter:         tarWriter,
		gzipWriter:        gzipWriter,
		sourceDirectory:   options.AssignmentBase,
	}

	return bundler, nil

}

// Close closes the remaining file descriptors for the tar archive
func (b *GzipBundler) Close() error {
	defer b.archive.Close()
	defer b.tarWriter.Close()
	return b.gzipWriter.Close()
}

// AddAssignmentToArchive implements writing the assignment's PDF to the tar archive
func (b *GzipBundler) AddAssignment() error {
	if b.tarWriter == nil {
		return errors.New("writer not created yet")
	}
	assignmentPdfName := fmt.Sprintf("%s.pdf", b.AssignmentBase)
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

// AddAuxilliaryFileToArchive opens a file descriptor for the file and
// writes it to the tar archive file
func (b *GzipBundler) AddAuxilliaryFile(filename string) error {
	if b.tarWriter == nil {
		return errors.New("writer not created yet")
	}
	file, err := os.Open(filepath.Join(b.sourceDirectory, filename))
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// remove root-level path prefix
	header.Name = filename

	err = b.tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(b.tarWriter, file)
	if err != nil {
		return err
	}

	return nil
}
