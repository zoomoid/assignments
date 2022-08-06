package release

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lithammer/dedent"
	"github.com/zoomoid/assignments/v1/internal/bundle"
)

var (
	GitlabCIEnvFileTemplates = dedent.Dedent(`
		ASSIGNMENT="{{.Assignment}}"
		TAG="{{.Tag}}"
		ARTIFACTS_ID="{{.ArtifactsId}}"
		ARCHIVE_NAME="{{.ArchiveName}}"
		PDF_NAME="{{.PdfName}}"
		PDF_ASSETS='{{.PdfAssets}}'
		ARCHIVE_ASSETS='{{.ArchiveAssets}}'
	`)

	GitlabCITemplate = dedent.Dedent(`
	stages:
		- build-utility
		- build
		- pre-release
		- release
	build:
		stage: build
		image: ghcr.io/zoomoid/assignments/cli:latest
		script:
			- assignmentctl build --all
		artifacts:
			paths:
				- dist/
			expire_in: 4 months
			reports:
				dotenv: artifacts.env
	release:
		stage: release
		image: ghcr.io/zoomoid/assignments/ci/gitlab:latest
		rules:
			- if: $CI_COMMIT_TAG && $CI_COMMIT_TAG =~ /^assignment-[0-1][0-9]$/
		script:
			- assignmentctl ci release gitlab > .env
			- source .env 
			- release-cli create 
				--tag-name $CI_COMMIT_TAG 
				--name "Assignment $ASSIGNMENT"
				--description "Release assignment $ASSIGNMENT for $CI_PROJECT_NAME from CI"
				--assets-link $ARCHIVE_ASSETS
				--assets-link $PDF_ASSETS
	`)
)

type Asset struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	LinkType string `json:"link_type,omitempty"`
	Filepath string `json:"filepath,omitempty"`
}

func TemplateGitlabCIEnvFile(artifactsDirectory string, archiveNameTemplate string, ad map[string]interface{}) (*bytes.Buffer, error) {
	tag := os.Getenv("CI_COMMIT_TAG")
	artifactsId := os.Getenv("CI_JOB_ID")
	projectURL := os.Getenv("CI_PROJECT_URL")
	assignment := strings.Replace(tag, "assignment-", "", 1)

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

	archiveName := filepath.Base(matches[0])
	pdfName := fmt.Sprintf("assignment-%s.pdf", assignment)

	tmpl, err := template.New("envfile").Parse(GitlabCIEnvFileTemplates)
	if err != nil {
		return nil, err
	}

	archiveArtifact := Asset{
		Name: "Submittable archive",
		URL:  fmt.Sprintf("%s/-/%s/artifacts/file/dist/%s", projectURL, artifactsId, archiveName),
	}

	pdfArtifact := Asset{
		Name: "PDF",
		URL:  fmt.Sprintf("%s/-/%s/artifacts/file/dist/%s", projectURL, artifactsId, pdfName),
	}

	marshalledArchiveArtifact, err := json.Marshal(archiveArtifact)
	if err != nil {
		return nil, err
	}
	marshalledPDFArtifact, err := json.Marshal(pdfArtifact)
	if err != nil {
		return nil, err
	}

	d := struct {
		Assignment    string
		Tag           string
		ArtifactsId   string
		PdfName       string
		ArchiveName   string
		PdfAssets     string
		ArchiveAssets string
	}{
		Assignment:    assignment,
		Tag:           tag,
		ArtifactsId:   artifactsId,
		PdfName:       pdfName,
		ArchiveName:   archiveName,
		PdfAssets:     string(marshalledPDFArtifact),
		ArchiveAssets: string(marshalledArchiveArtifact),
	}
	output := &bytes.Buffer{}
	err = tmpl.Execute(output, d)
	if err != nil {
		return output, err
	}
	return output, nil
}
