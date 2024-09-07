# zoomoid/assignments

## Preamble

Hand-in assignments are a common thing among computer science courses at
University. So much so, that there's enough pain with academia using LaTeX for
anything and everything and the creation of third-party tooling for making
working with LaTeX a little bit less painful is justified.

Over the course of several semesters, I've written a lot of LaTeX code and read
even more documentation and had my fair share of experience with building more
and more complex things, foremost a library of LaTeX document classes and
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
\sheet{} % adds a sheet identifier, automatically prefixed with "Exercise Sheet"
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

If you don't require building the assignments (simply because you can flawlessly write TeX code) or chose to run LaTeX in a container (e.g. `miktex/miktex`), the installation is either the same as if you run from the host system, _or_ you use
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

## Installation and Usage of `assignmentctl`

To learn more about the CLI that accompanies this project's LaTeX class, see [./cli](./cli/README.md).

For information about the older version of this tool, see [./pyassignmentctl](./pyassignmentctl/README.md).

## Usage with CI Pipelines

For details on how to use this project in Continuous Integration pipelines, see [./ci](./ci/README.md).

## Epilogue: The "Why" of assignmentctl

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
2. Easily template repeating code: The scaffolding essentially is required every week
   again. Instead of copying the source files
   and manually adjusting the parameters to fit the current assignment, having
   tools that generate the source files automatically reduces friction and
   mistakes that occur from copy-pasting over and over again. A single command
   is ideally all you need to start on a new assignment. Also it allows us to
   enforce a repeating structure for the directories of the assignments. Further
   downstream, we can then leverage the repeating structure to run further
   automation
3. Automatic bundling of assignment and additional files: If an assignment requires more
   than just a single PDF, e.g., code, using the fixed structure of the file system allows
   for bundling, again, with a single command. With `assignmentctl`, we automatically create
   archives from the CLI and can immediately upload these to the submission platform.
4. Integration into Git/VCS workflow: We can leverage versioning systems such as tags
   to mark exercises as being done and trigger particular actions on the final iteration output.

This resulted in building the first CLI for handling assignment templating, building, and bundling.
It was written in Python, as a singular script.
