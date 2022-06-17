# zoomoid/assignments

## Preamble

Hand-in assignments are a common thing among computer science courses at
University. So much so, that there's enough pain with academia using LaTeX for
anything and everything and the creation of third-party tooling for making
working with LaTeX a little bit less painful is justified.

Over the course of several semesters, I've written a lot of LaTeX code and read
even more documentation and had my fair share of experience with building more
and more complex things, foremost a library if LaTeX document classes and
packages for all sorts of stuff: <https://github.com/zoomoid/alphabetclasses>.

Several of these macros and packages only ever experienced one iteration, and
where immediately superseeded by another class or package. Finally, after years
of iterating over which macros where essential and designing a layout for
assignments that was generally accepted, `csassignments` was created.

## csassigments.cls

`csassignments` brings several macros for working on at least Computer Science
assignments, but you can always add more functionality. For more details,
consider the classes' documentation, either from source ([TeXDoc](./tex/latex/csassignments/csassignments.dtx)) or
online on CTAN (<https://www.ctan.org/pkg/csassignments>).

Use it either in the document's root directory, or add it to your personal texmf
directory:

```latex
\documentclass{csassignments}

\course{} % course will replace title in maketitle
\group{} % group works like a subtitle
\member[]{} % members are like authors, but with a separation into ID, name, and surname
\due{} % add a due date, automatically prefixed with "Handed in on"

\begin{document}
\maketitle % makes the documents header
\gradingtable % adds a grading table between title and exercises

\exercise[<points>]{<title>} % add an exercise

\subexercise[<points>]{<subtitle>} % add a subexercise
% (points are ignored, only top-level exercise points are added up to sum in grading table)

\end{document}
```

## Prerequisites for LaTeX

To compile your assignments, you'll have to have access to your local LaTeX installation to install packages.
Depending on your TeX distribution and operating system, this varies, also whether you install packages locally
or system-wide.

To reduce comlexity for all those different deployment methods, we push installing the LaTeX class down to your
TeX distro:

1. **MikTeX**: You won't need to do anything, MikTeX will download the class from your closest CTAN mirror. You might have to confirm installing the class
2. **TeXLive**: If something different to `texlive-full` is installed on your system, run `tlmgr install csassignments` to install the class for your user's TeX installation
3. **Other distros**: Follow a guide on how to install packages.

If you don't require building the assignments (simply because you can flawlessly write TeX code) or chose to run LaTeX in a container (e.g. `miktex/miktex`), the installation is either the same as if you run from the host system, *or* you use
the container image provided from this repository:

> Depending on your container runtime, you might have to adjust some of the command's arguments

```bash
# Run a bash inside the container to build your assignments with LaTeX directly 
#
# On UNIX
$ docker run -ti -v $(pwd):/work ghcr.io/zoomoid/assignments bash
# In PowerShell for Windows users
$ docker run -ti -v ${PWD}:/work ghcr.io/zoomoid/assignments bash
```

The image is Ubuntu-based, and installs the most recent version of MikTeX available. MikTeX will then install the `csassignments` class on-demand.

If you'd like to cache the installed packages for later use again, provide an additional (anonymous) volume for the
installation path of MikTeX: `-v miktex:/miktex/.miktex`.

## Installation

Download the binary fit for your OS from the Release page or build the CLI from source.

> TODO: add build instructions

Afterwards, move the binary created somewhere into your $PATH such that you can use it directly
in your command line.

## Usage

### Initialization

In a fresh directory (or repository), run

```bash
$ assignments bootstrap
```

If no further flags are set, the CLI will prompt you for several properties, namely:

1. Course name
2. Group name
3. Group Members, passed in the format of `"<Name>;<ID>,<Name 2>;<ID 2>,..."`

You can also provide those properties as flags to the `bootstrap` command, see the CLI's manual.

Afterwards, the directory will be initialized with a configuration file that contains all the entered
information, as well as additional metadata required for templates and more:

```yaml
# .assignments.yaml
spec:
  course: "Linear Algebra I"
  group: "Group Alpha"
  members:
    - id: "123456"
      name: "Max Mustermann"
    - id: "AB123456"
      name: "Erika Mustermann"
    - id: "69420"
      name: "Kim Took"
  includes: []
status:
  assignment: 1
```

### Starting a new Assignment

Calling

```bash
$ assignments generate
```

will create fresh scaffolding TeX file for the current assignment. If no further arguments are provided, the
configuration file will be used for the current assignment's number. The CLI will prompt you for the due date
of the assignment. If you chose to not provide one, or you simply don't know it, just press `Enter` to skip the field.

You can also pass the assignment's number to the `generate` command, which will also update the status field in your
`.assignments.yaml`.

If you just want to generate a new assignment without the side-effects of incrementing the status counter in the config file,
pass the `--no-increment` flag to `generate`.

If you need to override an existing assignment, there's also a `--force` flag, but beware that any existing files are overridden,
as the directory for the assignment is created from scratch.

Then, proceed to edit your assignment. Documentation of the `csassignments` class can be found at <https://www.ctan.org/pkg/csassignments>.

### Building Assignments

Run

```bash
# Compile a specific assignment
$ assignments build $ASSIGNMENT_NO
# Compile *all* assignments
$ assignments build --all
# Compile an assignment and override any already existing artifacts
$ assignments build --force $ASSIGNMENT_NO
```

## Tooling

To make collaboration easier and enable templating and automation of as much as
possible, in particular considering a weekly/bi-weekly schedule of assignments
per-course, we built several tools for helping with that.

1. CI pipelines: Leveraging the advantages of CI in the context of building
   assignment PDFs from source seemed only logical: It introduces
   reproducability and also, not everyone needs to install a fully-blown TeX
   toolchain into their system. Rather, with things such as containers, we can
   delegate building entirely to disposable containers, that run the build and
   export the PDF as artifacts from the pipeline (One might even say that the
   _only_ purpose of all this effort was alleviating the need for a TeX
   toolchain on my personal system entirely).
2. Easily template repeating code: The scaffolding you see above; it's
   essentially required every week again. Instead of copying the source files
   and manually adjusting the parameters to fit the current assignment, having
   tools that generate the source files automatically reduces friction and
   mistakes that occur from copy-pasting over and over again. A single command
   is ideally all you need to start on a new assignment. Also it allows us to
   enforce a repeating structure for the directories of the assignments. Further
   downstream, we can then leverage the repeating structure to run further
   automation
3. Automatic bundling of assignment artifacts: If an assignment requires more
   than just a single PDF, e.g., code, using the fixed structure of the file system allows
   for bundling, again, with a single command. With this, we automatically create ZIP archives from
   the CLI and can immediately upload these to the submission platform.
4. Integration into a consise Git worlflow: We can leverage things such as tags
   to mark exercises as being done and trigger particular actions on the final
   output.

This resulted in building the first CLI for handling assignment templating, building, and bundling.
It was written in Python, as a singular script.