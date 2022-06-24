package latexmk

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
	"go.uber.org/zap"
)

var (
	validAssignmentTexCode string = dedent.Dedent(`
		\documentclass{csassignments}
		\course{Online Algorithms}
		\group{TheCow}
		\member[369407]{Alexander Bartolomey}
		\member[366976]{Julius Rickert}
		\member[367129]{Adrian Hinrichs}
		\sheet{07}
		\due{June 03, 2022}
		
		\begin{document}
		\maketitle
		\gradingtable
		
		\exercise[4]{Modified Paging}
		
		\subexercise{}
		
		\begin{proof}
			
		Consider an instances \(I = P_{k+1},\dots,P_m,P_k,\dots,P_1,P_{k+1},\dots,P_m,P_k,\dots,P_1\)
		with \(p_1,\dots,p_k\) being initially selected, hence an adversary requesting all pages twice 
		in serial order.
		
		Then this causes the optimal algorithm to induce costs of \(\text{cost}(\textsc{Opt}(I)) = 2
		\cdot (2k + 1) = 2 \cdot (m - k + 1)\).
		
		\textsc{Lru} on the other hand makes a total of \(n = 2m = 3 \cdot 2k = 6k\) page faults, yielding 
		
		\[
			\lim_{k \rightarrow \infty} \frac{\text{cost}(\textsc{Lru}(I))}{\text{cost}(\textsc{Opt}(I))} = 
			\lim_{k \rightarrow \infty} \frac{6k}{4k+2} \rightarrow \frac{6}{4} = \frac{3}{2}.  
		\]
		
		For every \(c := \frac{3}{2} - \varepsilon\), \(\varepsilon > 0\), there exists a \(k_0 \in \N\) 
		such that for all \(k \geq k_0\), \textsc{Lru} is not \\ \(\frac{3}{2}-\varepsilon\)-competitive.
		Consequently, \textsc{Lru} is at best \(\frac{3}{2}\)-competitive.
		\end{proof}
		\end{document}
	`)
	cfg config.Configuration = config.Configuration{
		Spec: &config.ConfigurationSpec{
			Course: "Online Algorithms",
			Group:  "The Cow",
			Members: []config.GroupMember{
				{
					Name: "Alexander Bartolomey",
					ID:   "369407",
				}, {
					Name: "Julius Rickert",
					ID:   "366976",
				}, {
					Name: "Adrian Hinrichs",
					ID:   "367129",
				},
			},
			BuildOptions: &config.BuildOptions{
				Recipe: []config.Recipe{{
					Command: "pdflatex",
					Args:    []string{"-interaction=nonstopmode", "-file-line-error"},
				}},
			},
		},
		Status: &config.ConfigurationStatus{
			Assignment: 7,
		},
	}
)

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func makeAppContext(root string) (*context.AppContext, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer l.Sync()
	return &context.AppContext{
		Logger:        l.Sugar(),
		Cwd:           root,
		Root:          root,
		Configuration: &cfg,
	}, nil
}

func makeSourceFile(root string) (string, error) {
	dirName := filepath.Join(root)
	err := os.Mkdir(dirName, os.ModeDir)
	if err != nil {
		return "", err
	}

	fn := filepath.Join(dirName, "assignment.tex")
	f, err := os.Create(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(validAssignmentTexCode)
	if err != nil {
		return "", err
	}

	return dirName, nil
}

func TestValidRunnerFromRoot(t *testing.T) {
	workingDirectory := t.TempDir()

	targetDirectory, err := makeSourceFile(workingDirectory)
	if err != nil {
		t.Error(err)
	}

	ctx, err := makeAppContext(workingDirectory)
	if err != nil {
		t.Error(err)
	}

	testRunner, err := New(ctx, &RunnerOptions{
		TargetDirectory:   targetDirectory,
		OverrideArtifacts: true,
		Quiet:             false,
	})
	if err != nil {
		t.Error(err)
	}

	t.Run("Build", func(t *testing.T) {
		t.Run("assignmentNumber", func(t *testing.T) {
			r := testRunner.Clone().Build()
			t.Run("TargetDirectory=assignment", func(t *testing.T) {
				r.TargetDirectory = "assignment"
				i, err := r.assignmentNumber()
				if err != nil {
					// this is to be expected
					return
				}
				t.Error(fmt.Errorf("found assignment number where there shouldn't have, found %s, saw %s ", i, r.TargetDirectory))
			})
			t.Run("TargetDirectory=assignment-03", func(t *testing.T) {
				r.TargetDirectory = "assignment-03"
				s, err := r.assignmentNumber()
				if err != nil {
					t.Error(err)
					return
				}
				if s != "03" {
					t.Error(fmt.Errorf("Expected '%s', found '%s'", "03", s))
				}
			})
			t.Run("TargetDirectory=assignment-13", func(t *testing.T) {
				r.TargetDirectory = "assignment-13"
				s, err := r.assignmentNumber()
				if err != nil {
					t.Error(err)
					return
				}
				if s != "13" {
					t.Error(fmt.Errorf("Expected '%s', found '%s'", "13", s))
				}
			})
			t.Run("TargetDirectory=assignment-3", func(t *testing.T) {
				r.TargetDirectory = "assignment-3"
				s, err := r.assignmentNumber()
				if err != nil {
					t.Error(err)
					return
				}
				if s != "03" {
					t.Error(fmt.Errorf("Expected '%s', found '%s'", "03", s))
				}
			})
		})
		t.Run("makeBuildCommand", func(t *testing.T) {
			r := testRunner.Clone().Build()
			t.Run("predefined recipe", func(t *testing.T) {
				err = r.makeBuildCommands()
				if err != nil {
					t.Error(err)
					return
				}
				if len(r.Commands) != 1 {
					t.Error(fmt.Errorf("r.Commands is of length %d, expected %d ", len(r.Commands), 1))
					return
				}
				cmd := r.Commands[0]
				program := r.Configuration.Spec.BuildOptions.Recipe[0].Command
				if cmd.Args[0] != program {
					t.Error(fmt.Errorf("cmd runs %s, expected %s", cmd.Args[0], program))
					return
				}
			})
			t.Run("empty recipe", func(t *testing.T) {
				r.Configuration.Spec.BuildOptions.Recipe = []config.Recipe{}
				err = r.makeBuildCommands()
				if err != nil {
					t.Error(err)
					return
				}
				// defaults to single latexmk run
				if len(r.Commands) != len(r.Configuration.Spec.BuildOptions.Recipe) {
					t.Error(fmt.Errorf("r.Commands is of length %d, expected %d ", len(r.Commands), len(r.Configuration.Spec.BuildOptions.Recipe)))
					return
				}
				cmd := r.Commands[0]
				if cmd.Args[0] != defaultProgram {
					t.Error(fmt.Errorf("cmd runs %s, expected %s", cmd.Args[0], defaultProgram))
					return
				}
			})
			t.Run("pdflatex -> bibtex -> pdflatex recipe", func(t *testing.T) {
				r.Configuration.Spec.BuildOptions.Recipe = []config.Recipe{{
					Command: "pdflatex",
					Args:    []string{"-interaction=nonstopmode", "-file-line-error"},
				}, {
					Command: "bibtex",
				}, {
					Command: "pdflatex",
					Args:    []string{"-interaction=nonstopmode", "-file-line-error"},
				}}
				err = r.makeBuildCommands()
				if err != nil {
					t.Error(err)
					return
				}
				// defaults to single latexmk run
				if len(r.Commands) != len(r.Configuration.Spec.BuildOptions.Recipe) {
					t.Error(fmt.Errorf("r.Commands is of length %d, expected %d ", len(r.Commands), len(r.Configuration.Spec.BuildOptions.Recipe)))
					return
				}
				cmd := r.Commands[0]
				if cmd.Args[0] != defaultProgram {
					t.Error(fmt.Errorf("cmd runs %s, expected %s", cmd.Args[0], defaultProgram))
					return
				}
			})
		})
		t.Run("makeArtifactsDirectory", func(t *testing.T) {
			r := testRunner.Clone().Build()
			t.Run("does not exist", func(t *testing.T) {
				newRoot := t.TempDir() // get a fresh temporary directory
				r.Root = newRoot
				err := r.makeArtifactsDirectory()
				if err != nil {
					t.Error(err)
					return
				}
				// directory should exist now
				_, err = os.Stat(filepath.Join(newRoot, "dist"))
				if err != nil {
					t.Error(err)
					return
				}
				// Cleanup removes inner temp dir
			})
			t.Run("already exists", func(t *testing.T) {
				newRoot := t.TempDir()
				err := os.Mkdir(filepath.Join(newRoot, "dist"), os.ModeDir)
				if err != nil {
					t.Error(err)
					return
				}
				r.Root = newRoot
				// should return nil, as directory should already exist
				err = r.makeArtifactsDirectory()
				if err != nil {
					t.Error(err)
					return
				}
				// Cleanup removes inner temp dir
			})
		})
		// not actually building anything in between because that would
		// introduce a dependency on latex in testing - which is definitely
		// not what we want
		t.Run("exportArtifacts", func(t *testing.T) {
			r := testRunner.Clone().Build()
			// since we don't need an actual PDF at the src path,
			// just create a small test file named "assignment.pdf"
			newRoot := t.TempDir()
			r.Root = newRoot

			pdf, err := os.Create(filepath.Join(r.TargetDirectory, "assignment.pdf"))
			if err != nil {
				t.Error(err)
				return
			}
			_, err = pdf.WriteString(validAssignmentTexCode)
			if err != nil {
				t.Error(err)
				return
			}
			err = pdf.Close()
			if err != nil {
				t.Error(err)
				return
			}

			dest, err := r.exportArtifacts()
			if err != nil {
				t.Error(err)
				return
			}

			sfi, _ := os.Stat(filepath.Join(r.TargetDirectory, "assignment.pdf"))

			// check if assignment is in dist folder
			fi, err := os.Stat(dest)
			if err != nil {
				t.Error(err)
				return
			}

			if sfi.Size() != fi.Size() {
				t.Error(fmt.Errorf("src and dest files are not the same size, expected %d bytes, found %d bytes", sfi.Size(), fi.Size()))
			}
		})
	})
	t.Run("Clean", func(t *testing.T) {

	})
}
