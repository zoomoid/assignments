\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{assignments}[2022/06/01 Assignment Class]
\DeclareOption{sansserif}{\renewcommand{\familydefault}{\sfdefault}}
\DeclareOption*{\PassOptionsToClass{\CurrentOption}{article}}
\ProcessOptions\relax
\LoadClass[a4paper,11pt]{article}
\makeatletter
\RequirePackage[left=2.5cm,right=2.5cm,top=2cm,bottom=2cm,includeheadfoot]{geometry}
\RequirePackage[utf8]{inputenc}
\RequirePackage[english]{babel}
\RequirePackage{titlesec}
\RequirePackage{enumitem}
\RequirePackage{graphicx}
\RequirePackage{tocloft}
\RequirePackage{float}
\RequirePackage{ifthen}
\RequirePackage{translations}
\RequirePackage{csquotes}
\RequirePackage{fancyhdr}
\RequirePackage{microtype}
\RequirePackage{stmaryrd}
\RequirePackage{pdfpages}
\RequirePackage[hidelinks, unicode]{hyperref}
\RequirePackage{tikz}
\RequirePackage{microtype}
\RequirePackage[fleqn]{amsmath}
\RequirePackage{amsthm}
\RequirePackage{amssymb}
\RequirePackage{mathtools}
\RequirePackage{totcount}
\RequirePackage{changepage}
\RequirePackage{etoolbox}
\RequirePackage{environ}
\RequirePackage{multicol}
\RequirePackage{tgpagella}
\RequirePackage{suffix}
\RequirePackage{helvet}

\usetikzlibrary{arrows,calc}

\DeclareLanguage{english}
\DeclareLanguage{german}
\select@language{english}

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

\edef\exerciseLabelGradingTable{}
\edef\gradingTableExercisePoints{}

\edef\PTableA{}
\edef\PTableB{}
\edef\PTHead{}

\newcommand{\@addPT}[2]{
  \edef\gradingTableExerciseLabel{\GetTranslation{exerciseAbbrev} #1 &}
  \edef\gradingTableExercisePoints{\qquad / $#2$ &}
  \xappto\PTableA{\gradingTableExerciseLabel}
  \xappto\PTableB{\gradingTableExercisePoints}
  \xappto\PTHead{c | }
}

\newcommand{\gradingtable}{%
  \ifdef{\PTHeadSaved}{%
    {\Large\bfseries Grading Table}\\%
    \begin{table}[h!]
      \renewcommand\arraystretch{2}%
      \begin{tabular}{\expandonce{\PTHeadSaved} c}
        \PTableASaved {\Large$\Sigma$}\\ \hline
        \PTableBSaved \qquad / $\totalpointsSaved$
      \end{tabular}
      \renewcommand{\arraystretch}{1}%
    \end{table}
  }{\texttt{Compile again}\ClassWarning{Compile again to properly display grading table}}%
}

\newcommand{\@pointsPerExercise}[1]{
  \ifstrempty{#1}{}{{\quad\large\bfseries/#1 \GetTranslation{points}}}
}
\newcommand\exercise[2][]{%
  \penalty150%
  \ifblank{#1}{\def\points{0}}{\def\points{#1}}%
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
  \notblank{#1}{%
    \addtocounter{totalpoints}{\points}%
    \@addPT{\arabic{exercisenumber}}{\points}%
  }{}
}
\WithSuffix\newcommand\exercise*[2][]{%
  \penalty150%
  \ifblank{#1}{\def\points{0}}{\def\points{#1}}%
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
}
\AtEndDocument{
  \makeatletter
  \immediate\write\@mainaux{\string\gdef\string\PTableASaved{\PTableA}}%
  \immediate\write\@mainaux{\string\gdef\string\PTableBSaved{\PTableB}}%
  \immediate\write\@mainaux{\string\gdef\string\PTHeadSaved{\PTHead}}%
  \immediate\write\@mainaux{\string\gdef\string\totalpointsSaved{\arabic{totalpoints}}}%
  \makeatother
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
\def\Vhrulefill{\leavevmode\leaders\hrule height 0.7ex depth \dimexpr0.4pt-0.7ex\hfill\kern0pt}
\newcommand{\exerciseRules}[1][]{
  \renewcommand{\@exerciseRules}{\hspace{1em}\Vhrulefill}
  \ifx#1{all}{\renewcommand{\@subexerciseRules}{\hspace{0em}\Vhrulefill}\else{}\fi
  }
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
}

\newcommand{\@course}{Assignment}
\newcommand{\@sheet}{}
\newcommand{\@group}{}
\newcommand{\@prefixId}{\GetTranslation{studentId} }
\newcommand{\@members}{}
\newcommand{\@member}[2][]{\notblank{#1}{#2, \notblank{#1}{\@prefixId{#1}}{} \\}{#2}}
\newcommand{\@due}{}
\newcommand{\@prefixAuthor}{\GetTranslation{prefixAuthor}}
\newcommand{\@prefixDate}{\GetTranslation{prefixDate}}
\newcommand{\course}[1]{\expandafter\notblank\expandafter{#1}{\def\@course{#1}}{}}
\newcommand{\sheet}[1]{\def\@sheet{\GetTranslation{sheet} #1}}
\newcommand{\group}[1]{\renewcommand{\@group}{#1}}
\newcommand{\due}[1]{\renewcommand{\@due}{#1}}
\newcommand{\member}[2][]{
  \expandafter\renewcommand\expandafter\@members\expandafter{\@members\@member[#1]{#2}}%
  \expandafter\renewcommand\expandafter\pdfmembers\expandafter{\pdfmembers\pdfmember[#1]{#2}}%
}
\newcommand{\PrefixId}[1]{\renewcommand{\@prefixId}{#1}}
\newcommand{\PrefixAuthor}[1]{\renewcommand{\@prefixAuthor}{#1}}
\newcommand{\PrefixDate}[1]{\renewcommand{\@prefixDate}{#1}}
\renewcommand{\and}{\\}
\renewcommand{\maketitle}{%
  \thispagestyle{empty}%
  {\bfseries\Huge\@course}%
  \par%
  \expandafter\notblank\expandafter\@sheet{{\LARGE\@sheet\\}}{}%
  \par%
  \expandafter\notblank\expandafter\@members{%
    {\footnotesize\bfseries\@prefixAuthor}\\%
    \hspace*{4pt}{\large{\vspace*{4pt}\begin{tabular}[t]{@{}l}%
    \@members%
    \end{tabular}}}\\%
  }{}%
  \expandafter\notblank\expandafter\@due{%
    {\footnotesize\bfseries\@prefixDate}\\%
    \vspace*{4pt}{\hspace*{4pt}{\large\@due}}\\%
  }{}%
  \vskip \baselineskip%
  \pagestyle{fancy}%
}
\renewcommand{\headrulewidth}{1pt}
\setlength\headheight{4\baselineskip}
\rhead{
  \begin{tabular}{r}
    \textbf{\@course} \\
    \@sheet \\
    {\small\@due}
  \end{tabular}
}
\lhead{
  \small{%
    \begin{tabular}{l}
      \@members
    \end{tabular}%
  }
}

\newcommand{\@pdfauthorprefix}{\expandafter\notblank\expandafter{\@group}{Group \@group: }{}\relax}
\newcommand{\pdfmembers}{}
\newcommand{\pdfmember}[2][]{\expandafter\notblank\expandafter{#1}{#2, {#1};}{#2}}

\AtBeginDocument{\hypersetup{
  pdftitle = {\texorpdfstring{Assignment \@sheet}{}},
  pdfauthor = {\texorpdfstring{\@pdfauthorprefix\pdfmembers}{}},
  pdfsubject = {\texorpdfstring{\@course}{}}
}}

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
\newcommand{\modring}[2]{#1/\!#2\:}
\newcommand{\matadd}[3]{\ensuremath{\xmapsto{\mathrm{add}_{#1,#2,#3}}}}
\newcommand{\matmul}[2]{\ensuremath{\xmapsto{\mathrm{mul}_{#1,#2}}}}
\newcommand{\matswap}[2]{\ensuremath{\xmapsto{\mathrm{sw}_{#1,#2}}}}
\let\oforall\forall
\let\oexists\exists
\renewcommand{\forall}{\ensuremath{\hskip 2pt \oforall \hskip 2pt}}
\renewcommand{\exists}{\ensuremath{\hskip 2pt \oexists \hskip 2pt}}
\newcommand{\ceil}[1]{\ensuremath{\left\lceil #1 \right \rceil}}
\newcommand{\floor}[1]{\ensuremath{\left\lfloor #1 \right \rfloor}}
\newcommand{\abs}[1]{\ensuremath{\left\vert#1\right\vert}}
\newcommand{\rfrac}[2]{{}^{#1}\!/_{#2}}
\renewcommand{\mod}{\ \mathrm{mod}\ }
\newcommand{\rel}[1]{\ensuremath{\mathrel{#1}}}
\newcommand{\QED}{\hfill\(\square\)}
\newtheoremstyle{it}{3pt}{3pt}{}{}{\itshape}{:}{.5em}{}
\theoremstyle{it}
\makeatother
\endinput
