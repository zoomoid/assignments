# csassignments

This class file is designed for assignments of university courses, specifically
in the field of computer science, but can easily be adapted to other fields. It
provides macros for counting exercise points, adding a grading table, a helpful
document title block including several non-standard information, and page
headers with compressed author details.

It supports German and English language by default, by requirement of the
author's university.

Additional to the layout utilities, the class provides several commonly used
macros for computer science topics, namely several mathematical operators,
mathmode utilities and special environments for proofs and theorem. Those are a
condensed version of macros defined in <https://github.com/zoomoid/AlphabetClasses>.

## User Guide

Use the class like any default document class:

```latex
  \documentclass[<option>]{csassignments}

  \course{}
  \group{}
  \due{}

  \member{}{}{}
  % ...

  \begin{document}
  \maketitle
  \gradingtable

  \exercise[<Exercise Points>]{<Exercise Title>}
  \subexercise{}
  \end{document}
\end{verbatim}
```

Because it inherits from `article`, you can pass down any options that
`article` understands.
