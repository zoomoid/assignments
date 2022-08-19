# pyassignmentctl

> This is project is archived and should no longer be used. It is without maintenance and also without
> proper testing of edge cases and can potentially cause harm to your files. If your keen of using it
> make sure to run it inside a clean working tree of your repository to be able to revert changes.

pyassignmentctl is the reference implementation for the assignmentctl CLI, originally written in Python3.

It is a singular script built to "do it all" with argparse, from templating TeX files to running latexmk,
bundling to zip and running in CI to export artifacts.

It is however, not really flexible, as it was created with some strict assumptions, originally only tailored
for a singular specific course, later adapted to another (which caused a lot of hassle actually).

It does most of its work by running shell commands, i.e., running latexmk, running `zip` for bundling etc.

There also is very little configurability: you cannot alter the default TeX file template without altering
the code, you cannot alter the filename schema for bundling the files, and generally the configuration in
.ini format is really inflexible.

These are all points that were adressed in the design of the newer, better CLI that is `assignmentctl`.
