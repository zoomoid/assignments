/*
Copyright 2022 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	// ArchiveName created in Bundler.New()
	ArchiveName string
	// SourceDirectory is the directory where the additional files and
	// the assignment's source originate
	SourceDirectory string
	// Artifacts directory is the directory where the PDF originates from
	ArtifactsDirectory string
}

type TarBundler struct {
	// context.AppContext is a reference to a cloned AppContext
	*context.AppContext
	// embedded struct of the options passed to NewGzipBundler
	*TarBundlerOptions
	// backend is the selected backend, here, always BundlerBackendTarGzip, for aligning file endings
	backend BundlerBackend
	// writer is a pointer to the writer context for tar files. It is configured to pipe its
	// output to the gzip writer for compression
	writer *tar.Writer
	// archive is the file descriptor opened to write the gzipped tarball to. Calling Bundler.Close()
	// will close it, as well as the writers linked to it
	archive *os.File
	// sourceDirectory contains the directory from which to include additional files from.
	// It is the base for the Include patterns in Configuration.Spec.BundleOptions.
	sourceDirectory string
	// files contains paths to all auxilliary files to be added to an archive
	files []additionalFile
}

// Compile-time check for TarBundler implementing the Bundler interface
var _ Bundler = &TarBundler{}

// NewTarBundler makes a new tar bundler that uses the archive/tar module.
// It returns a TarBundler that implements the Bundler interface
//
// If the archive file descriptor cannot be created, NewTarBundler returns an error.
func NewTarBundler(ctx *context.AppContext, files []additionalFile, options *TarBundlerOptions) (*TarBundler, error) {
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
		sourceDirectory:   options.SourceDirectory,
	}

	return bundler, nil

}

// Close closes the remaining file descriptors for the tar archive
func (b *TarBundler) Close() error {
	defer b.archive.Close()
	return b.writer.Close()
}

// Type returns the static type of the bundler
func (b *TarBundler) Type() BundlerBackend {
	return BundlerBackendTar
}

// AddAssignmentToArchive implements writing the assignment's PDF to the tar archive
func (b *TarBundler) AddAssignment() error {
	if b.writer == nil {
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
func (b *TarBundler) addAuxilliaryFile(file additionalFile) error {
	if b.writer == nil {
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

	err = b.writer.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(b.writer, fd)
	if err != nil {
		return err
	}

	return nil
}
