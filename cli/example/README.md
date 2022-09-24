# Configuration

Configuration of the local directory is done in YAML format. `assignmentctl
init` (or `bootstrap`) creates a minimal configuration file in the working
directory. If you want a configuration file that also contains all the default
values, run

```bash
assignmentctl init --full
```

An exemplary configuration file looks like this:

```yaml
# ./.assignments.yaml

# spec contains the main configuration fields left to the user
spec:
  # course name
  course: "Demo Course"
  # group name
  group: "Group A"
  # members:
  members:
    # each group member, with name and id
    - name: "Max Mustermann"
      id: "123456" # like student immatriculation ID
  # options for the 'generate' subcommand
  generate:
    # list all directories to create when using assignmentctl generate
    create:
      - feedback
      - code
  # options for the 'build' subcommand
  build:
    # recipe is a workflow description for building latex documents
    # There is a set of keywords that are expanded with Golang template syntax
    # For documentation of thsose, see below
    recipe:
      - command: latexmk
        args:
          - -interaction=nonstopmode
          - -pdf
          - -file-line-error
          - -shell-escape
          - -outdir="{{.OUTDIR}}"
          - "{{.DOC}}"
    # Configuration for cleanup
    cleanup:
      # cleanup by deleting all files that match the glob pattern
      glob:
        - "*.log"
        - "*.aux"
  bundle:
    # Name template for the bundles created.
    # _id and _format are derived automatically and should thus be treated as "internal"
    template: assignment-{{._id}}.{{._format}}
    # data is a map that contains any additional data used in the template
    data: {}
    # include contains glob patterns of all additional files to include in the bundle
    include:
      - code/**
      - figures/**
# status keeps track of the repository's state
status:
  # The current assignment number that we are currently at. The
  # model assumes monotonicity of assignment numbers
  assignment: 1
```

## Expansion of Dynamic Fields

To make recipes more flexible, we need some variables that are expanded at
runtime, e.g. filenames for LaTeX commands to run. Because Golang templates are
convenient tools for templating, we construct a struct that looks like this:

```go
type substitutionContext struct {
  DOC              string
  DOCEXT           string
  DIR              string
  TMPDIR           string
  WORKSPACE_FOLDER string
  RELATIVE_DIR     string
  RELATIVE_DOC     string
  OUTDIR           string
}
```

This defines all the names that are available for expansion in the recipe's
arguments and their semantics are as followed:

1. `DOC` is the *absolute* path of the document to build **without its file
   extension**. This is required for most latex/pdflatex commands
2. `DOCEXT` is the same as `DOC` but with the file extension (i.e. ".tex")
   included.
3. `DIR` is the *absolute* path to the directory of the file to build.
4. `TMPDIR` is the OS's temporary directory, suitable for storing any
   intermediate files that may be disposed of
5. `WORKSPACE_FOLDER` is the folder that the current command runs in
6. `RELATIVE_DIR` is the *relative* path to the directory of the file to build
7. `RELATIVE_DOC` is the *relative* path to the document to build
8. `OUTDIR` is the directory to which to build, and in our scenario always the
   same as the directory of the source file

You can use those in your arguments like usual Golang templates and they will be
expanded if found in any of the arguments.
