\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{assignments}[2018/10/13 Assignment Class]
\DeclareOption*{\PassOptionsToClass{\CurrentOption}{article}}
\ProcessOptions \relax
\LoadClass[a4paper,11pt]{article}
\makeatletter
\usepackage[left=2.5cm,right=2.5cm,top=2cm,bottom=2cm,includeheadfoot]{geometry}
\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{titlesec,enumitem,graphicx,tocloft,float,ifthen,translations,csquotes,fancyhdr}

% Mathesymbole und ähnliches
\usepackage{microtype}
\usepackage{stmaryrd}

% Grafiken und PDFs einfügen
\usepackage{graphicx}
\usepackage{pdfpages}

% PDF-Tools
\usepackage[hidelinks, unicode]{hyperref}

% Abbildungen
\usepackage{tikz}
\usetikzlibrary{arrows,calc}
\usepackage{microtype}

\select@language{english}

% TRANSLATIONS
\DeclareLanguage{english}
\DeclareLanguage{german}

\ProvideTranslation{german}{prefixAuthor}{Eingereicht von}
\ProvideTranslation{english}{prefixAuthor}{Handed in by}
\ProvideTranslationFallback{prefixAuthor}{Handed in by}

\ProvideTranslation{german}{prefixDate}{am}
\ProvideTranslation{english}{prefixDate}{on}
\ProvideTranslationFallback{prefixDate}{on}

\ProvideTranslation{ngerman}{exercise}{Aufgabe}
\ProvideTranslation{english}{exercise}{Exercise}
\ProvideTranslationFallback{exercise}{Exercise}

\ProvideTranslation{ngerman}{exerciseAbbrev}{Aufg.}
\ProvideTranslation{english}{exerciseAbbrev}{Ex.}
\ProvideTranslationFallback{exerciseAbbrev}{Ex.}

\ProvideTranslation{ngerman}{subexercise}{Teilaufgabe}
\ProvideTranslation{english}{subexercise}{Subexercise}
\ProvideTranslationFallback{subexercise}{Subexercise}

\ProvideTranslation{ngerman}{points}{Punkte}
\ProvideTranslation{english}{points}{Points}
\ProvideTranslationFallback{points}{Points}

\ProvideTranslation{ngerman}{sheet}{Abgabe}
\ProvideTranslation{english}{sheet}{Exercise Sheet}
\ProvideTranslationFallback{sheet}{Exercise Sheet}

% FONTS
\usepackage{tgpagella}
% \usepackage{helvet}
% \renewcommand{\familydefault}{\sfdefault}

% LENGTHS
\setlength{\parindent}{0pt}
\setlength{\parskip}{0.6em}

% GRADING TABLE
\RequirePackage{totcount}
\RequirePackage{changepage}
\RequirePackage{etoolbox}
\RequirePackage{environ}
\RequirePackage{multicol}
\newtotcounter{exercisenumber}
\newtotcounter{subexercisenumber}
\newtotcounter{totalpoints}
\newcounter{task}

\newcommand{\addPT}[2]
{
    \edef\tempA{\GetTranslation{exerciseAbbrev} #1 &}
    \edef\tempB{\qquad / $#2$ &}
    \expandafter\gappto\expandafter\PTableA\expandafter{\tempA}
    \expandafter\gappto\expandafter\PTableB\expandafter{\tempB}
    \expandafter\gappto\expandafter\PTHead\expandafter{c | }
}

\newcommand{\gradingtable}{
    {\Large\bfseries Grading Table}
    \begin{table}[h!]
        \ifcsname PTableASaved\endcsname
        \renewcommand\arraystretch{2}
        \begin{tabular}{\PTHeadSaved c}
            \PTableASaved {\Large$\Sigma$}\\ \hline
            \PTableBSaved \qquad / $\totalpointsSaved$
        \end{tabular}
        \renewcommand{\arraystretch}{1}
        \else
            \texttt{Compile again}\ClassWarning{Compile again to proper display grading table}
        \fi
    \end{table}
}

\newcommand{\@pointsPerExercise}[1]{
    \ifstrempty{#1}{

    }{
        {\large\bfseries/#1 \GetTranslation{points}}
    }
}

% TASK
\newcommand{\exercise}[2][]{
    \penalty150
    \ifstrempty{#1}{\def\points{0}}{\def\points{#1}}
    \stepcounter{exercisenumber}
    \setcounter{subexercisenumber}{0}
    \par
    \vspace{1.5\baselineskip}
    {
        {\Large\bfseries\GetTranslation{exercise}\ \arabic{exercisenumber}}%
        \@exerciseLabel{#2}
        \@exerciseRules%
        \@pointsPerExercise{#1}
    }
    \par
    \addtocounter{totalpoints}{\points}
    \addPT{\arabic{exercisenumber}}{\points}
}

% Exercise rules
% fill remaining flushright space with a rule
\def\Vhrulefill{\leavevmode\leaders\hrule height 0.7ex depth \dimexpr0.4pt-0.7ex\hfill\kern0pt}
\newcommand{\exerciseRules}{
    \renewcommand{\@exerciseRules}{\hspace{1em}\Vhrulefill}
    \renewcommand{\@subexerciseRules}{\hspace{0em}\Vhrulefill}
}
\newcommand{\@exerciseRules}{}
\newcommand{\@subexerciseRules}{}

% Exercise label macro
\newcommand{\@exerciseLabel}[1]{
    \quad\ifstrempty{#1}{}{\bfseries(#1)}
}

% NO POINTS MACRO
% Strips points template for each exercise from title
\newcommand{\nopoints}{
    \renewcommand{\@pointsPerExercise}{}
    \renewcommand{\gradingtable}{}
    \renewcommand{\exercise}{\@exerciseNoPoints}
}

% Placeholder for \nopoints macro
\newcommand{\@exerciseNoPoints}[1]{
    \def\points{0}
    \stepcounter{exercisenumber}
    \setcounter{subexercisenumber}{0}
    \par
    \vspace{1.5\baselineskip}
    {
        {\Large\bfseries\GetTranslation{exercise}\ \arabic{exercisenumber}}%
        \@exerciseLabel{#1}
        \@exerciseRules
    }
    \par
    \addtocounter{totalpoints}{\points}
    \addPT{\arabic{exercisenumber}}{\points}
}

% Subexercise
% This is purely semantical, as subexercises do not have specific
% valuations and are not counted in the grading table
\newcommand{\subexercise}[2][]{
    \penalty150
    \stepcounter{subexercisenumber}
    \par
    \vspace{0.5\baselineskip}
    {
        {\bfseries\GetTranslation{subexercise}\ (\alph{subexercisenumber})}%
        \quad\ifstrempty{#2}{}{(#2)}%
        \@subexerciseRules%
    }
    \par
}

% SHEET
\newcommand{\printsheet}{?}
\newcommand{\sheet}[1]{\renewcommand{\printsheet}{\GetTranslation{sheet} #1}}

% DUE DATE
\newcommand{\printdue}{?}
\newcommand{\due}[1]{\renewcommand{\printdue}{#1}}

% COURSE
\newcommand{\printcourse}{?}
\newcommand{\course}[1]{\renewcommand{\printcourse}{#1}}

% GROUP
\newcommand{\printgroup}{?}
\newcommand{\group}[1]{\renewcommand{\printgroup}{#1}}

% GROUP MEMBERS
\newcommand{\printmembers}{}
\newcommand{\printmember}[3]{{#2} {#3}, Matr.Nr. {#1} \\}

\newcommand{\pdfmembers}{}
\newcommand{\pdfmember}[3]{{#1} {#3}, {#2}; }

\newcommand{\member}[3]{%
\expandafter\renewcommand\expandafter\printmembers\expandafter{\printmembers\printmember{#1}{#2}{#3}}%
\expandafter\renewcommand\expandafter\pdfmembers\expandafter{\pdfmembers\pdfmember{#1}{#2}{#3}}%
}

\newcommand{\@prefixAuthor}{\GetTranslation{prefixAuthor}}
\newcommand{\PrefixAuthor}[1]{\renewcommand{\@prefixAuthor}{#1}}

\newcommand{\@prefixDate}{\GetTranslation{prefixDate}}
\newcommand{\PrefixDate}[1]{\renewcommand{\@prefixDate}{#1}}

% \renewcommand{\@author}{\printmembers}
% \renewcommand{\@title}{\@course}
\renewcommand{\and}{\\}
\renewcommand{\maketitle}{
    \thispagestyle{empty}
    {\bfseries\Huge\printcourse}
    \par
    {\LARGE\printsheet}
    \par
    {\footnotesize\bfseries\@prefixAuthor}\\
    \hspace*{4pt}{\large{\vspace*{4pt}\begin{tabular}[t]{@{}l}
        \printmembers
    \end{tabular}}}\\
    {\footnotesize\bfseries\@prefixDate}\\
    \vspace*{4pt}{\hspace*{4pt}{\large\printdue}}\\
    \vskip \baselineskip
    \pagestyle{fancy}
}
% HEADER BAR
\renewcommand{\headrulewidth}{1pt}
\setlength\headheight{4\baselineskip}
\rhead{
    \begin{tabular}{r}
        \textbf{\printcourse} \\
        \printsheet \\
        {\small\printdue}
    \end{tabular}
}
\lhead{
    \small{\begin{tabular}{l}
        \printmembers
    \end{tabular}}
}

% PDF METADATA
\AtBeginDocument{\hypersetup{
    pdftitle = {Exercise Sheet \printsheet},
    pdfauthor = {Group \printgroup: \pdfmembers},
    pdfsubject = {Functional Programming}
}}

% AUXILLIARY FILES
\AtEndDocument{
  \makeatletter
  \immediate\write\@mainaux{\string\gdef\string\PTableASaved{\PTableA}}%
  \immediate\write\@mainaux{\string\gdef\string\PTableBSaved{\PTableB}}%
  \immediate\write\@mainaux{\string\gdef\string\PTHeadSaved{\PTHead}}%
  \immediate\write\@mainaux{\string\gdef\string\totalpointsSaved{\arabic{totalpoints}}}%
  \makeatother
}

% MATH
% These are mostly copied from amath:
% https://github.com/zoomoid/AlphabetClasses
\usepackage[fleqn]{amsmath}
\usepackage{amsthm}
\usepackage{amssymb}
\RequirePackage{mathtools}
% Number Spaces
\newcommand{\N}{\ensuremath{\mathbf{N}}}
\newcommand{\Z}{\ensuremath{\mathbf{Z}}}
\newcommand{\R}{\ensuremath{\mathbf{R}}}
\newcommand{\Q}{\ensuremath{\mathbf{Q}}}
\newcommand{\C}{\ensuremath{\mathbf{C}}}
\newcommand{\F}{\ensuremath{\mathbb{F}}}
\newcommand{\primefield}{\ensuremath{\mathbb{P}}}
% Set & Equation Semantics
\newcommand{\derivative}[1]{\ensuremath{\frac{d}{d#1}}}
\newcommand{\divides}{\ensuremath{\ |\ }}
\newcommand{\property}{\ensuremath{\ |\ }}
\newcommand{\modulofield}[2]{#1/\!#2\:}
% Combinatorical & Algebraic Sets
\DeclareMathOperator{\Perm}{Perm}
\DeclareMathOperator{\MComb}{MComb}
\DeclareMathOperator{\Comb}{Comb}
\DeclareMathOperator{\Pot}{Pot}
\DeclareMathOperator{\Map}{Map}
\DeclareMathOperator{\Hom}{Hom}
\DeclareMathOperator{\Ker}{Ker}
\DeclareMathOperator{\Intpol}{Intpol}
\DeclareMathOperator{\Pol}{Pol}
\DeclareMathOperator{\Sol}{Sol}
\DeclareMathOperator{\Col}{Col}
\DeclareMathOperator{\Bin}{Bin}
\DeclareMathOperator{\charac}{char}
% Matrix operations
\newcommand{\matadd}[3]{\ensuremath{\xmapsto{\mathrm{add}_{#1,#2,#3}}}}
\newcommand{\matmul}[2]{\ensuremath{\xmapsto{\mathrm{mul}_{#1,#2}}}}
\newcommand{\matswap}[2]{\ensuremath{\xmapsto{\mathrm{sw}_{#1,#2}}}}
% Logic quantors
\let\oforall\forall
\let\oexists\exists
\renewcommand{\forall}{\ensuremath{\hskip 2pt \oforall \hskip 2pt}}
\renewcommand{\exists}{\ensuremath{\hskip 2pt \oexists \hskip 2pt}}
\newcommand{\bigforall}{\mbox{\Large $\mathsurround3pt\forall$}}
\newcommand{\bigexists}{\mbox{\Large $\mathsurround3pt\exists$}}
% General Purpose functions
\newcommand{\ceil}[1]{\ensuremath{\left\lceil #1 \right \rceil}}
\newcommand{\floor}[1]{\ensuremath{\left\lfloor #1 \right \rfloor}}
\newcommand{\abs}[1]{\ensuremath{\left\vert#1\right\vert}}
\newcommand{\columnvector}[1]{\begin{pmatrix}#1\end{pmatrix}}
\newcommand*{\rfrac}[2]{{}^{#1}\!/_{#2}}
\renewcommand{\mod}{\ \mathrm{mod}\ }
% Semantic Relation Symbol
\newcommand{\rel}[1]{\ensuremath{\mathrel{#1}}}
% Large Brackets
\newcommand{\largebrackets}[1]{\ensuremath{\left(#1\right)}}
% Environments & Semantics
\newcommand{\QED}{\hfill\(\square\)}

% CODE using minted
%\usepackage[cache=false]{minted}

\newcommand{\Opt}{\textsc{Opt}}
\newcommand{\Adv}{\textsc{Adv}}
\newcommand{\Rand}{\textsc{Rand}}

\usepackage{amsthm}
\newtheoremstyle{it}% 〈name〉
{3pt}% 〈Space above〉1
{3pt}% 〈Space below 〉1
{}% 〈Body font〉
{}% 〈Indent amount〉2
{\itshape}% 〈Theorem head font〉
{:}% 〈Punctuation after theorem head 〉
{.5em}% 〈Space after theorem head 〉3
{}% 〈Theorem head spec (can be left empty, meaning ‘normal’ )〉
\theoremstyle{it}

\makeatother
\endinput
