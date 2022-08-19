# assignments/cli

`assignmentsctl` CLI to use for managing an assignment directory from start to
hand-in.

## Installation

Download the latest binary for your system from the Releases page and unzip the
archive to somewhere in your `$PATH`. Afterwards you can use the CLI in any
directory.

If you'd like to not install anything into your host system, you can always run
the CLI inside a container: run `docker run -ti ghcr.io/zoomoid/assignments/cli ...` and add a mount argument to mount your assignments directory from your
host.

## Usage

It all starts with the initialization of a new course directory:

```bash
# start a new course directory
$ mkdir linear-algebra-2022 && cd linear-algebra-2022
# run the initialization
$ assignmentctl init # or assignmentctl bootstrap
# You can also immediately create a git repository in the directory
$ assignmentctl init --git
```

After filling in the prompts with your specific values, this creates a
configuration file that marks the root of the assignments directory. You can use
the assignmentctl CLI in any subdirectory of this root and you will always use
the closest configuration file, i.e., `.assignments.yaml`.

For details on how to customize the configuration file, see
[./hack/README.md](./hack/README.md).

The following commands assume state, as we store the current assignment's number
inside the config file. When called without arguments, `assignmentctl` will use
the status field in the config file for reference.

To start a new assignment, use `assignmentctl generate`. The `generate` command
templates a TeX file using the `csassignments` class and fills out some preamble
fields automatically from config.

To build the newly created and filled-out assignment, run `assignmentctl build`.
This will run latexmk by default (or any other recipe specified in the
configuration file), and spit out a PDF in the `dist/` directory. For more
information about the `build` command, see `assignmentctl build --help`.

To create a bundle of all files required for the assignments, e.g., additional
code files, figures not embedded in the PDF and more, run `assignmentctl bundle`. It will bundle all files specified in the config file for a specific
assignment in a zip or a tarball. For more information about the `bundle`
command, see `assignmentctl bundle --help`.

The same workflow is used for any subsequent assignments:

1. Run `assignmentctl generate` for a fresh TeX file.
2. Work on the assignment.
3. Build the assignment PDF with latexmk or any other tool.
4. Optionally, bundle the assignment into a zip archive or a tarball.

## Commands

Commands are listed in order of their common usage pattern.

<details>
  <summary><code>assignmentctl bootstrap</code></summary>
  <code>

    Running bootstrap will create a local configuration file for the current
    working directory containing minimal information about the course,
    passed in either as flags or during interactive prompts.

    The configuration file can be customized further afterwards, e.g., by
    adding different build recipes and bundling options. For this, see
    documentation.

    Usage:
      assignmentctl bootstrap [flags]

    Aliases:
      bootstrap, init, initialize

    Flags:
          --course string      Course name
          --full               Include all defaults in configuration file
          --git                Create a git repository in the current directory and commit the configuration file immediately
          --group string       Group name
      -h, --help               help for bootstrap
          --includes strings   Custom TeX includes for the template. Paths are relative to the REPOSITORY root, not the actual assignment source file
          --members strings    Group members, as comma-separated <Name>;<ID> tuples

    Global Flags:
      -v, --verbose   Sets logging verbosity level to high

  </code>
</details>

<details>
  <summary><code>assignmentctl generate</code></summary>
  <code>
    The command generates a new assignment, either given by number
    as an argument to the command, or otherwise from the local
    configuration file, which keeps track of the upstream assignment.

    Generating (or templating) a new assignment requires a due date.
    As this is usually given, you can either use the --due flag, or
    wait for the CLI to prompt you. If however the due date is *not*
    provided by the assignment, just pressing ENTER during the prompt
    will leave it empty and thus not printed in the assignment's
    header.

    You can make the command skip incrementing the status counter in
    the local configuration file by passing the --no-increment flag.

    If there already exists an assignment in the target directory,
    the command will exit with an error. If however you pass the
    --force flag, any files in the target directory will be overriden.
    Be careful!

    The default template for new assignments looks like this:

    \documentclass{csassignments}
    {{- range $_, $input := .Includes -}}
    \input{ {{- $input -}} }
    {{ end }}
    \course{ {{- .Course -}} }
    \group{ {{- .Group | default "" -}} }
    \sheet{ {{- .Sheet | default "" -}} }
    \due{ {{- .Due | default "" -}} }
    {{- range $_, $member := .Members }}
    {{- $firstname := ($member.Name | splitList " " | initial | join " ") | default "" -}}
    {{- $lastname := ($member.Name | splitList " " | last) | default "" -}}
    \member{ {{- $firstname -}} }{ {{- $lastname -}} }{ {{- $member.ID -}} }
    {{ end }}
    \begin{document}
    \maketitle
    \gradingtable

    % Start the assignment here

    \end{document}

    You can provide your own template from the configuration file, by
    setting .spec.template to a Golang template. You can use any Sprig
    template function in your custom template.

    The command creates a new directory from the current assignment number,
    as well as all directories defined in the .spec.generate.create list.

    Usage:
      assignmentctl generate [flags]

    Flags:
          --due string     Due date of the assignment to generate. If not provided, you'll be prompted for a due date
      -f, --force          Overrides any existing assignment source files
      -h, --help           help for generate
          --no-increment   Skip incrementing assignment number in configuration

    Global Flags:
      -v, --verbose   Sets logging verbosity level to high

  </code>
</details>
<details>
  <summary><code>assignmentctl build</code></summary>
  <code>

    The command builds a selected assignment, either from arguments,
    or from the state of the local configuration file using latexmk
    with the underlying LaTeX distro. After successful build, the
    artifact files are copied to a common output directory, commonly
    ./dist/.

    To build *all* assignments found in the working directory, add the
    --all (or -a) flag.

    After compilation, the command also cleans up any intermediate files
    created by the LaTeX compiler. By default, the cleanup will be done
    directly in the file system, by using Glob patterns for a large
    set of intermediate TeX files. You can override this behaviour in
    two ways:

    1) you can specify a different set of Glob patterns in the config at
    .spec.build.cleanup.glob.patterns. If you'd like to run the cleanup
    recursively, also set .spec.build.cleanup.glob.recursive to true.
    Note that the Glob patterns are not merged with the default ones:
    If you provide your own, these are the complete ones to cleanup

    2) you can change the execution from using Globs to running commands,
    e.g. latexmk -C: For this, set .spec.build.cleanup.command.recipe
    accordingly

    Note that .spec.build.cleanup.command and .spec.build.cleanup.glob
    are mutually exclusive. Presence of both will cause the CLI to throw
    an error.

    If you use the build command in a setup different to one-off runs,
    for which you might want to keep the files for later runs again to save
    times, you can use --keep to preserve those intermediate files.

    You can suppress the output of spawned shell commands by passing
    --quiet, or -q.

    To adjust the build recipe for compilation, add a recipe to your
    configuration file at .spec.build.recipe. Recipes are order-preservent
    lists of commands with arguments in YAML format. A recipe consists
    of Tools, which must at least contain a .command string, and may
    include arbitrary .args as a YAML list.

    Usage:
      assignmentctl build [flags]

    Flags:
      -a, --all           Build all assignments in assignment-*/
      -f, --file string   Specify a file to build, will override any derived behaviour from the repository's configmap
          --force         Override any existing assignments with the same name
      -h, --help          help for build
          --keep          Skip latexmk -C cleaning up all files in the source directory
          --quiet         Suppress output from latexmk subprocesses

    Global Flags:
      -v, --verbose   Sets logging verbosity level to high

  </code>
</details>

<details>
  <summary><code>assignmentctl bundle</code></summary>
  <code>

    The command builds a selected assignment, either from arguments,
    or from the state of the local configuration file using latexmk
    with the underlying LaTeX distro. After successful build, the
    artifact files are copied to a common output directory, commonly
    ./dist/.

    To build *all* assignments found in the working directory, add the
    --all (or -a) flag.

    After compilation, the command also cleans up any intermediate files
    created by the LaTeX compiler. By default, the cleanup will be done
    directly in the file system, by using Glob patterns for a large
    set of intermediate TeX files. You can override this behaviour in
    two ways:

    1) you can specify a different set of Glob patterns in the config at
    .spec.build.cleanup.glob.patterns. If you'd like to run the cleanup
    recursively, also set .spec.build.cleanup.glob.recursive to true.
    Note that the Glob patterns are not merged with the default ones:
    If you provide your own, these are the complete ones to cleanup

    2) you can change the execution from using Globs to running commands,
    e.g. latexmk -C: For this, set .spec.build.cleanup.command.recipe
    accordingly

    Note that .spec.build.cleanup.command and .spec.build.cleanup.glob
    are mutually exclusive. Presence of both will cause the CLI to throw
    an error.

    If you use the build command in a setup different to one-off runs,
    for which you might want to keep the files for later runs again to save
    times, you can use --keep to preserve those intermediate files.

    You can suppress the output of spawned shell commands by passing
    --quiet, or -q.

    To adjust the build recipe for compilation, add a recipe to your
    configuration file at .spec.build.recipe. Recipes are order-preservent
    lists of commands with arguments in YAML format. A recipe consists
    of Tools, which must at least contain a .command string, and may
    include arbitrary .args as a YAML list.

    Usage:
      assignmentctl build [flags]

    Flags:
      -a, --all           Build all assignments in assignment-*/
      -f, --file string   Specify a file to build, will override any derived behaviour from the repository's configmap
          --force         Override any existing assignments with the same name
      -h, --help          help for build
          --keep          Skip latexmk -C cleaning up all files in the source directory
          --quiet         Suppress output from latexmk subprocesses

    Global Flags:
      -v, --verbose   Sets logging verbosity level to high
    λ /git/containers/assignments/cli/hack/ main* assignmentctl bundle --help

    Bundling compiles all files relevant for an assignment into an archive
    format. The backend defaults to zip, but can be set to tarball by
    passing the --tar flag. If you want to use tar and gzip, use --gzip.

    By default, every bundle includes at least the assignment's PDF from
    the ./dist/ directory. If you want to add further files or directories
    see the list .spec.bundle.include in your configuration file. It lets
    you specify files explicitly, or a glob pattern for multiple files,
    e.g. "code/*" or "figures/*.pdf". It is meant to complement the list
    of directories to create when using the generate command. The bundle will
    preserve the structure of the files included, and will have the PDF
    located at the archive's root.

    You can customize how the filename for the archive is generated. For this,
    you can set .spec.bundle.template to be an arbitrary Golang text template
    (including the use of sprig text functions). Just note that this is
    limited by what file paths are supported by your operating system, so
    don't get too crazy. The map in .spec.bundle.data is passed down to
    the template's execution for data binding.

    The default archive template is "assignment-{{._id}}.{{._format}}". Note
    the _id field: this is internally augmented from the command's arguments
    or the configuration's status field (or, in case of usage of --all, all
    available assignments in the repository). "format" is derived from the
    selected backend's common file extension, but respects overrides from
    the map at .spec.bundle.data, so you can also pick your own file extension
    without overriding the entire template.

    Usage:
      assignmentctl bundle [flags]

    Flags:
      -a, --all     Bundle all assignments
      -f, --force   Override any existing archives with the same name
          --gzip    Use tar and gzip as backend for archive bundling
      -h, --help    help for bundle
          --tar     Use tar as a backend for archive bundling

    Global Flags:
      -v, --verbose   Sets logging verbosity level to high

  </code>
</details>

<details>
  <summary><code>assignmentctl ci bootstrap {gitlab,github}</code></summary>
  <code>

    Run this command to quickly template CI files for the supported SCM providers,
    namely Gitlab and Github. Afterwards, you can customize them to your liking.

    To learn more about the CI integration, see the documentation at
    https://github.com/zoomoid/assignments/blob/main/ci/README.md

    Usage:
      assignmentctl ci bootstrap [flags]

    Flags:
      -f, --file string   Write the template directly to a file instead of Stdout
      -h, --help          help for bootstrap

    Global Flags:
      -v, --verbose   Sets logging verbosity level to high

  </code>
</details>

<details>
  <summary><code>assignmentctl ci release {gitlab,github}</code></summary>
  <code>

    The command is meant for usage inside CI pipelines to create release objects 
    for Gitlab and exports several environment variables in a file that are 
    required for the job running Gitlab's release-cli or Github's CLI, respectively.

    You can run this command outside of CI pipelines, note however that it is highly
    dependant on the ENV variables being available in the runner context. You will
    have to provide either $CI_COMMIT_TAG with $CI_JOB_ID and $CI_PROJECT_URL or
    $GITHUB_REF_NAME for the command to output correct .env files.

    Usage:
      assignmentctl ci release [flags]

    Flags:
      -f, --file string   Write the template directly to a file instead of Stdout
      -h, --help          help for release

    Global Flags:
      -v, --verbose   Sets logging verbosity level to high

  </code>
</details>

## Building from Source

You can build the CLI from source if you have go and make installed:

```bash
# Clone the repository
$ git clone https://github.com/zoomoid/assignments
# store the latest version from git tag in a variable
$ VERSION=$(git describe --abbrev=0 --tags)
# Install the go modules
$ go mod download
# Build the executable using make
$ make build
```

Afterwards, put the executable in `./bin/assignmentctl` to somewhere in your
`$PATH` to use it in the terminal.

## Contributing and Collaborating

The project is open for any sensible contribution. If you'd like to participate
in the development of the assignmentctl CLI, fork the repository and open a pull
request with your changes.

## Development Timeline

As of now, August 2022, the main functionality of the `assignmentctl` CLI is more or less done.
It handles creation, building, and bundling, and can assist you with running jobs in
CI pipelines on Gitlab and Github. The choice of those platforms is by availability.
I'm open to extension to other CI providers.

In the future, this section could contain an outlook on further development, if anything big comes up.
