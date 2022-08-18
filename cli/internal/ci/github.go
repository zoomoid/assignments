package ci

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lithammer/dedent"
)

var (
	GithubActionEnvFileTemplate = dedent.Dedent(`
    ASSIGNMENT="{{.Assignment}}"
    PDF_ASSETS='{{.PdfAssets}}'
    ARCHIVE_ASSETS='{{.ArchiveAssets}}'
  `)

	GithubActionTemplate = dedent.Dedent(`
    name: Assignmentctl workflow

    on:
      push:
        branches:
          - "*"
        tags:
          - assignment-[0-9][0-9]+
    jobs:
      build:
        name: Build assignments
        container: ghcr.io/zoomoid/assignments/runner:latest
        steps:
          - uses: actions/checkout@v2

          - name: Build assignments to ./dist/
            run: assignmentctl build --all

          - name: Bundle assignments in ./dist/
            run: assignmentctl bundle --all

          - uses: actions/upload-artifacts@v3
            with:
              name: assignments
              path: ${{ github.workspace }}/dist/

      release:
        name: Release assignment
        container: ghrc.io/zoomoid/assignments/ci/github:latest
        needs: build
        if: ${{ github.event_name == 'create' && github.ref_type == 'tag' }}
        steps:
          - uses: actions/download-artifact@v3
            with:
              name: assignments

          - name: Create pre-release file with release data
            run: assignmentctl ci release github > .env

          - run: source .env

          - name: Create release with github-cli
            env:
              GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
            run: |
              gh release create ${{ github.ref_name }}
                --title "Assignment $ASSIGNMENT"
                --notes "Release assignment $ASSIGNMENT for ${{ github.ref_name }} from CI"
                $ARCHIVE_ASSETS
                $PDF_ASSETS
  `)
)

type Artifact struct {
	Name string
	Path string
}

func (a *Artifact) ToString() string {
	return fmt.Sprintf("%s#%s", a.Path, a.Name)
}

func TemplateGithubActionsEnvFile(artifactsDirectory string, archiveNameTemplate string, ad map[string]interface{}) (*bytes.Buffer, error) {
	tag := os.Getenv("GITHUB_REF_NAME")
	assignment := strings.Replace(tag, "assignment-", "", 1)

	artifacts, err := archiveAndPdfName(assignment, artifactsDirectory, archiveNameTemplate, ad)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("envfile").Parse(strings.TrimSpace(GithubActionEnvFileTemplate))
	if err != nil {
		return nil, err
	}

	pdfArtifacts := Artifact{
		Path: filepath.Join(artifactsDirectory, artifacts.PDF),
		Name: "PDF",
	}
	archiveArtifact := Artifact{
		Path: filepath.Join(artifactsDirectory, artifacts.Archive),
		Name: "Submittable archive",
	}

	d := struct {
		Assignment    string
		PdfAssets     string
		ArchiveAssets string
	}{
		Assignment:    assignment,
		PdfAssets:     pdfArtifacts.ToString(),
		ArchiveAssets: archiveArtifact.ToString(),
	}
	output := &bytes.Buffer{}
	err = tmpl.Execute(output, d)
	if err != nil {
		return output, err
	}
	return output, nil
}
