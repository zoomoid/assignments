package bundle

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zoomoid/assignments/v1/internal/context"
)

type TarBundlerOptions struct {
	ArchiveName        string
	AssignmentBase     string
	ArtifactsDirectory string
}

type TarBundler struct {
	*context.AppContext
	backend BundlerBackend
	writer  *tar.Writer
	*TarBundlerOptions
	archive         *os.File
	files           []string
	sourceDirectory string
}

var _ Bundler = &TarBundler{}

func NewTarBundler(ctx *context.AppContext, files []string, options *TarBundlerOptions) (*TarBundler, error) {
	archive, err := os.Create(filepath.Join(options.ArtifactsDirectory, options.ArchiveName))
	if err != nil {
		return nil, err
	}

	writer := tar.NewWriter(archive)

	bundler := &TarBundler{
		backend:           BundlerBackendTar,
		AppContext:        ctx,
		TarBundlerOptions: options,
		writer:            writer,
		files:             files,
		sourceDirectory:   options.AssignmentBase,
	}

	return bundler, nil

}

// Close closes the remaining file descriptors for the tar archive
func (b *TarBundler) Close() error {
	defer b.archive.Close()
	return b.writer.Close()
}

// AddAssignmentToArchive implements writing the assignment's PDF to the tar archive
func (b *TarBundler) AddAssignment() error {
	if b.writer == nil {
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

	err = b.writer.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(b.writer, src)
	if err != nil {
		return err
	}
	return nil
}

func (b *TarBundler) AddAuxilliaryFiles() error {
	for _, f := range b.files {
		if err := b.addAuxilliaryFile(f); err != nil {
			return err
		}
	}
	return nil
}

// addAuxilliaryFile opens a file descriptor for the file and
// writes it to the tar archive file
func (b *TarBundler) addAuxilliaryFile(filename string) error {
	if b.writer == nil {
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

	err = b.writer.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(b.writer, file)
	if err != nil {
		return err
	}

	return nil
}
