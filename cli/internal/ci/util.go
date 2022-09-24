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
