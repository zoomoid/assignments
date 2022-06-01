\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{assignments}[2022/06/01 Assignment Class]
\DeclareOption*{\PassOptionsToClass{\CurrentOption}{article}}
\ProcessOptions \relax
\LoadClass[a4paper,11pt]{article}
\makeatletter
\usepackage[left=2.5cm,right=2.5cm,top=2cm,bottom=2cm,includeheadfoot]{geometry}
\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{titlesec}
\usepackage{enumitem}
\usepackage{graphicx}
\usepackage{tocloft}
\usepackage{float}
\usepackage{ifthen}
\usepackage{translations}
\usepackage{csquotes}
\usepackage{fancyhdr}
\usepackage{microtype}
\usepackage{stmaryrd}
\usepackage{pdfpages}
\usepackage[hidelinks, unicode]{hyperref}
\usepackage{tikz}
\usepackage{microtype}
\usepackage[fleqn]{amsmath}
\usepackage{amsthm}
\usepackage{amssymb}
\usepackage{mathtools}
\usepackage{totcount}
\usepackage{changepage}
\usepackage{etoolbox}
\usepackage{environ}
\usepackage{multicol}
\usepackage{tgpagella}
\usepackage{suffix}
% \usepackage{helvet}
% \renewcommand{\familydefault}{\sfdefault}

\usetikzlibrary{arrows,calc}

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

\ProvideTranslation{ngerman}{studentId}{Matr.Nr.}
\ProvideTranslation{english}{studentId}{Student ID}
\ProvideTranslationFallback{studentId}{Student ID}

% LENGTHS
\setlength{\parindent}{0pt}
\setlength{\parskip}{0.6em}

% GRADING TABLE
\newtotcounter{exercisenumber}
\newtotcounter{subexercisenumber}
\newtotcounter{totalpoints}
\newcounter{task}

\newcommand{\addPT}[2]{
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
\ifstrempty{#1}{}{{\quad\large\bfseries/#1 \GetTranslation{points}}}
}

% TASK
\newcommand\exercise[2][]{%
\penalty150%
\ifstrempty{#1}{\def\points{0}}{\def\points{#1}}%
\stepcounter{exercisenumber}%
\setcounter{subexercisenumber}{0}%
\par%
\vspace{1.5\baselineskip}%
{%
{\Large\bfseries\GetTranslation{exercise}\ \arabic{exercisenumber}}%
\@exerciseLabel{#2}%
\@exerciseRules%
\notblank{#1}{\@pointsPerExercise{#1}}{}%
}%
\par%
\addtocounter{totalpoints}{\points}%
\addPT{\arabic{exercisenumber}}{\points}%
}

\WithSuffix\newcommand\exercise*[1]{%
\penalty150%
\def\points{0}
\stepcounter{exercisenumber}%
\setcounter{subexercisenumber}{0}%
\par%
\vspace{1.5\baselineskip}%
{%
{\Large\bfseries\GetTranslation{exercise}\ \arabic{exercisenumber}}%
\@exerciseLabel{#1}%
\@exerciseRules%
}%
\par%
\addtocounter{totalpoints}{\points}%
}

\def\Vhrulefill{\leavevmode\leaders\hrule height 0.7ex depth \dimexpr0.4pt-0.7ex\hfill\kern0pt}
\newcommand{\exerciseRules}{
\renewcommand{\@exerciseRules}{\hspace{1em}\Vhrulefill}
\renewcommand{\@subexerciseRules}{\hspace{0em}\Vhrulefill}
}
\newcommand{\@exerciseRules}{}
\newcommand{\@subexerciseRules}{}

\newcommand{\@exerciseLabel}[1]{
\quad\ifstrempty{#1}{}{\bfseries(#1)}
}

\newcommand{\nopoints}{
\renewcommand{\@pointsPerExercise}{}
\renewcommand{\gradingtable}{}
\renewcommand{\exercise}{\@exerciseNoPoints}
}

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

\newcommand{\subexercise}[1][]{
\penalty150
\stepcounter{subexercisenumber}
\par
\vspace{0.5\baselineskip}
{
{\bfseries\GetTranslation{subexercise}\ (\alph{subexercisenumber})}%
\quad\ifstrempty{#1}{}{(#1)}%
\@subexerciseRules%
}
\par
}
\newcommand{\printsheet}{}
\newcommand{\sheet}[1]{\def\printsheet{\GetTranslation{sheet} #1}}
\newcommand{\printdue}{}
\newcommand{\due}[1]{\renewcommand{\printdue}{#1}}
\newcommand{\printcourse}{Assignment}
\newcommand{\course}[1]{\expandafter\notblank\expandafter{#1}{\def\printcourse{#1}}{}}
\newcommand{\printgroup}{}
\newcommand{\group}[1]{\renewcommand{\printgroup}{#1}}
\newcommand{\@prefixId}{\GetTranslation{studentId} }
\newcommand{\PrefixId}[1]{\renewcommand{\@prefixId}{#1}}
\newcommand{\printmembers}{}
\newcommand{\printmember}[2][]{\notblank{#1}{#2, \notblank{#1}{\@prefixId{#1}}{} \\}{#2}}
\newcommand{\pdfmembers}{}
\newcommand{\pdfmember}[2][]{\expandafter\notblank\expandafter{#1}{#2, {#1};}{#2}}
\newcommand{\member}[2][]{
\expandafter\renewcommand\expandafter\printmembers\expandafter{\printmembers\printmember[#1]{#2}}%
\expandafter\renewcommand\expandafter\pdfmembers\expandafter{\pdfmembers\pdfmember[#1]{#2}}%
}
\newcommand{\@prefixAuthor}{\GetTranslation{prefixAuthor}}
\newcommand{\PrefixAuthor}[1]{\renewcommand{\@prefixAuthor}{#1}}
\newcommand{\@prefixDate}{\GetTranslation{prefixDate}}
\newcommand{\PrefixDate}[1]{\renewcommand{\@prefixDate}{#1}}
% \renewcommand{\@author}{\printmembers}
% \renewcommand{\@title}{\@course}
\renewcommand{\and}{\\}
\renewcommand{\maketitle}{%
\thispagestyle{empty}%
{\bfseries\Huge\printcourse}%
\par%
\expandafter\notblank\expandafter\printsheet{{\LARGE\printsheet\\}}{}%
\par%
\expandafter\notblank\expandafter\printmembers{%
{\footnotesize\bfseries\@prefixAuthor}\\%
\hspace*{4pt}{\large{\vspace*{4pt}\begin{tabular}[t]{@{}l}%
\printmembers%
\end{tabular}}}\\%
}{}%
\expandafter\notblank\expandafter\printdue{%
{\footnotesize\bfseries\@prefixDate}\\%
\vspace*{4pt}{\hspace*{4pt}{\large\printdue}}\\%
}{}%
\vskip \baselineskip%
\pagestyle{fancy}%
}
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

\newcommand{\@pdfauthorprefix}{\expandafter\notblank\expandafter{\printgroup}{Group \printgroup: }{}\relax}

\AtBeginDocument{\hypersetup{
pdftitle = {\texorpdfstring{Assignment \printsheet}{}},
pdfauthor = {\texorpdfstring{\@pdfauthorprefix\pdfmembers}{}},
pdfsubject = {\texorpdfstring{\printcourse}{}}
}}

\AtEndDocument{
\makeatletter
\immediate\write\@mainaux{\string\gdef\string\PTableASaved{\PTableA}}%
\immediate\write\@mainaux{\string\gdef\string\PTableBSaved{\PTableB}}%
\immediate\write\@mainaux{\string\gdef\string\PTHeadSaved{\PTHead}}%
\immediate\write\@mainaux{\string\gdef\string\totalpointsSaved{\arabic{totalpoints}}}%
\makeatother
}

\newcommand{\N}{\ensuremath{\mathbf{N}}}
\newcommand{\Z}{\ensuremath{\mathbf{Z}}}
\newcommand{\R}{\ensuremath{\mathbf{R}}}
\newcommand{\Q}{\ensuremath{\mathbf{Q}}}
\newcommand{\C}{\ensuremath{\mathbf{C}}}
\newcommand{\F}{\ensuremath{\mathbb{F}}}
\newcommand{\primefield}{\ensuremath{\mathbb{P}}}
\newcommand{\derivative}[1]{\ensuremath{\frac{d}{d#1}}}
\newcommand{\divides}{\ensuremath{\ |\ }}
\newcommand{\property}{\ensuremath{\ |\ }}
\newcommand{\modulofield}[2]{#1/\!#2\:}
\newcommand{\matadd}[3]{\ensuremath{\xmapsto{\mathrm{add}_{#1,#2,#3}}}}
\newcommand{\matmul}[2]{\ensuremath{\xmapsto{\mathrm{mul}_{#1,#2}}}}
\newcommand{\matswap}[2]{\ensuremath{\xmapsto{\mathrm{sw}_{#1,#2}}}}
\let\oforall\forall
\let\oexists\exists
\renewcommand{\forall}{\ensuremath{\hskip 2pt \oforall \hskip 2pt}}
\renewcommand{\exists}{\ensuremath{\hskip 2pt \oexists \hskip 2pt}}
\newcommand{\bigforall}{\mbox{\Large $\mathsurround3pt\forall$}}
\newcommand{\bigexists}{\mbox{\Large $\mathsurround3pt\exists$}}
\newcommand{\ceil}[1]{\ensuremath{\left\lceil #1 \right \rceil}}
\newcommand{\floor}[1]{\ensuremath{\left\lfloor #1 \right \rfloor}}
\newcommand{\abs}[1]{\ensuremath{\left\vert#1\right\vert}}
\newcommand{\columnvector}[1]{\begin{pmatrix}#1\end{pmatrix}}
\newcommand*{\rfrac}[2]{{}^{#1}\!/_{#2}}
\renewcommand{\mod}{\ \mathrm{mod}\ }
\newcommand{\rel}[1]{\ensuremath{\mathrel{#1}}}
\newcommand{\largebrackets}[1]{\ensuremath{\left(#1\right)}}
\newcommand{\QED}{\hfill\(\square\)}
\newcommand{\Opt}{\textsc{Opt}}
\newcommand{\Adv}{\textsc{Adv}}
\newcommand{\Rand}{\textsc{Rand}}
\newtheoremstyle{it}{3pt}{3pt}{}{}{\itshape}{:}{.5em}{}
\theoremstyle{it}
\makeatother
\endinput
