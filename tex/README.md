# assignments/tex

This directory structure mimics a TDS structure such that we can copy the structure to a container image or the user's home directory for local installation.

To build the documentation PDF, run `pdflatex -interaction=nonstopmode csassignments.dtx`. Clean up using e.g. `latexmk -c csassignments.dtx`.

To publish new versions of the package, bundle a zip file of the `csassignments` directory in the following structure:

```text
csassignments-vX.Y.Z.zip
  csassignments/
    README.md
    csassignments.dtx
    csassignments.ins
    csassignments.pdf
```

When there's nothing else in the directory, run `zip -r csassignments-vX.Y.Z.zip csassignments` from the `tex/latex` directory.
