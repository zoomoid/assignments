package ci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/lithammer/dedent"
)

var (
	GitlabCIEnvFileTemplate = dedent.Dedent(`
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
    - build
    - release
  build:
    stage: build
    image: ghcr.io/zoomoid/assignments/runner:latest
    script:
      - assignmentctl build --all
    artifacts:
      paths:
        - dist/
      expire_in: 4 months
  release:
    stage: release
    image: ghcr.io/zoomoid/assignments/ci/gitlab:latest
    rules:
      - if: $CI_COMMIT_TAG && $CI_COMMIT_TAG =~ /^assignment-[0-9][0-9]+$/
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

	artifacts, err := archiveAndPdfName(assignment, artifactsDirectory, archiveNameTemplate, ad)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("envfile").Parse(strings.TrimSpace(GitlabCIEnvFileTemplate))
	if err != nil {
		return nil, err
	}

	archiveArtifact := Asset{
		Name: "Submittable archive",
		URL:  fmt.Sprintf("%s/-/%s/artifacts/file/dist/%s", projectURL, artifactsId, artifacts.Archive),
	}

	pdfArtifact := Asset{
		Name: "PDF",
		URL:  fmt.Sprintf("%s/-/%s/artifacts/file/dist/%s", projectURL, artifactsId, artifacts.PDF),
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
		PdfName:       artifacts.PDF,
		ArchiveName:   artifacts.Archive,
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
