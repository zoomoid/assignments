# assignments/ci

To support releasing assignments from CI runners in combination with the
`assignmentctl ci release` command, create either a Gitlab CI file or a Github
action in your repository:

## Gitlab

Your pipeline file `.gitlab-ci.yml` should look something like this:

```yaml
# .gitlab-ci.yml

stages:
  - build
  - release
build:
  stage: build
  image: ghcr.io/zoomoid/assignments/runner:latest
  script:
    - assignmentctl build --all
    - assignmentctl bundle --all
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
```

Note the two stages, `build` and `release`. Building with `assignmentctl` will
create PDF files and bundles in the exported `./dist` directory. The release
stage depends on those files to be present in order to add them to the release;
Gitlab links to Job artifacts from CI pipelines and this coupling is leveraged
in `assignmentctl ci release gitlab` to create the necessary URLs in
`$ARCHIVE_ASSETS` and `$PDF_ASSETS` for the release.

If you do not want artifacts to be added to the release, simply remove the lines
`--assets-link $ARCHIVE_ASSETS` and `--assets-link $PDF_ASSETS` from the job's
script.

The release job is only ever run if you push a new tag of the form
`assignment-[0-9][0-9]+`, e.g. `assignment-03`. The build job is run on every
push to the repository.

If you do not want to build all PDFs (by far the most time-consuming task due to
`latexmk`) in every pipeline, you can scope the job differently, e.g. also only
run it on pushes of a tag (as done in the release job), or only build a specific
assignment (either statically or dynamically), you could add
`NO=$(echo $CI_COMMIT_TAG | sed -rn "s/assignment-([0-9][0-9]+)/\1/p")` to capture
the specific assignment denoted by the tag name and run `assignmentctl build $NO`
instead of `assignmentctl build --all` (and same for bundling, respectively, but
this is far less time-consuming than running `latexmk` at least 2 times on each
assignment). If you'd choose to do assignments branch-based with branch names of
the form `assignment-XY`, you could do the same as above but with
`$CI_COMMIT_BRANCH` instead of `$CI_COMMIT_TAG`.

## Github

Github's action model is similar to Gitlab's pipelines, but slightly different
in specifics. Here's a template action to be stored at
`./.github/workflows/build.yml`:

```yaml
# .github/workflows/build.yml
name: Assignmentctl workflow

on:
  push:
    branches:
      - "*"
  push:
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
    if: ${{ github.event_name == 'create' && github.ref_type == 'tag'  }}
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
```

You can customize the actions the same way you would for Gitlab CI pipelines if
you require more strict filters for events such as only building on pushes of
tags or only building off of specific branches and re-using the ref names in the
steps.

Here, Github Actions are a lot more flexible in terms of inclusion of
environment than Gitlab and in theory should be more easily customizable.

## Container Images

To support the two workflow types, we have built specific container images that
include the specific CLI needed to interact with the API of the respective SCM
provider. Both images are based on the default `assignmentctl` CLI Alpine Linux
image and just add the respective provider's CLI.
