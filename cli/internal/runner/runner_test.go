package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/zoomoid/assignments/v1/internal/config"
	"github.com/zoomoid/assignments/v1/internal/context"
	"github.com/zoomoid/assignments/v1/internal/util"
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
	dirName := fmt.Sprintf("assignment-%s", util.AddLeadingZero(cfg.Status.Assignment))
	err := os.Mkdir(filepath.Join(root, dirName), 0777)
	if err != nil {
		return "", err
	}

	fn := filepath.Join(root, dirName, "assignment.tex")
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

func TestRunner(t *testing.T) {
	workingDirectory := t.TempDir()
	targetDirectory, err := makeSourceFile(workingDirectory)
	if err != nil {
		t.Error(err)
		return
	}

	ctx, err := makeAppContext(workingDirectory)
	if err != nil {
		t.Error(err)
		return
	}

	testRunner, err := New(ctx, &RunnerOptions{
		TargetDirectory:   targetDirectory,
		OverrideArtifacts: false,
		Quiet:             false,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("SetRoot", func(t *testing.T) {
		newRoot := t.TempDir()
		testRunner.SetRoot(newRoot)

		if testRunner.root != newRoot {
			t.Error(fmt.Errorf("expected runner root to be %s, found %s", newRoot, testRunner.root))
			return
		}
		if !strings.HasPrefix(testRunner.ArtifactsDirectory(), newRoot) {
			t.Error(fmt.Errorf("expected runner artifacts directory to have prefix %s, found %s", newRoot, testRunner.ArtifactsDirectory()))
			return
		}
		if !strings.HasPrefix(testRunner.TargetDirectory(), newRoot) {
			t.Error(fmt.Errorf("expcted runner target directory to have prefix %s, found %s", newRoot, testRunner.TargetDirectory()))
		}
	})

	t.Run("SetTargetDirectory(abs)", func(t *testing.T) {
		newRoot := t.TempDir()
		newTargetDirectory := filepath.Join(newRoot, "assignment-02")
		testRunner.SetTargetDirectory(newTargetDirectory)
		if !filepath.IsAbs(testRunner.targetDirectory) {
			t.Error(fmt.Errorf("expected testRunner.TargetDirectory() to be absolute, found %s", testRunner.TargetDirectory()))
			return
		}
		if testRunner.TargetDirectory() != newTargetDirectory {
			t.Error(fmt.Errorf("expected testRunner.TargetDirectory() to be %s, found %s", newTargetDirectory, testRunner.TargetDirectory()))
			return
		}
		if testRunner.targetDirectory != testRunner.TargetDirectory() {
			t.Error(fmt.Errorf("expected testRunner.targetDirectory to be equal to %s, found %s", testRunner.TargetDirectory(), testRunner.targetDirectory))
			return
		}
	})

	t.Run("SetTargetDirectory(rel)", func(t *testing.T) {
		newTargetDirectory := "assignment-03"
		testRunner.SetTargetDirectory(newTargetDirectory)
		if testRunner.targetDirectory != newTargetDirectory {
			t.Error(fmt.Errorf("expected targetDirectory to be %s, found %s", newTargetDirectory, testRunner.targetDirectory))
			return
		}
		if strings.HasPrefix(testRunner.TargetDirectory(), testRunner.root) {
			t.Error(fmt.Errorf("expected testRunner.TargetDirectory() have prefix %s, found %s", testRunner.root, testRunner.TargetDirectory()))
			return
		}
	})

	t.Run("SetArtifactsDirectory", func(t *testing.T) {

	})
}
