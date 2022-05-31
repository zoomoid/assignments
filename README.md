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
assignments that was generally accepted, `assignments` was created.

## assigments.cls

`assignments` brings several macros for working on at least Computer Science
assignments, but you can always add more functionality. For more details,
consider the classes' documentation.

Use it either in the document's root directory, or add it to your personal texmf
directory:

```latex
\documentclass{assignments}
\title{}
\author{}
\date{}

\course{} % course will replace title in maketitle
\group{} % group works like a subtitle
\member{}{}{} % members are like authors, but with a separation into ID, name, and surname
\due{} % add a due date, automatically prefixed with "Handed in on"

\begin{document}
\maketitle % makes the documents header
\gradingtable{} % adds a grading table between title and exercises

\exercise[<points>]{<title>} % add an exercise

\subexercise[<points>]{<subtitle>} % add a subexercise
% (points are ignored, only top-level exercise points are added up to sum in grading table)

\end{document}
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