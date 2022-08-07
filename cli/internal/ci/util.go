package ci

import (
	"fmt"
	"path/filepath"

	"github.com/zoomoid/assignments/v1/internal/bundle"
)

type Artifacts struct {
	PDF     string
	Archive string
}

func archiveAndPdfName(assignment string, artifactsDirectory string, archiveNameTemplate string, ad map[string]interface{}) (*Artifacts, error) {
	ad["_id"] = assignment
	ad["_format"] = "*" // glob the archive name later so that the actual bundle's format is irrelevant
	archiveGlobName, err := bundle.MakeArchiveName(archiveNameTemplate, ad)
	if err != nil {
		return nil, err
	}

	absArchiveGlobName := filepath.Join(artifactsDirectory, archiveGlobName)

	matches, err := filepath.Glob(absArchiveGlobName)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("cannot find archive in artifacts directory")
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("archive name is ambiguous, can only export a single archive per tag")
	}

	archive := filepath.Base(matches[0])
	pdf := fmt.Sprintf("assignment-%s.pdf", assignment)

	return &Artifacts{
		PDF:     pdf,
		Archive: archive,
	}, nil
}
