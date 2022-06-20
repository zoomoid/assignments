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
	ArchiveName        string
	AssignmentBase     string
	ArtifactsDirectory string
}

type ZipBundler struct {
	*context.AppContext
	backend BundlerBackend
	writer  *zip.Writer
	*ZipBundlerOptions
	archive         *os.File
	sourceDirectory string
}

var _ Bundler = &ZipBundler{}

func NewZipBundler(ctx *context.AppContext, options *ZipBundlerOptions) (*ZipBundler, error) {

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
		sourceDirectory:   options.AssignmentBase,
	}

	return bundler, nil
}

// Close closes the remaining file descriptors for the zip archive
func (b *ZipBundler) Close() error {
	defer b.archive.Close()
	return b.writer.Close()
}

// AddAssignmentToArchive implements writing the assignment's PDF to the zip archive
func (b *ZipBundler) AddAssignment() error {
	if b.writer == nil {
		return errors.New("writer not created yet")
	}
	assignmentPdfName := fmt.Sprintf("%s.pdf", b.AssignmentBase)
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

// AddAuxilliaryFileToArchive opens a file descriptor for the file and
// writes it to the zip archive file
func (b *ZipBundler) AddAuxilliaryFile(filename string) error {
	if b.writer == nil {
		return errors.New("writer not created yet")
	}
	file, err := os.Open(filepath.Join(b.sourceDirectory, filename))
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
